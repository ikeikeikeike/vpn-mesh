package main

import (
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"

	"github.com/multiformats/go-multiaddr"
)

func newP2P(port int) (host.Host, error) {
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
