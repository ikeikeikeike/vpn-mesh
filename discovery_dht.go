package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

type discoveryDHT struct {
	host       host.Host
	dht        *dht.IpfsDHT
	rendezvous string
}

func (n *discoveryDHT) run(address string) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()

		rd := routing.NewRoutingDiscovery(n.dht)
		util.Advertise(ctx, rd, n.rendezvous)

		peerCh, err := rd.FindPeers(ctx, n.rendezvous)
		if err != nil {
			fmt.Println("DHT FindPeers failed:", err)
			continue
		}

		for p := range peerCh {
			if p.ID == n.host.ID() || len(p.Addrs) == 0 {
				continue
			}

			switch n.host.Network().Connectedness(p.ID) {
			case network.NotConnected:
				if err := n.host.Connect(ctx, p); err != nil {
					// fmt.Println("DHT Connection failed:", p.ID, ">>", err)
					continue
				}
				fmt.Printf("An address is now joined to vpn-mesh by DHT: %s\n", p.ID)
			default:
				if err := discoveryWriter(ctx, n.host, address, p); err != nil {
					// fmt.Println("DHT writer failed:", p.ID, ">>", err)
					continue
				}
			}
		}
	}
}

func newDHT(ctx context.Context, h host.Host, rendezvous string) (*discoveryDHT, error) {
	kadDHT, err := dht.New(ctx, h)
	if err != nil {
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

	if err = kadDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}
	// cid, err := cid.NewPrefixV1(cid.Raw, mh.IDENTITY).Sum([]byte(rendezvous))
	// if err != nil {
	// 	return nil, err
	// }
	// if err := kadDHT.Provide(ctx, cid, true); err != nil {
	// 	return nil, err
	// }

	ddht := &discoveryDHT{
		host:       h,
		dht:        kadDHT,
		rendezvous: rendezvous,
	}
	return ddht, nil
}
