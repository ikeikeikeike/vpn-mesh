package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	// Arguments
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	// P2P Host
	h, err := newP2P(args.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	// Tun interface
	i, err := newif(args.Interface)
	if err != nil {
		log.Fatal(err)
	}
	defer i.Close()
	err = applyif(args.Interface, args.Network)
	if err != nil {
		log.Fatal(err)
	}

	// For local-network
	dMDNS, err := newMDNS(h, args.Rendezvous)
	if err != nil {
		log.Fatal(err)
	}
	// For global-network
	dDHT, err := newDHT(ctx, h, args.Rendezvous)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Peer: %s\n", h.ID())
	fmt.Printf("VPN Address: %s\n\n", args.Network)

	// Discover
	go dMDNS.run(args.Network)
	go dDHT.run(args.Network)

	// Packet
	h.SetStreamHandler(MeshProtocol.ID(), meshHandler(i))
	h.SetStreamHandler(DiscoveryProtocol.ID(), discoveryHandler)
	meshBridge(ctx, h, i)
}
