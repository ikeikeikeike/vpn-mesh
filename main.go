package main

import (
	"context"
	"crypto/rand"
	"flag"
	"io"
	"log"
	"os"
)

const Protocol = "/dmsg/0.0.1"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")
	help := flag.Bool("help", false, "Display help")
	flag.Parse()
	if *help {
		os.Exit(0)
	}

	var r io.Reader = rand.Reader
	h, err := newHost(*sourcePort, r)
	if err != nil {
		log.Println(err)
		return
	}

	if *dest == "" {
		startPeer(ctx, h, handleStream)
	} else {
		rw, err := startPeerAndConnect(ctx, h, *dest)
		if err != nil {
			log.Println(err)
			return
		}

		go writeIO(rw)
		go readIO(rw)
	}

	select {}
}
