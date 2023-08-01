//go:build !darwin
// +build !darwin

package main

import (
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

func newif(name string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
		// PlatformSpecificParams: water.PlatformSpecificParams{Persist: !c.NetLinkBootstrap},
	}
	config.Name = name

	return water.New(config)
}

func applyif(name, address string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}
	if err := netlink.LinkSetMTU(link, MTU); err != nil {
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
