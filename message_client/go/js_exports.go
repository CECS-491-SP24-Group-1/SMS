package main

import (
	"strings"
	"syscall/js"

	"wraith.me/clientside_crypto/lib"
	"wraith.me/clientside_crypto/util"

	ccrypto "wraith.me/message_server/crypto"
)

func main() {
	//Create a blocking channel
	c := make(chan struct{})

	//Golang -> JS
	//js.Global().Call("alert", "Hello from Golang!")

	//--Export functions for JS
	//Constructors
	js.Global().Set("ed25519Keygen", js.FuncOf(ed25519Keygen))
	js.Global().Set("ed25519FromBytes", js.FuncOf(ed25519FromBytes))
	js.Global().Set("ed25519FromJSON", js.FuncOf(ed25519FromJSON))
	js.Global().Set("ed25519FromSK", js.FuncOf(ed25519FromSK))

	//Methods
	js.Global().Set("ed25519Equal", js.FuncOf(ed25519Equal))
	js.Global().Set("ed25519JSON", js.FuncOf(ed25519JSON))
	js.Global().Set("ed25519Sign", js.FuncOf(ed25519Sign))
	js.Global().Set("ed25519String", js.FuncOf(ed25519String))
	js.Global().Set("ed25519Verify", js.FuncOf(ed25519Verify))

	//Functions

	//Block the channel to prevent termination; termination will cause errors in the JS RT
	<-c
}

/* //-- Constructors */
// ed25519Keygen() -> JSONObject[Ed25519KP]
func ed25519Keygen(_ js.Value, _ []js.Value) interface{} {
	//Create the Go object and marshal to JSON
	json := lib.Ed25519Keygen().JSON()

	//Call `JSON.parse()` on the string to derive an object usable in JS
	return js.Global().Get("JSON").Call("parse", json)
}

// ed25519FromBytes(sk []byte, pk []byte) -> JSONObject[Ed25519KP]
func ed25519FromBytes(_ js.Value, args []js.Value) interface{} {
	//Extract the secret and public keys from the args
	sk := util.JSArray2GoByteArray(args[0], lib.ED25519_LEN)
	pk := util.JSArray2GoByteArray(args[1], lib.ED25519_LEN)

	//Create an object from the bytes
	obj := lib.Ed25519FromBytes(sk, pk)

	//Call `JSON.parse()` to derive an object usable in JS
	return js.Global().Get("JSON").Call("parse", obj.JSON())
}

// ed25519FromJSON(jsons string) -> JSONObject[Ed25519KP]
func ed25519FromJSON(_ js.Value, args []js.Value) interface{} {
	obj := ed25519ObjFromJSON(args[0])
	return js.Global().Get("JSON").Call("parse", obj.JSON())
}

// func ed25519FromSK(sk []byte) Ed25519KP
func ed25519FromSK(_ js.Value, args []js.Value) interface{} {
	//Extract the secret from the args
	barr := util.JSArray2GoByteArray(args[0], lib.ED25519_LEN)

	//Derive a keypair object from the private and return it
	obj := lib.Ed25519FromSK(barr)
	return js.Global().Get("JSON").Call("parse", obj.JSON())
}

/* //-- Methods */
// ed25519Equal(us Ed25519KP, them Ed25519KP) -> bool
func ed25519Equal(_ js.Value, args []js.Value) interface{} {
	//Get the 2 objects to compare as JSON strings
	usstr := util.Val2Any[string](js.Global().Get("JSON").Call("stringify", args[0]))
	themstr := util.Val2Any[string](js.Global().Get("JSON").Call("stringify", args[1]))

	//Derive Go objects from the JSON strings
	us, erra := lib.Ed25519FromJSON(usstr)
	them, errb := lib.Ed25519FromJSON(themstr)

	//Ensure both objects parsed successfully before doing equality checks
	return erra == nil && errb == nil && us.Equal(them)
}

// ed25519JSON(this Ed25519KP) -> string
func ed25519JSON(_ js.Value, args []js.Value) interface{} {
	return ed25519ObjFromJSON(args[0]).JSON()
}

// ed25519JSON(this Ed25519KP, msg string) -> string
func ed25519Sign(_ js.Value, args []js.Value) interface{} {
	//Get the args as well-typed items
	obj := ed25519ObjFromJSON(args[0])
	msg := util.Val2Any[string](args[1])

	//Calculate the signature of the message
	sig := obj.Sign([]byte(msg))
	/*
		ta := js.Global().Get("Uint8Array").New(len(sig))
		ta.Call("set", util.GenerifyArray(sig[:]))
	*/

	//Return the signature
	return sig.String()
}

// func ed25519String(this Ed25519KP) -> string
func ed25519String(_ js.Value, args []js.Value) interface{} {
	return ed25519ObjFromJSON(args[0]).String()
}

// ed25519JSON(this Ed25519KP, msg string, sig string) -> bool
// TODO: ingest the signature as a string
func ed25519Verify(_ js.Value, args []js.Value) interface{} {
	//Get the args as well-typed items
	obj := ed25519ObjFromJSON(args[0])
	msg := util.Val2Any[string](args[1])
	sig := util.Val2Any[string](args[2])

	//Verify the signature and return the results
	return obj.Verify([]byte(msg), ccrypto.MustFromString[ccrypto.Signature](ccrypto.ParseSignature, sig))
}

/* //-- Utility functions */
//Converts a JSONObject representation of an Ed25519 keypair to its equivalent Go counterpart.
func ed25519ObjFromJSON(arg js.Value) lib.Ed25519KP {
	//Create a new Ed25519 keypair object
	obj := lib.Ed25519KP{}

	//Ingest the incoming JSONObject and stringify it, cleaning it up in the process
	str := util.Val2Any[string](js.Global().Get("JSON").Call("stringify", arg))
	strClean := strings.ReplaceAll(str, "\\", "")
	strClean = strClean[1 : len(strClean)-1] //This strips off the beginning and ending quotes

	//Attempt to derive a keypair object from the JSON
	//Only replace the original object if there are no errors
	parsed, err := lib.Ed25519FromJSON(strClean)
	if err == nil {
		obj = parsed
	}

	//Return the resultant Go object
	return obj
}
