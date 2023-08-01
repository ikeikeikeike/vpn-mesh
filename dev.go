package main

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

const (
	MeshProtocol      Protocol = "/mesh/0.0.1"
	DiscoveryProtocol Protocol = "/mesh/discovery/0.0.1"
)

type Protocol string

func (p Protocol) ID() protocol.ID {
	return protocol.ID(string(p))
}

const (
	MTU = 1420
)

var (
	PeerTable = map[string]peer.ID{}
)
