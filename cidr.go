package main

import "fmt"

type Cidr struct {
	IPAddr           string `yaml:"ip_addr,omitempty" json:"ip_addr,omitempty"`
	BroadcastAddr    string `yaml:"broadcast_addr,omitempty" json:"broadcast_addr,omitempty"`
	NetworkAddr      string `yaml:"network_addr,omitempty" json:"network_addr,omitempty"`
	SubnetMask       string `yaml:"subnet_mask,omitempty" json:"subnet_mask,omitempty"`
	Prefix           int    `yaml:"prefix,omitempty" json:"prefix,omitempty"`
	HostIdentifier   int    `yaml:"host_identifier,omitempty" json:"host_identifier,omitempty"`
	TotalHosts       int    `yaml:"total_hosts,omitempty" json:"total_hosts,omitempty"`
	ipAddrOctets     []int
	subnetMaskOctets []int
}

func (c Cidr) GetAddr(name string) []int {
	addr := getSlice()
	i := 0
	for i < len(c.ipAddrOctets) {
		if name == "broadcast" {
			addr[i] = c.ipAddrOctets[i] | (255 - c.subnetMaskOctets[i])
		} else if name == "network" {
			addr[i] = c.ipAddrOctets[i] & c.subnetMaskOctets[i]
		}
		i++
	}
	return addr
}

func (c Cidr) GetSubnet() []int {
	subnet := getSlice()
	var bits int
	i := 0
	p := c.Prefix
	for p > 0 {
		if p > 8 {
			bits = 8
			p -= 8
		} else {
			bits = p
			p = 0
		}
		subnet[i] = btoi(bits)
		i++
	}
	return subnet
}

func (c Cidr) GetTotalHosts() int {
	hostBits := 32 - c.Prefix
	// Subtract both the network and broadcast addresses.
	return powerOf2(hostBits) - 2
}

func (c Cidr) String() string {
	return fmt.Sprintf(
		`   Network prefix: %d
  Host identifier: %d
      Total hosts: %d
       IP address: %s
      Subnet mask: %s
  Network address: %s
Broadcast address: %s
`,
		c.Prefix,
		c.HostIdentifier,
		c.TotalHosts,
		c.IPAddr,
		c.SubnetMask,
		c.NetworkAddr,
		c.BroadcastAddr,
	)
}
