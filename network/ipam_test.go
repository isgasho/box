package network

import (
	"net"
	"testing"
)

func allocate(t *testing.T, ipam *IPAM, expected string, left uint32) {
	a, err := ipam.Allocate()
	u := ipam.Unallocated()
	if u != left {
		t.Fatalf("Unallocated %d, expected %d", u, left)
	}
	if err != nil {
		if expected != "" {
			t.Fatalf("Unexpected error for %s", expected)
		}
		return
	}
	if !a.Equal(net.ParseIP(expected)) {
		t.Fatalf("Address %s, expected %s", a, expected)
	}
}
func free(t *testing.T, ipam *IPAM, addr string, left uint32) {
	ipam.Free(net.ParseIP(addr))
	u := ipam.Unallocated()
	if u != left {
		t.Fatalf("Unallocated %d, expected %d", u, left)
	}
}
func create(t *testing.T, cidr string, left uint32) *IPAM {
	ipam, err := NewIPAM("", cidr)
	if err != nil {
		t.Fatalf("Failed to create ipam %s", cidr)
	}
	i := ipam.Unallocated()
	if i != left {
		t.Fatalf("Unallocated %d, expected %d", i, left)
	}
	return ipam
}
func reserve(t *testing.T, ipam *IPAM, addr string, expectedErr bool, left uint32) {
	err := ipam.Reserve(net.ParseIP(addr))
	if err != nil {
		if !expectedErr {
			t.Fatalf("Unexpected error for %s", addr)
		}
	}
	u := ipam.Unallocated()
	if u != left {
		t.Fatalf("Unallocated %d, expected %d", u, left)
	}
}
func reserveFirstAndLast(t *testing.T, ipam *IPAM, left uint32) {
	ipam.ReserveFirstAndLast()
	u := ipam.Unallocated()
	if u != left {
		t.Fatalf("Unallocated %d, expected %d", u, left)
	}
}

func TestBasic(t *testing.T) {
	ipam, err := NewIPAM("", "malformed")
	if err == nil {
		t.Fatalf("Could create a malformed ipam")
	}

	ipam = create(t, "1000::/127", 2)
	allocate(t, ipam, "1000::", 1)
	allocate(t, ipam, "1000::1", 0)
	allocate(t, ipam, "", 0)
	free(t, ipam, "1000::", 1)
	free(t, ipam, "1000::", 1)
	free(t, ipam, "1000::", 1)
	allocate(t, ipam, "1000::", 0)

	ipam = create(t, "10.10.10.0/29", 8)
	allocate(t, ipam, "10.10.10.0", 7)
	allocate(t, ipam, "10.10.10.1", 6)
	allocate(t, ipam, "10.10.10.2", 5)
	allocate(t, ipam, "10.10.10.3", 4)
	allocate(t, ipam, "10.10.10.4", 3)
	allocate(t, ipam, "10.10.10.5", 2)
	allocate(t, ipam, "10.10.10.6", 1)
	allocate(t, ipam, "10.10.10.7", 0)

	free(t, ipam, "10.10.10.3", 1)
	free(t, ipam, "10.10.10.5", 2)

	allocate(t, ipam, "10.10.10.3", 1)

	free(t, ipam, "10.10.10.0", 2)
	free(t, ipam, "10.10.10.1", 3)
	free(t, ipam, "10.10.10.2", 4)
	free(t, ipam, "10.10.10.3", 5)
	allocate(t, ipam, "10.10.10.5", 4)

	ipam = create(t, "100.10.1.0/29", 8)
	reserveFirstAndLast(t, ipam, 6)
	allocate(t, ipam, "100.10.1.1", 5)
	free(t, ipam, "100.10.1.1", 6)
	allocate(t, ipam, "100.10.1.2", 5)
	allocate(t, ipam, "100.10.1.3", 4)
	allocate(t, ipam, "100.10.1.4", 3)
}
