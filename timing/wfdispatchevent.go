package timing

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/gcn3"
	"gitlab.com/akita/gcn3/timing/wavefront"
)

// WfDispatchEvent is the event that the dispatcher dispatches a wavefront
type WfDispatchEvent struct {
	*akita.EventBase

	ManagedWf  *wavefront.Wavefront
	IsLastInWG bool
	MapWGReq   *gcn3.MapWGReq
}

// NewWfDispatchEvent creates a new WfDispatchCompletionEvent
func NewWfDispatchEvent(
	t akita.VTimeInSec,
	handler akita.Handler,
	Wf *wavefront.Wavefront,
) *WfDispatchEvent {
	evt := new(WfDispatchEvent)
	evt.EventBase = akita.NewEventBase(t, handler)
	evt.ManagedWf = Wf
	return evt
}
