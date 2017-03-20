package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/performancecopilot/speed"
)

func main() {
	cpuIndom, err := speed.NewPCPInstanceDomain(
		"CPU Metrics",
		[]string{"CGoCalls", "Goroutines"},
	)
	if err != nil {
		panic(err)
	}

	cpuMetric, err := speed.NewPCPInstanceMetric(
		speed.Instances{
			"CGoCalls":   0,
			"Goroutines": 0,
		},
		"cpu",
		cpuIndom,
		speed.Int64Type,
		speed.CounterSemantics,
		speed.OneUnit,
	)
	if err != nil {
		panic(err)
	}

	memIndom, err := speed.NewPCPInstanceDomain(
		"Memory Metrics",
		[]string{
			"Alloc", "TotalAlloc", "Sys", "Lookups", "Mallocs", "Frees", "HeapAlloc",
			"HeapSys", "HeapIdle", "HeapInuse", "HeapReleased", "HeapObjects", "StackInuse",
			"StackSys", "MSpanInuse", "MSpanSys", "MCacheInuse", "MCacheSys", "BuckHashSys",
			"GCSys", "OtherSys", "NextGC", "LastGC", "PauseTotalNs", "PauseNs", "PauseEnd",
			"NumGC", "NumForcedGC",
		},
	)
	if err != nil {
		panic(err)
	}

	memInsts := speed.Instances{}
	for _, v := range memIndom.Instances() {
		memInsts[v] = 0
	}
	memMetric, err := speed.NewPCPInstanceMetric(
		memInsts,
		"mem",
		memIndom,
		speed.Uint64Type,
		speed.CounterSemantics,
		speed.OneUnit,
	)
	if err != nil {
		panic(err)
	}

	client, err := speed.NewPCPClient("runtime")
	if err != nil {
		panic(err)
	}

	client.MustRegister(cpuMetric)
	client.MustRegister(memMetric)
	client.MustStart()
	defer client.MustStop()

	mStats := runtime.MemStats{}

	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
			cpuMetric.MustSetInstance(runtime.NumCgoCall(), "CGoCalls")
			cpuMetric.MustSetInstance(runtime.NumGoroutine(), "Goroutines")

			runtime.ReadMemStats(&mStats)
			memMetric.MustSetInstance(mStats.Alloc, "Alloc")
			memMetric.MustSetInstance(mStats.TotalAlloc, "TotalAlloc")
			memMetric.MustSetInstance(mStats.Sys, "Sys")
			memMetric.MustSetInstance(mStats.Mallocs, "Mallocs")
			memMetric.MustSetInstance(mStats.Frees, "Frees")
			memMetric.MustSetInstance(mStats.HeapAlloc, "HeapAlloc")
			memMetric.MustSetInstance(mStats.HeapSys, "HeapSys")
			memMetric.MustSetInstance(mStats.HeapIdle, "HeapIdle")
			memMetric.MustSetInstance(mStats.HeapInuse, "HeapInuse")
			memMetric.MustSetInstance(mStats.HeapReleased, "HeapReleased")
			memMetric.MustSetInstance(mStats.HeapObjects, "HeapObjects")
			memMetric.MustSetInstance(mStats.StackInuse, "StackInuse")
			memMetric.MustSetInstance(mStats.StackSys, "StackSys")
			memMetric.MustSetInstance(mStats.MSpanInuse, "MSpanInuse")
			memMetric.MustSetInstance(mStats.MSpanSys, "MSpanSys")
			memMetric.MustSetInstance(mStats.MCacheInuse, "MCacheInuse")
			memMetric.MustSetInstance(mStats.MCacheSys, "MCacheSys")
			memMetric.MustSetInstance(mStats.BuckHashSys, "BuckHashSys")
			memMetric.MustSetInstance(mStats.GCSys, "GCSys")
			memMetric.MustSetInstance(mStats.OtherSys, "OtherSys")
			memMetric.MustSetInstance(mStats.NextGC, "NextGC")
			memMetric.MustSetInstance(mStats.LastGC, "LastGC")
			memMetric.MustSetInstance(mStats.PauseTotalNs, "PauseTotalNs")
			memMetric.MustSetInstance(mStats.PauseNs[(mStats.NumGC+255)%256], "PauseNs")
			memMetric.MustSetInstance(mStats.PauseEnd[(mStats.NumGC+255)%256], "PauseEnd")
			memMetric.MustSetInstance(uint64(mStats.NumGC), "NumGC")
			memMetric.MustSetInstance(uint64(mStats.NumForcedGC), "NumForcedGC")
		}
	}()
	defer ticker.Stop()

	fmt.Println("The metric is currently mapped as mmv.runtime, to stop the mapping, press enter")
	_, _ = os.Stdin.Read(make([]byte, 1))
}
