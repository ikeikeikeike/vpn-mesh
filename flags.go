package main

import (
	"flag"

	"github.com/pkg/errors"
)

type args struct {
	Rendezvous string
	Network    string
	Interface  string
	Port       int
}

func parseArgs() (*args, error) {
	a := &args{}

	flag.StringVar(&a.Rendezvous, "rv", "", "Rendezvous string like the only master key")
	flag.StringVar(&a.Network, "net", "", "vpn-mesh host e.g. 10.1.1.1/24, 10.1.1.2/24")
	flag.StringVar(&a.Interface, "ifce", "utun5", "vpn-mesh VPN(TUN) interface name")
	flag.IntVar(&a.Port, "port", 6868, "vpn-mesh port\n")

	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, r := range []string{"rv", "net"} {
		if !seen[r] {
			return nil, errors.Errorf("missing required -%s argument/flag\n", r)
		}
	}

	return a, nil
}
