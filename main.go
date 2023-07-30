package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
)

var (
	peerTable = map[string]peer.ID{}
)

func main() {
	ctx := context.Background()

	args, err := parseArgs()
	if err != nil {
		panic(err)
	}
	h, err := newP2P(port1)
	if err != nil {
		panic(err)
	}
	defer h.Close()
	i, err := newif(name1, type1)
	if err != nil {
		panic(err)
	}
	defer i.Close()
	dMDNS, err := newMDNS(h, args.Token)
	if err != nil {
		panic(err)
	}
	// dDHT, err := newDHT(ctx, h, args.Token)
	// if err != nil {
	// 	panic(err)
	// }
	err = applyif(name1, args.Network, MTU)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Peer: %s\n", h.ID())
	fmt.Printf("VPN Address: %s\n\n", args.Network)

	go dMDNS.run(ctx, args.Network)
	// go dDHT.run(ctx, args.Network)

	h.SetStreamHandler(MeshProtocol.ID(), meshHandler(i))
	h.SetStreamHandler(MDNSProtocol.ID(), mdnsHandler)
	h.SetStreamHandler(DHTProtocol.ID(), dhtHandler)

	meshBridge(ctx, h, i)
}
