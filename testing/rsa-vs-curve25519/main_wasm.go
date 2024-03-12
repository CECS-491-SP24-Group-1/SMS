package main

import (
	"fmt"
	"syscall/js"
	"time"

	"wraith.me/testing/rsa-vs-curve-25519/lib"
)

func main() {
	//Create a blocking channel
	c := make(chan struct{})

	//Export functions for JS
	js.Global().Set("ed25519Keygen", js.FuncOf(ed25519Keygen))
	js.Global().Set("rsaKeygen", js.FuncOf(rsaKeygen))
	js.Global().Set("hello", js.FuncOf(hello))
	js.Global().Set("benchEd25519", js.FuncOf(benchEd25519))
	js.Global().Set("benchRSA", js.FuncOf(benchRSA))

	//Block the channel to prevent termination; termination will cause errors in the JS RT
	<-c
}

// hello() -> void
func hello(_ js.Value, _ []js.Value) interface{} {
	tyme := time.Now()
	js.Global().Call("alert", "Hello from Golang!")
	fmt.Printf("Time since alert: %s\n", time.Since(tyme))
	return nil
}

// ed25519Keygen() -> JSONObject[Ed25519KP]
func ed25519Keygen(_ js.Value, _ []js.Value) interface{} {
	//Create the Go object and marshal to JSON
	json := lib.Ed25519Keygen().JSON()

	//Call `JSON.parse()` on the string to derive an object usable in JS
	return js.Global().Get("JSON").Call("parse", json)
}

// rsaKeygen(bitSize: int) -> JSONObject[RSAKP]
func rsaKeygen(_ js.Value, args []js.Value) interface{} {
	//Get the key size from the args
	keysize := args[0].Int()

	//Create the Go object and marshal to JSON
	json := lib.RSAKeygen(keysize).JSON()

	//Call `JSON.parse()` on the string to derive an object usable in JS
	return js.Global().Get("JSON").Call("parse", json)
}

// benchEd25519(warmups: uint, runs: uint) -> JSONObject[Bench]
func benchEd25519(_ js.Value, args []js.Value) interface{} {
	//Get the warmup and runs count
	warmups := uint(args[0].Int())
	runs := uint(args[1].Int())

	//Do the benchmark
	bench := lib.NewBench(warmups, runs, func() any {
		lib.Ed25519Keygen()
		time.Sleep(1 * time.Nanosecond) //To prevent discards
		return nil
	})
	bench.Run()

	//Call `JSON.parse()` on the string to derive an object usable in JS
	return js.Global().Get("JSON").Call("parse", bench.String())
}

// benchEd25519(warmups: uint, runs: uint, keySize: int) -> JSONObject[Bench]
func benchRSA(_ js.Value, args []js.Value) interface{} {
	//Get the warmup and runs count
	warmups := uint(args[0].Int())
	runs := uint(args[1].Int())
	keySize := args[2].Int()

	//Do the benchmark
	bench := lib.NewBench(warmups, runs, func() any {
		lib.RSAKeygen(keySize)
		time.Sleep(1 * time.Nanosecond) //To prevent discards
		return nil
	})
	bench.Run()

	//Call `JSON.parse()` on the string to derive an object usable in JS
	return js.Global().Get("JSON").Call("parse", bench.String())
}
