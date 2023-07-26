package main

import (
	"bufio"
	"context"
	"encoding/binary"
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

func runMDNS(ctx context.Context, h host.Host, address string, peerChan chan peer.AddrInfo) {
	packetSize := make([]byte, 4)
	binary.BigEndian.PutUint32(packetSize, uint32(len(address)))

	for {
		peer := <-peerChan
		fmt.Println("Found peer:", peer, ", connected from: ", h.ID())

		if err := h.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}

		stream, err := h.NewStream(ctx, peer.ID, DiscoverProtocol.ID())

		if err != nil {
			fmt.Println("Stream open failed", err)
			continue
		}
		writer := bufio.NewWriter(stream)

		if _, err := writer.Write(packetSize); err != nil {
			fmt.Printf("Error sending message length: %v\n", err)
			continue
		}
		if _, err := writer.WriteString(address); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			continue
		}
		if err := writer.Flush(); err != nil {
			fmt.Printf("Error flushing writer: %v\n", err)
			continue
		}

		fmt.Println("Connected to:", peer)
	}
}
