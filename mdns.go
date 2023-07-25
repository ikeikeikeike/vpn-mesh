package main

import (
	"bufio"
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

func newMDNS(h host.Host, rendezvous string) (chan peer.AddrInfo, error) {
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	ser := mdns.NewMdnsService(h, rendezvous, n)
	if err := ser.Start(); err != nil {
		return nil, err
	}

	return n.PeerChan, nil
}

func runMDNS(ctx context.Context, h host.Host, peerChan chan peer.AddrInfo) {
	for { // allows multiple peers to join
		println("mdns", 0)
		peer := <-peerChan // will block until we discover a peer
		fmt.Println("Found peer:", peer, ", connecting")

		if err := h.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}

		// open a stream, this stream will be handled by handleStream other end
		stream, err := h.NewStream(ctx, peer.ID, Protocol)

		if err != nil {
			fmt.Println("Stream open failed", err)
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			fmt.Println("Connected to:", peer, rw)
		}
	}

}
