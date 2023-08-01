package main

import (
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	tcp "github.com/libp2p/go-libp2p/p2p/transport/tcp"
)

func newP2P(port int) (host.Host, error) {
	pkey, err := privKey()
	if err != nil {
		return nil, err
	}

	addrs := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
		fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", port),
		fmt.Sprintf("/ip6/::/tcp/%d", port),
		fmt.Sprintf("/ip6/::/udp/%d/quic", port),
	)

	// Default Behavior: https://pkg.go.dev/github.com/libp2p/go-libp2p#New
	return libp2p.New(
		addrs,
		libp2p.Identity(pkey),
		// libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.DefaultMuxers,
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.FallbackDefaults,
	)
}

func privKey() (crypto.PrivKey, error) {
	// name := ".pkey"

	// Restore pkey
	// if _, err := os.Stat(name); !os.IsNotExist(err) {
	// 	dat, err := ioutil.ReadFile(name)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	return crypto.UnmarshalPrivateKey(dat)
	// }

	pkey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}
	// Store Key
	// privBytes, err := crypto.MarshalPrivateKey(pkey)
	// if err != nil {
	// 	return nil, err
	// }
	// if err := ioutil.WriteFile(name, privBytes, 0644); err != nil {
	// 	return nil, err
	// }

	return pkey, nil
}
