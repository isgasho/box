// Copyright 2019 Nordix foundation

// Borrowed from https://github.com/Nordix/simple-ipam

// Package ipam is a very simple IPAM it administers a single CIDR range, e.g "1000::/124".
//
// The functions are NOT thread safe.
package network

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/mikioh/ipaddr"
)

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

// IPAM holds the ipam state
type IPAM struct {
	fn        string
	cidr      *net.IPNet
	cursor    *ipaddr.Cursor
	allocated map[uint32]bool
}

// New creates a new IPAM for the passed CIDR.
// Error if the passed CIDR is invalid.
func NewIPAM(fn, cidr string) (*IPAM, error) {
	cursor, err := ipaddr.Parse(cidr)
	if err != nil {
		return nil, err
	}

	_, ipnet, _ := net.ParseCIDR(cidr)
	ipam := &IPAM{
		allocated: make(map[uint32]bool),
		cursor:    cursor,
		cidr:      ipnet,
		fn:        fn,
	}

	ipam.load()

	return ipam, nil
}

func (i *IPAM) load() error {
	data, err := ioutil.ReadFile(i.fn)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, i)
}

func (i *IPAM) save() error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(i.fn, data, 0600)
}

func (i *IPAM) MarshalJSON() ([]byte, error) {
	data := struct {
		CIDR      string
		Allocated map[uint32]bool
	}{
		CIDR:      i.cidr.String(),
		Allocated: make(map[uint32]bool),
	}

	for ip, inuse := range i.allocated {
		data.Allocated[ip] = inuse
	}

	return json.Marshal(data)
}

func (i *IPAM) UnmarshalJSON(bs []byte) error {
	data := struct {
		CIDR      string
		Allocated map[uint32]bool
	}{
		Allocated: make(map[uint32]bool),
	}

	if err := json.Unmarshal(bs, &data); err != nil {
		return err
	}

	cursor, err := ipaddr.Parse(data.CIDR)
	if err != nil {
		return err
	}
	_, ipnet, _ := net.ParseCIDR(data.CIDR)

	i.cidr = ipnet
	i.cursor = cursor

	for ip, inuse := range data.Allocated {
		i.allocated[ip] = inuse
	}

	return nil
}

// Allocate allocates a new address.
// An error is returned if there is no addresses left.
func (i *IPAM) Allocate() (net.IP, error) {
	defer i.save()
	if i.Unallocated() < 1 {
		return nil, fmt.Errorf("No addresses left")
	}
	for {
		p := i.cursor.Next()
		ip := ip2int(p.IP)
		if _, ok := i.allocated[ip]; !ok {
			i.allocated[ip] = true
			return p.IP, nil
		}
	}
}

// Free frees an allocated address.
// To free a non-allocated address is a no-op.
func (i *IPAM) Free(a net.IP) {
	defer i.save()
	delete(i.allocated, ip2int(a))
}

// Unallocated returns the number of unallocated addresses.
func (i *IPAM) Unallocated() uint32 {
	o, _ := i.cidr.Mask.Size()
	return uint32(1) << uint(32-o)
}

// Reserve reserves an address.
// Error if the address is outside the CIDR or if the address is allocated already.
func (i *IPAM) Reserve(a net.IP) error {
	defer i.save()
	if !i.cidr.Contains(a) {
		return fmt.Errorf("Address outside the cidr")
	}
	ip := ip2int(a)
	if _, ok := i.allocated[ip]; ok {
		return fmt.Errorf("Address already allocated")
	}
	i.allocated[ip] = true
	return nil
}

// ReserveFirstAndLast reserves the first and last address.
// These are valid addresses but some programs may refuse to use them.
// Note that the number of Unallocated addresses may become zero.
func (i *IPAM) ReserveFirstAndLast() {
	i.Reserve(i.cursor.First().IP)
	i.Reserve(i.cursor.Last().IP)
}
