package main

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/songgao/water"

	"github.com/multiformats/go-multiaddr"
)

func genHost(port int) (host.Host, error) {
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}
	maddrs, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	if err != nil {
		return nil, err
	}

	return libp2p.New(
		libp2p.ListenAddrs(maddrs),
		libp2p.Identity(prvKey),
		// libp2p.EnableAutoRelay(),
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),
	)
}

func streamHandler(ifce *water.Interface) func(stream network.Stream) {
	return func(stream network.Stream) {
		if len(peerTable) == 0 {
			stream.Reset()
			return
		}

		if len(peerTable) > 0 {
			found := false
			for _, p := range peerTable {
				if p.String() == stream.Conn().RemotePeer().String() {
					found = true
				}
			}
			if !found {
				stream.Reset()
				return
			}
		}

		if _, err := io.Copy(ifce.ReadWriteCloser, stream); err != nil {
			stream.Reset()
		}
	}
}

const Protocol = "/dmsg/0.0.1"

const (
	port1 = 6868      // for dev
	name1 = "utun5"   // for dev
	type1 = water.TUN // for dev
)

var (
	peerTable map[string]peer.ID
)
