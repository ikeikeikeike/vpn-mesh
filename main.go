package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	args, err := parseArgs()
	if err != nil {
		panic(err)
	}
	h, err := newP2P(args.Port)
	if err != nil {
		panic(err)
	}
	defer h.Close()
	i, err := newif(args.Interface)
	if err != nil {
		panic(err)
	}
	defer i.Close()
	dMDNS, err := newMDNS(h, args.Rendezvous)
	if err != nil {
		panic(err)
	}
	dDHT, err := newDHT(ctx, h, args.Rendezvous)
	if err != nil {
		panic(err)
	}
	err = applyif(args.Interface, args.Network)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Peer: %s\n", h.ID())
	fmt.Printf("VPN Address: %s\n\n", args.Network)

	go dMDNS.run(args.Network)
	go dDHT.run(args.Network)

	h.SetStreamHandler(MeshProtocol.ID(), meshHandler(i))
	h.SetStreamHandler(DiscoveryProtocol.ID(), discoveryHandler)

	meshBridge(ctx, h, i)
}
