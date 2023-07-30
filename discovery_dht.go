package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

type discoveryDHT struct {
	PeerCh <-chan peer.AddrInfo
	host   host.Host
}

func (n *discoveryDHT) run(ctx context.Context, address string) {
	packetSize := make([]byte, 4)
	binary.BigEndian.PutUint32(packetSize, uint32(len(address)))

	for p := range n.PeerCh {
		if p.ID == n.host.ID() || len(p.Addrs) == 0 {
			continue
		}

		switch n.host.Network().Connectedness(p.ID) {
		case network.NotConnected:
			if err := n.host.Connect(ctx, p); err != nil {
				fmt.Println("DHT Connection failed:", p.ID, ">>", err)
				continue
			}
		case network.Connected:
			stream, err := n.host.NewStream(ctx, p.ID, DHTProtocol.ID())
			if err != nil {
				fmt.Println("DHT Stream open failed:", p.ID, ">>", err)
				continue
			}
			writer := bufio.NewWriter(stream)

			if _, err := writer.Write(packetSize); err != nil {
				fmt.Printf("DHT Error sending message length: %v\n", err)
				continue
			}
			if _, err := writer.WriteString(address); err != nil {
				fmt.Printf("DHT Error sending message: %v\n", err)
				continue
			}
			if err := writer.Flush(); err != nil {
				fmt.Printf("DHT Error flushing writer: %v\n", err)
				continue
			}
		}
	}
}

func newDHT(ctx context.Context, h host.Host, rendezvous string) (*discoveryDHT, error) {
	kadDHT, err := dht.New(ctx, h)
	if err != nil {
		return nil, err
	}
	if err = kadDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}
	maddr, err := multiaddr.NewMultiaddr(
		"/ip4/104.131.131.82/udp/4001/quic/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	)
	if err != nil {
		return nil, err
	}
	boots := append(dht.DefaultBootstrapPeers, maddr)

	var wg sync.WaitGroup
	for _, pa := range boots {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(pa)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerinfo); err != nil {
				fmt.Printf("DHT Bootstrap Connection failed: %v\n", err)
			}
		}()
	}
	wg.Wait()

	rd := routing.NewRoutingDiscovery(kadDHT)
	util.Advertise(ctx, rd, rendezvous)

	peerCh, err := rd.FindPeers(ctx, rendezvous)
	if err != nil {
		return nil, err
	}

	ddht := &discoveryDHT{
		PeerCh: peerCh,
		host:   h,
	}
	return ddht, nil
}

func dhtHandler(stream network.Stream) {
	packetSize := make([]byte, 4)

	for {
		if _, err := stream.Read(packetSize); err != nil {
			fmt.Printf("DHT Error reading length from stream: %v\n", err)
			stream.Close()
			return
		}

		address := make([]byte, binary.BigEndian.Uint32(packetSize))
		if _, err := stream.Read(address); err != nil {
			fmt.Printf("DHT Error reading message from stream: %v\n", err)
			stream.Close()
			return
		}
		ip, _, err := net.ParseCIDR(string(address))
		if err != nil {
			stream.Close()
			return
		}

		peerTable[ip.String()] = stream.Conn().RemotePeer()
		fmt.Printf("An address is now joined to vpn-mesh by DHT: %s\n", address)
	}
}
