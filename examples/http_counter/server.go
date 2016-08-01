package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/performancecopilot/speed"
)

// TODO: replace the raw metric with a Counter once defined

var metric speed.SingletonMetric

func main() {
	var err error
	metric, err = speed.NewPCPSingletonMetric(
		0,
		"http.requests",
		speed.Int32Type,
		speed.CounterSemantics,
		speed.OneUnit,
		"Number of Requests",
	)
	if err != nil {
		panic(err)
	}

	client, err := speed.NewPCPClient("example", speed.ProcessFlag)
	if err != nil {
		panic(err)
	}

	client.MustRegister(metric)

	client.MustStart()
	defer client.MustStop()

	http.HandleFunc("/increment", handleIncrement)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	fmt.Println("To stop the server press enter")
	_, _ = os.Stdin.Read(make([]byte, 1))
	os.Exit(0)
}

func handleIncrement(w http.ResponseWriter, r *http.Request) {
	v := metric.Val().(int32)
	v++
	metric.MustSet(v)
	fmt.Fprintf(w, "incremented\n")
}
