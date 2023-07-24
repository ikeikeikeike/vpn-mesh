package main

import "context"

func main() {
	h, err := genHost(port1)
	if err != nil {
		panic(err)
	}
	defer h.Close()

	ifce, err := createInterface(name1, type1)
	if err != nil {
		panic(err)
	}
	defer ifce.Close()

	h.SetStreamHandler(Protocol, streamHandler(ifce))

	readPackets(context.Background(), h, ifce)
}
