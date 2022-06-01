package main

import (
	"flag"

	"gitlab.com/akita/mgpusim/v2/samples/runner"
)

var benchmarkFlag = flag.String("benchmark", "fir",
	"Which benchmark to run")

func main() {
	flag.Parse()
	runner := new(runner.Runner).ParseFlag().Init()

	benchmark := SelectBenchmark(
		*benchmarkFlag, runner.Driver())

	runner.AddBenchmark(benchmark)
	runner.Run()
}
