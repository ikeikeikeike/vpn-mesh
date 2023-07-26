package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	address := os.Getenv("ADDRESS")
	if address == "" {
		panic(address)
	}

	h, err := newP2P(port1)
	if err != nil {
		panic(err)
	}
	defer h.Close()

	fmt.Printf("NodeID:%s\nNodeAddresses:%s\nVPNAddresse:%s\n", h.ID(), h.Addrs(), address)

	i, err := createInterface(name1, type1)
	if err != nil {
		panic(err)
	}
	defer i.Close()

	h.SetStreamHandler(MeshProtocol.ID(), meshHandler(i))
	h.SetStreamHandler(DiscoverProtocol.ID(), discoverHandler)

	peerChan, err := newMDNS(h, rendezvous)
	if err != nil {
		panic(err)
	}
	prepareInterface(name1, address, MTU)

	ctx := context.Background()
	go runMDNS(ctx, h, address, peerChan)

	readPackets(ctx, h, i)
}
