package dispatching

import (
	"fmt"
	"sync"

	"gitlab.com/akita/akita/v2/monitoring"
	"gitlab.com/akita/akita/v2/sim"
	"gitlab.com/akita/mgpusim/v2/kernels"
	"gitlab.com/akita/mgpusim/v2/protocol"
	"gitlab.com/akita/mgpusim/v2/timing/cp/internal/resource"
	"gitlab.com/akita/util/v2/tracing"
)

type projectionDecisionMaker interface {
	StartKernel(*protocol.LaunchKernelReq)
	StartWG(wg *protocol.MapWGReq)
	EndWG(now sim.VTimeInSec, wg *protocol.MapWGReq)
	DoProjection() bool
	MeanWGTerminationGap() float64
}

type windowedProjectionDecisionMaker struct {
	skipHead              int
	windowSize            int
	wavefrontPerWG        uint32
	wgTerminationGaps     []float64
	wgCount               int
	ciSizeInSigma         float64
	errorMargin           float64
	mean                  float64
	stdev                 float64
	ciLow, ciHigh         float64
	lastWGTerminationTime float64
	projectionMode        bool
}

func newWindowedProjectionDecisionMaker() *windowedProjectionDecisionMaker {
	waveSize := 64 * 40 // 64 CU, 40 WG per CU
	// waveSize := 64
	numWaves := 4
	safeGuard := 0.2
	return &windowedProjectionDecisionMaker{
		skipHead:       int(float64(waveSize) * (1 + safeGuard)),
		windowSize:     int(float64(waveSize*numWaves) * (1 + safeGuard)),
		ciSizeInSigma:  3.0,
		errorMargin:    0.02,
		projectionMode: false,
	}
}

func (dm *windowedProjectionDecisionMaker) StartKernel(
	req *protocol.LaunchKernelReq,
) {
	dm.projectionMode = false
	dm.wgTerminationGaps = make([]float64, 0)
	dm.wgCount = 0

	threadPerWavefront := uint32(64)
	gridSize := req.Packet.GridSizeX *
		req.Packet.GridSizeY *
		req.Packet.GridSizeZ
	blockSize := req.Packet.WorkgroupSizeX *
		req.Packet.WorkgroupSizeY *
		req.Packet.WorkgroupSizeZ
	dm.wavefrontPerWG = gridSize / uint32(blockSize) / threadPerWavefront

	waveSize := int(64 * 40 / dm.wavefrontPerWG) // 64 CU, 40 WG per CU
	// waveSize := 64
	numWaves := 4
	safeGuard := 0.2
	dm.skipHead = int(float64(waveSize) * (1 + safeGuard))
	dm.windowSize = int(float64(waveSize*numWaves) * (1 + safeGuard))
}

func (dm *windowedProjectionDecisionMaker) StartWG(
	req *protocol.MapWGReq,
) {
	// Do nothing.
}

func (dm *windowedProjectionDecisionMaker) EndWG(
	now sim.VTimeInSec,
	req *protocol.MapWGReq,
) {
	gap := float64(now)
	if dm.lastWGTerminationTime != 0 {
		gap = float64(now) - dm.lastWGTerminationTime
	}
	dm.lastWGTerminationTime = float64(now)

	dm.wgCount++
	if dm.wgCount <= dm.skipHead {
		return
	}

	dm.wgTerminationGaps = append(dm.wgTerminationGaps, gap)

	if len(dm.wgTerminationGaps) > dm.windowSize {
		dm.wgTerminationGaps = dm.wgTerminationGaps[1:]
	}

	dm.processStats(float64(now))
}

func (dm *windowedProjectionDecisionMaker) processStats(now float64) {
	if len(dm.wgTerminationGaps) < dm.windowSize {
		return
	}

	// Calculate mean and stdev.
	var sum float64
	for _, gap := range dm.wgTerminationGaps {
		sum += gap
	}
	dm.mean = sum / float64(len(dm.wgTerminationGaps))

	var sum2 float64
	for _, gap := range dm.wgTerminationGaps {
		sum2 += (gap - dm.mean) * (gap - dm.mean)
	}
	dm.stdev = sum2 / float64(len(dm.wgTerminationGaps))

	// Calculate confidence interval.
	dm.ciLow = dm.mean - dm.ciSizeInSigma*dm.stdev
	dm.ciHigh = dm.mean + dm.ciSizeInSigma*dm.stdev

	ciRange := dm.ciHigh - dm.ciLow
	ciRangePercent := ciRange / dm.mean

	if ciRangePercent < dm.errorMargin {
		dm.projectionMode = true
	}

	// fmt.Printf("%.10f, %d, %f, %f, %f, %f, %f, %f, %t\n",
	// 	now, dm.wgCount, dm.mean, dm.stdev, dm.ciLow, dm.ciHigh, ciRange, ciRangePercent, dm.projectionMode)
}

