package main

import (
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
)

var topic = "thisismytopic"

func main() {
	fmt.Println("Hello world")

	ipfs := shell.NewShell("localhost:5001")

	version, _, err := ipfs.Version()

	if err != nil {
		panic(err)
	}

	fmt.Printf("IPFS Version = %s\n", version)

	// object, err := ipfs.ObjectGet("QmTf765bQweUropqo9yZiKTFh8BhYN8Bmwk3zZGM9TWkKX")

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(object)

	pubsub, err := ipfs.PubSubSubscribe(topic)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on %s\n", topic)

	data, err := pubsub.Next()
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s] %s\n", data.From(), data.Data())

	data, err = pubsub.Next()
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s] %s\n", data.From(), data.Data())

	data, err = pubsub.Next()
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s] %s\n", data.From(), data.Data())

	err = pubsub.Cancel()
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}
