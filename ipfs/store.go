package ipfs

import (
	// "errors"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/wschenk/archiver"
)

type IpfsService struct {
	shell *shell.Shell
}

func CreateService() *IpfsService {
	ipfs := shell.NewShell("localhost:5001")
	return &IpfsService{
		shell: ipfs,
	}
}

func (ipfs *IpfsService) IsPinned(hash string) (pinned bool, err error) {
	pinInfo, err := ipfs.shell.Pins()

	if err != nil {
		return false, err
	}

	info := pinInfo[hash]

	fmt.Printf("pin: %s: %v\n", hash, info)

	if info.Type == "" {
		return false, nil
	}
	return true, nil
}

func (ipfs *IpfsService) Get(hash string, locationPath string) error {
	return ipfs.shell.Get(hash, locationPath)
}

func (ipfs *IpfsService) Put(repo archiver.Repository) (hash string, err error) {
	return ipfs.shell.AddDir(repo.Dir())
}

/*

	version, _, err := ipfs.Version()

	if err != nil {
		panic(err)
	}

	fmt.Printf("IPFS Version = %s\n", version)

	// object, err := ipfs.ObjectGet("QmTf765bQweUropqo9yZiKTFh8BhYN8Bmwk3zZGM9TWkKX")

	// if err != nil {
	//  panic(err)
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
*/
