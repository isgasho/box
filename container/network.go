package container

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"

	"github.com/prologic/box/filesystem"
	"github.com/prologic/box/network"
)

const (
	netnsPath       = "/var/lib/box/netns"
	ipamFile        = "/var/lib/box/ipam.json"
	bridgeInterface = "box0"
)

var ipam *network.IPAM

func init() {
	var err error
	ipam, err = network.NewIPAM(ipamFile, "172.30.0.0/16")
	if err != nil {
		panic(fmt.Errorf("error creating IPAM driver: %w", err))
	}

	ipam.ReserveFirstAndLast()
	ipam.Reserve(net.ParseIP("172.30.0.1"))
}

func (c *Container) SetupNetwork(bridge string) (filesystem.Unmounter, error) {
	nsMountTarget := filepath.Join(netnsPath, c.Digest)
	vethName := fmt.Sprintf("veth%.7s", c.Digest)
	peerName := fmt.Sprintf("P%s", vethName)
	masterName := bridgeInterface

	if err := network.SetupVirtualEthernet(vethName, peerName); err != nil {
		return nil, err
	}
	if err := network.LinkSetMaster(vethName, masterName); err != nil {
		return nil, err
	}
	unmount, err := network.MountNewNetworkNamespace(nsMountTarget)
	if err != nil {
		return unmount, err
	}
	if err := network.LinkSetNsByFile(nsMountTarget, peerName); err != nil {
		return unmount, err
	}

	// Change current network namespace to setup the veth
	unset, err := network.SetNetNSByFile(nsMountTarget)
	if err != nil {
		return unmount, nil
	}
	defer unset()

	ctrEthName := "eth0"

	ctrEthIPAddr, release, err := c.GetIP()
	if err != nil {
		return func() error { unmount(); release(); return nil }, err
	}

	if err := network.LinkRename(peerName, ctrEthName); err != nil {
		return func() error { unmount(); release(); return nil }, err
	}
	if err := network.LinkAddAddr(ctrEthName, ctrEthIPAddr); err != nil {
		return func() error { unmount(); release(); return nil }, err
	}
	if err := network.LinkSetup(ctrEthName); err != nil {
		return func() error { unmount(); release(); return nil }, err
	}
	if err := network.LinkAddGateway(ctrEthName, "172.30.0.1"); err != nil {
		return func() error { unmount(); release(); return nil }, err
	}
	if err := network.LinkSetup("lo"); err != nil {
		return func() error { unmount(); release(); return nil }, err
	}

	return func() error { unmount(); release(); return nil }, nil
}

func (c *Container) SetNetworkNamespace() (network.Unsetter, error) {
	netns := filepath.Join(netnsPath, c.Digest)
	return network.SetNetNSByFile(netns)
}

func (c *Container) GetIP() (string, network.Releaser, error) {
	genLinkLocalAddress := func() string {
		a, _ := strconv.ParseInt(c.Digest[:2], 10, 64)
		b, _ := strconv.ParseInt(c.Digest[62:], 10, 64)
		return fmt.Sprintf("169.254.%d.%d/16", a, b)
	}

	ip, err := ipam.Allocate()
	if err != nil {
		return genLinkLocalAddress(), func() error { return nil }, err
	}

	return fmt.Sprintf("%s/16", ip), func() error { ipam.Free(ip); return nil }, nil
}
