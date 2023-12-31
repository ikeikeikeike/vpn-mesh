package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/pkg/errors"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

func meshHandler(i *water.Interface) func(stream network.Stream) {
	return func(stream network.Stream) {
		// TODO
		// if _, ok := RevLookup[stream.Conn().RemotePeer()]; !ok {
		// 	stream.Reset()
		// 	return
		// }
		defer stream.Close()

		var packet = make([]byte, MTU)
		var packetSize = make([]byte, 2)
		for {
			_, err := stream.Read(packetSize)
			if err != nil {
				return
			}

			size := binary.LittleEndian.Uint16(packetSize)

			var plen uint16 = 0
			for plen < size {
				tmp, err := stream.Read(packet[plen:size])
				plen += uint16(tmp)
				if err != nil {
					return
				}
			}

			i.Write(packet[:size])
		}
	}
}

func meshBridge(ctx context.Context, h host.Host, i *water.Interface) {
	activeStreams := make(map[string]network.Stream)

	var packet = make([]byte, MTU)
	for {
		plen, err := i.Read(packet)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dst := net.IPv4(packet[16], packet[17], packet[18], packet[19]).String() // TODO: use ethernet.Frame

		if stream, ok := activeStreams[dst]; ok {
			if err := binary.Write(stream, binary.LittleEndian, uint16(plen)); err == nil {
				if _, err := stream.Write(packet[:plen]); err == nil {
					continue
				}
			}
			stream.Close()
			delete(activeStreams, dst)
		}

		if peer, ok := PeerTable[dst]; ok {
			stream, err := h.NewStream(ctx, peer, MeshProtocol.ID())
			if err != nil {
				continue
			}
			if err := binary.Write(stream, binary.LittleEndian, uint16(plen)); err != nil {
				stream.Close()
				continue
			}
			if _, err := stream.Write(packet[:plen]); err != nil {
				stream.Close()
				continue
			}

			activeStreams[dst] = stream
		}
	}
}

func getFrame(i *water.Interface) (ethernet.Frame, error) {
	var frame ethernet.Frame
	frame.Resize(MTU)

	n, err := i.Read([]byte(frame))
	if err != nil {
		return frame, errors.Wrap(err, "could not read from interface")
	}

	frame = frame[:n]
	return frame, nil
}
