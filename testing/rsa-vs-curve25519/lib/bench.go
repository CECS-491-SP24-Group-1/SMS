package lib

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Bench struct {
	Warmups       uint   `json:"warmups"`
	Runs          uint   `json:"runs"`
	Average       int64  `json:"average"`
	AveragePretty string `json:"average_pretty"`
	function      func() any
	splits        []time.Duration
}

func NewBench(warmups uint, runs uint, function func() any) Bench {
	out := Bench{
		Warmups:  warmups,
		Runs:     runs,
		function: function,
		splits:   make([]time.Duration, warmups+runs),
	}
	return out
}

func (b *Bench) Run() {
	//Loop x times
	for i := range b.splits {
		//Get the start time
		start := time.Now()

		//Run the function
		runtime.KeepAlive(b.function())

		//Get the end time and add it to the split
		delta := time.Since(start)
		b.splits[i] = delta
		fmt.Printf("split: %s\n", b.splits[i].String())
	}

	//Calculate the average of the splits, but only starting from the non-warmup splits
	var avg int64 = 0
	for _, split := range b.splits[b.Warmups:] {
		avg += split.Nanoseconds()
	}

	b.Average = avg / int64(b.Runs)
	b.AveragePretty = time.Duration(b.Average).String()
}

func (b Bench) String() string {
	json, _ := json.Marshal(b)
	return string(json)
}
