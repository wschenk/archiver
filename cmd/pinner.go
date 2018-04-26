package main

import (
	"fmt"
	"github.com/wschenk/archiver"
	"github.com/wschenk/archiver/ipfs"
)

func main() {
	pinService := ipfs.CreateService()

	hash := "QmTf765bQweUropqo9yZiKTFh8BhYN8Bmwk3zZGdM9TWkKX"
	checkPinned(pinService, hash)
}

func checkPinned(s archiver.ArchiveStore, hash string) {
	fmt.Printf("Checking to see if %s is pinned...\n", hash)

	pinned, err := s.IsPinned(hash)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Pinned = %s\n", pinned)
}
