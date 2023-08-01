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
)

func discoveryHandler(stream network.Stream) {
	packetSize := make([]byte, 4)

	for {
		if _, err := stream.Read(packetSize); err != nil {
			fmt.Printf("Error reading length from stream: %v\n", err)
			stream.Close()
			return
		}

		address := make([]byte, binary.BigEndian.Uint32(packetSize))
		if _, err := stream.Read(address); err != nil {
			fmt.Printf("Error reading message from stream: %v\n", err)
			stream.Close()
			return
		}
		ip, _, err := net.ParseCIDR(string(address))
		if err != nil {
			stream.Close()
			return
		}

		PeerTable[ip.String()] = stream.Conn().RemotePeer()
	}
}

func discoveryWriter(ctx context.Context, h host.Host, address string, p peer.AddrInfo) error {
	stream, err := h.NewStream(ctx, p.ID, DiscoveryProtocol.ID())
	if err != nil {
		return fmt.Errorf("Stream open failed %s: %w", p.ID, err)
	}
	packetSize := make([]byte, 4)
	binary.BigEndian.PutUint32(packetSize, uint32(len(address)))

	writer := bufio.NewWriter(stream)

	if _, err := writer.Write(packetSize); err != nil {
		return fmt.Errorf("Error sending message length: %w", err)
	}
	if _, err := writer.WriteString(address); err != nil {
		return fmt.Errorf("Error sending message: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("Error flushing writer: %w", err)
	}

	return nil
}
