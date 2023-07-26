//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

func createInterface(name string, dtype water.DeviceType) (*water.Interface, error) {
	config := water.Config{
		DeviceType: dtype,
	}
	config.Name = name

	return water.New(config)
}

func prepareInterface(name, address string, mtu int) error {
	ip, _, err := net.ParseCIDR(address)
	if err != nil {
		return err
	}
	if err := ifconfig(name, "mtu", fmt.Sprintf("%d", MTU)); err != nil {
		return err
	}
	if err := ifconfig(name, "inet", ip.String(), "10.1.1.1", "up"); err != nil {
		return err
	}

	return nil
}

func ifconfig(args ...string) error {
	cmd := exec.Command("ifconfig", args...)
	return cmd.Run()
}