func (dm *windowedProjectionDecisionMaker) DoProjection() bool {
	return dm.projectionMode
}

func (dm *windowedProjectionDecisionMaker) MeanWGTerminationGap() float64 {
	return dm.mean
}

type nextWGEvent struct {
	time    sim.VTimeInSec
	handler sim.Handler
}

func (e nextWGEvent) Time() sim.VTimeInSec {
	return e.time
}

func (e nextWGEvent) Handler() sim.Handler {
	return e.handler
}

func (e nextWGEvent) IsSecondary() bool {
	return false
}

// A wgProjectionDispatcher only dispatches work-groups at the beginning of
// a kernel. It will project the execution fot the rest of the kernel.
type wgProjectionDispatcher struct {
	sync.Mutex
	sim.HookableBase

	engine sim.Engine

	cp                       tracing.NamedHookable
	name                     string
	respondingPort           sim.Port
	dispatchingPort          sim.Port
	alg                      algorithm
	dispatching              *protocol.LaunchKernelReq
	currWG                   dispatchLocation
	cycleLeft                int
	numWGDispatchedInSimMode int
	numWGCompletedInSimMode  int
	inflightWGs              map[string]dispatchLocation
	originalReqs             map[string]*protocol.MapWGReq
	latencyTable             []int
	constantKernelOverhead   int

	projectionDecisionMaker            projectionDecisionMaker
	freq                               sim.Freq
	projectionStarted                  bool
	projectionEnded                    bool
	numWGCompletedWhenProjectionStart  int
	numWGDispatchedWhenProjectionStart int
	numWGCompletedInProjectionMode     int

	monitor     *monitoring.Monitor
	progressBar *monitoring.ProgressBar
}

// Name returns the name of the dispatcher
func (d *wgProjectionDispatcher) Name() string {
	return d.name
}

// RegisterCU allows the dispatcher to dispatch work-groups to the CU.
func (d *wgProjectionDispatcher) RegisterCU(cu resource.DispatchableCU) {
	d.alg.RegisterCU(cu)
}

// IsDispatching checks if the dispatcher is dispatching another kernel.
func (d *wgProjectionDispatcher) IsDispatching() bool {
	return d.dispatching != nil
}

// StartDispatching lets the dispatcher to start dispatch another kernel.
func (d *wgProjectionDispatcher) StartDispatching(
	req *protocol.LaunchKernelReq,
) {
	d.mustNotBeDispatchingAnotherKernel()

	d.alg.StartNewKernel(kernels.KernelLaunchInfo{
		CodeObject: req.HsaCo,
		Packet:     req.Packet,
		PacketAddr: req.PacketAddress,
		WGFilter:   req.WGFilter,
	})
	d.dispatching = req

	d.numWGDispatchedInSimMode = 0
	d.numWGCompletedInSimMode = 0

	d.projectionDecisionMaker.StartKernel(req)

	d.initializeProgressBar(req.ID)
}

func (d *wgProjectionDispatcher) initializeProgressBar(kernelID string) {
	if d.monitor != nil {
		d.progressBar = d.monitor.CreateProgressBar(
			fmt.Sprintf("At projection dispatcher %s, Kernel: %s, ", d.Name(), kernelID),
			uint64(d.alg.NumWG()),
		)
	}
}

func (d *wgProjectionDispatcher) mustNotBeDispatchingAnotherKernel() {
	if d.IsDispatching() {
		panic("dispatcher is dispatching another request")
	}
}

func (d *wgProjectionDispatcher) Handle(e sim.Event) error {
	d.Lock()
	defer d.Unlock()

	switch e := e.(type) {
	case nextWGEvent:
		return d.handleNextWGEvent(e)
	}

	panic("never")
}

func (d *wgProjectionDispatcher) handleNextWGEvent(e nextWGEvent) error {
	now := e.Time()

	d.numWGCompletedInProjectionMode++

	if d.isSimulationCompleted() && d.isProjectionCompleted() {
		ok := d.completeKernel(now)
		if !ok {
			panic("kernel is not completed")
		}

		return nil
	}

	if d.isProjectionCompleted() {
		return nil
	}

	newEvt := nextWGEvent{
		time: now +
			sim.VTimeInSec(d.projectionDecisionMaker.MeanWGTerminationGap()),
		handler: d,
	}
	d.engine.Schedule(newEvt)

	if d.progressBar != nil {
		// fmt.Printf("proj wg complete, %.10f\n", now)
		d.progressBar.IncrementFinished(1)
	}

	return nil
}

