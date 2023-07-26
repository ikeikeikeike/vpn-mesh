//go:build !darwin
// +build !darwin

package main

import (
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

func createInterface(name string, dtype water.DeviceType) (*water.Interface, error) {
	config := water.Config{
		DeviceType: dtype,
		// PlatformSpecificParams: water.PlatformSpecificParams{Persist: !c.NetLinkBootstrap},
	}
	config.Name = name

	return water.New(config)
}

func prepareInterface(name, address string, mtu int) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}
	if err := netlink.LinkSetMTU(link, mtu); err != nil {
		return err
	}
	if err := netlink.AddrAdd(link, addr); err != nil {
		return err
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	return nil
}
