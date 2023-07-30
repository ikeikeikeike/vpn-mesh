package main

import (
	"github.com/libp2p/go-libp2p-core/protocol"

	"github.com/songgao/water"
)

const (
	MeshProtocol Protocol = "/mesh/0.0.1"
	MDNSProtocol Protocol = "/mesh/mdns/0.0.1"
	DHTProtocol  Protocol = "/mesh/dht/0.0.1"
)

type Protocol string

func (p Protocol) ID() protocol.ID {
	return protocol.ID(string(p))
}

const (
	MTU   = 1420
	port1 = 6868      // for dev
	name1 = "utun5"   // for dev
	type1 = water.TUN // for dev
)
