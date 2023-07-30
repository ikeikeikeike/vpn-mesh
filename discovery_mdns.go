package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
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

func (n *discoveryMDNS) run(ctx context.Context, address string) {
	packetSize := make([]byte, 4)
	binary.BigEndian.PutUint32(packetSize, uint32(len(address)))

	for {
		p := <-n.PeerCh
		if err := n.host.Connect(ctx, p); err != nil {
			fmt.Println("MDNS Connection failed:", p.ID, err)
			continue
		}
		stream, err := n.host.NewStream(ctx, p.ID, MDNSProtocol.ID())
		if err != nil {
			fmt.Println("MDNS Stream open failed", p.ID, err)
			continue
		}
		writer := bufio.NewWriter(stream)

		if _, err := writer.Write(packetSize); err != nil {
			fmt.Printf("MDNS Error sending message length: %v\n", err)
			continue
		}
		if _, err := writer.WriteString(address); err != nil {
			fmt.Printf("MDNS Error sending message: %v\n", err)
			continue
		}
		if err := writer.Flush(); err != nil {
			fmt.Printf("MDNS Error flushing writer: %v\n", err)
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

func mdnsHandler(stream network.Stream) {
	packetSize := make([]byte, 4)

	for {
		if _, err := stream.Read(packetSize); err != nil {
			fmt.Printf("MDNS Error reading length from stream: %v\n", err)
			stream.Close()
			return
		}

		address := make([]byte, binary.BigEndian.Uint32(packetSize))
		if _, err := stream.Read(address); err != nil {
			fmt.Printf("MDNS Error reading message from stream: %v\n", err)
			stream.Close()
			return
		}
		ip, _, err := net.ParseCIDR(string(address))
		if err != nil {
			stream.Close()
			return
		}

		peerTable[ip.String()] = stream.Conn().RemotePeer()
		fmt.Printf("An address is now joined to vpn-mesh by MDNS: %s\n", address)
	}
}
