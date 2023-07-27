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

	mdns, err := newMDNS(h, args.Token)
	if err != nil {
		panic(err)
	}
	fmt.Printf("VPN Address: %s\n\n", args.Network)

	ctx := context.Background()
	go mdns.run(ctx, args.Network)

	h.SetStreamHandler(MeshProtocol.ID(), meshHandler(i))
	h.SetStreamHandler(DiscoverProtocol.ID(), discoverHandler)

	applyif(name1, args.Network, MTU)

	meshBridge(ctx, h, i)
}
