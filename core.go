package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pkg/errors"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/songgao/packets/ethernet"
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

const MTU = 1420

func readPackets(ctx context.Context, h host.Host, ifce *water.Interface) {
	activeStreams := make(map[string]network.Stream)

	var packet = make([]byte, MTU)
	for {
		plen, err := ifce.Read(packet)
		if err != nil {
			log.Println(err)
			continue
		}

		dst := net.IPv4(packet[16], packet[17], packet[18], packet[19]).String()

		stream, ok := activeStreams[dst]
		if ok {
			err = binary.Write(stream, binary.LittleEndian, uint16(plen))
			if err == nil {
				_, err = stream.Write(packet[:plen])
				if err == nil {
					continue
				}
			}
			stream.Close()
			delete(activeStreams, dst)
		}

		if peer, ok := peerTable[dst]; ok {
			stream, err = h.NewStream(ctx, peer, Protocol)
			if err != nil {
				continue
			}
			err = binary.Write(stream, binary.LittleEndian, uint16(plen))
			if err != nil {
				stream.Close()
				continue
			}
			_, err = stream.Write(packet[:plen])
			if err != nil {
				stream.Close()
				continue
			}

			activeStreams[dst] = stream
		}
	}
}

func getFrame(ifce *water.Interface) (ethernet.Frame, error) {
	var frame ethernet.Frame
	frame.Resize(MTU)

	n, err := ifce.Read([]byte(frame))
	if err != nil {
		return frame, errors.Wrap(err, "could not read from interface")
	}

	frame = frame[:n]
	return frame, nil
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
