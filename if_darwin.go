//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

func newif(name string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name

	return water.New(config)
}

func applyif(name, address string) error {
	ip, _, err := net.ParseCIDR(address)
	if err != nil {
		return err
	}
	if err := ifconfig(name, "mtu", fmt.Sprintf("%d", MTU)); err != nil {
		return err
	}
	if err := ifconfig(name, "inet", ip.String(), "10.1.1.1", "up"); err != nil { // XXX: 10.1.1.1
		return err
	}

	return nil
}

func ifconfig(args ...string) error {
	cmd := exec.Command("ifconfig", args...)
	return cmd.Run()
}
