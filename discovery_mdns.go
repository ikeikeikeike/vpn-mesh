package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryMDNS struct {
	PeerCh chan peer.AddrInfo
	host   host.Host
}

func (n *discoveryMDNS) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerCh <- pi
}

func (n *discoveryMDNS) run(address string) {
	for {
		p := <-n.PeerCh
		if p.ID == n.host.ID() {
			continue
		}
		ctx := context.Background()

		if err := n.host.Connect(ctx, p); err != nil {
			// fmt.Println("MDNS Connection failed:", p.ID, ">>", err)
			continue
		}
		fmt.Printf("An address is now joined to vpn-mesh by MDNS: %s\n", p.ID)

		if err := discoveryWriter(ctx, n.host, address, p); err != nil {
			// fmt.Println("MDNS writer failed:", p.ID, ">>", err)
			continue
		}
	}
}

func newMDNS(h host.Host, rendezvous string) (*discoveryMDNS, error) {
	n := &discoveryMDNS{
		host:   h,
		PeerCh: make(chan peer.AddrInfo),
	}

	ser := mdns.NewMdnsService(h, rendezvous, n)
	if err := ser.Start(); err != nil {
		return nil, err
	}

	return n, nil
}
