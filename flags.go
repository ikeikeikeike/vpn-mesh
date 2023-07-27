package main

import (
	"flag"

	"github.com/pkg/errors"
)

type args struct {
	Token   string
	Network string
}

func parseArgs() (*args, error) {
	a := &args{}

	flag.StringVar(&a.Token, "token", "", "The only master token")
	flag.StringVar(&a.Network, "network", "", "vpn-mesh host e.g. 10.1.1.1/24,  10.1.1.2/24\n")
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, r := range []string{"token", "network"} {
		if !seen[r] {
			return nil, errors.Errorf("missing required -%s argument/flag\n", r)
		}
	}

	return a, nil
}
