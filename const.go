package main

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"

	"github.com/songgao/water"
)

const MTU = 1420

const (
	MeshProtocol     Protocol = "/mesh/0.0.1"
	DiscoverProtocol Protocol = "/mesh/discover/0.0.1"
)

type Protocol string

func (p Protocol) ID() protocol.ID {
	return protocol.ID(string(p))
}

const (
	port1 = 6868      // for dev
	name1 = "utun5"   // for dev
	type1 = water.TUN // for dev

	rendezvous = "mememe"
)

var (
	peerTable map[string]peer.ID
)