func (d *wgProjectionDispatcher) isProjectionCompleted() bool {
	if !d.projectionStarted {
		return true
	}

	return d.numWGCompletedInProjectionMode+
		d.numWGCompletedWhenProjectionStart >=
		d.alg.NumWG()
}

func (d *wgProjectionDispatcher) isSimulationCompleted() bool {
	return d.numWGCompletedInSimMode >= d.numWGDispatchedInSimMode
}

// Tick updates the state of the dispatcher.
func (d *wgProjectionDispatcher) Tick(now sim.VTimeInSec) (madeProgress bool) {
	d.Lock()
	defer d.Unlock()

	if d.cycleLeft > 0 {
		d.cycleLeft--
		return true
	}

	if d.dispatching != nil {
		madeProgress = d.dispatchNextWG(now) || madeProgress
	}

	madeProgress = d.processMessagesFromCU(now) || madeProgress

	return madeProgress
}

func (d *wgProjectionDispatcher) processMessagesFromCU(
	now sim.VTimeInSec,
) bool {
	item := d.dispatchingPort.Peek()
	if item == nil {
		return false
	}

	msg, ok := item.(*protocol.WGCompletionMsg)
	if !ok {
		return false
	}

	location, ok := d.inflightWGs[msg.RspTo]
	if !ok {
		return false
	}

	d.numWGCompletedInSimMode++

	if d.isSimulationCompleted() && d.isProjectionCompleted() {
		ok := d.completeKernel(now)
		if !ok {
			panic("kernel is not completed")
		}
	}

	d.alg.FreeResources(location)
	delete(d.inflightWGs, msg.RspTo)

	d.dispatchingPort.Retrieve(now)

	originalReq := d.originalReqs[msg.RspTo]
	delete(d.originalReqs, msg.RspTo)
	tracing.TraceReqFinalize(originalReq, now, d)
	d.projectionDecisionMaker.EndWG(now, originalReq)

	if !d.projectionStarted && d.progressBar != nil {
		// fmt.Printf("wg complete, %.10f\n", now)
		d.progressBar.MoveInProgressToFinished(1)
	}

	return true
}

func (d *wgProjectionDispatcher) completeKernel(now sim.VTimeInSec) (
	madeProgress bool,
) {
	req := d.dispatching

	if req == nil {
		return true
	}

	req.Src, req.Dst = req.Dst, req.Src
	req.SendTime = now
	err := d.respondingPort.Send(req)

	if err == nil {
		d.dispatching = nil

		if d.monitor != nil {
			d.monitor.CompleteProgressBar(d.progressBar)
		}

		tracing.TraceReqComplete(req, now, d.cp)

		return true
	}

	req.Src, req.Dst = req.Dst, req.Src
	return false
}

func (d *wgProjectionDispatcher) dispatchNextWG(
	now sim.VTimeInSec,
) (madeProgress bool) {
	if !d.currWG.valid {
		if !d.alg.HasNext() {
			return false
		}

		d.currWG = d.alg.Next()
		if !d.currWG.valid {
			return false
		}
	}

	if d.projectionDecisionMaker.DoProjection() {
		if !d.projectionStarted {
			d.projectionStarted = true
			d.numWGDispatchedWhenProjectionStart = d.numWGDispatchedInSimMode
			d.numWGCompletedWhenProjectionStart = d.numWGCompletedInSimMode

			timePerWG := d.projectionDecisionMaker.MeanWGTerminationGap()
			evt := nextWGEvent{
				time:    now + sim.VTimeInSec(timePerWG),
				handler: d,
			}
			d.engine.Schedule(evt)
		}
	} else {
		return d.sendWGReq(now)
	}

	return false
}

func (d *wgProjectionDispatcher) sendWGReq(now sim.VTimeInSec) (
	madeProgress bool,
) {
	reqBuilder := protocol.MapWGReqBuilder{}.
		WithSrc(d.dispatchingPort).
		WithDst(d.currWG.cu).
		WithSendTime(now).
		WithPID(d.dispatching.PID).
		WithWG(d.currWG.wg)
	for _, l := range d.currWG.locations {
		reqBuilder = reqBuilder.AddWf(l)
	}
	req := reqBuilder.Build()
	err := d.dispatchingPort.Send(req)

	if err == nil {
		d.currWG.valid = false
		d.numWGDispatchedInSimMode++
		d.inflightWGs[req.ID] = d.currWG
		d.originalReqs[req.ID] = req
		d.cycleLeft = d.latencyTable[len(d.currWG.locations)]

		if d.progressBar != nil {
			d.progressBar.IncrementInProgress(1)
		}

		tracing.TraceReqInitiate(req, now, d,
			tracing.MsgIDAtReceiver(d.dispatching, d.cp))

		d.projectionDecisionMaker.StartWG(req)

		return true
	}

	return false
}
