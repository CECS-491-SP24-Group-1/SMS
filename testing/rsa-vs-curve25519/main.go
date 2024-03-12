package main

import (
	"fmt"
	"time"

	"wraith.me/testing/rsa-vs-curve-25519/lib"
)

func main() {
	rsa := lib.RSAKeygen(1024)

	//fmt.Printf("%+v\n", rsa.SK)

	fmt.Printf("SK size: %d\n", len(rsa.SK))
	fmt.Printf("PK size: %d\n", len(rsa.PK))

	tyme := time.Now()
	lib.Ed25519Keygen()
	delta := time.Since(tyme)
	fmt.Printf("EdKG: %s\n", delta.String())

	bench1 := lib.NewBench(5, 5, func() any {
		//return lib.RSAKeygen(2048)
		time.Sleep(250 * time.Millisecond)
		return nil
	})
	//bench1.Run()
	fmt.Println(bench1.String())

	bench2 := lib.NewBench(5, 5, func() any {
		lib.Ed25519Keygen()
		time.Sleep(250 * time.Millisecond)
		return nil
	})
	bench2.Run()
	fmt.Println(bench2.String())
}
