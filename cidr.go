package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Cidr struct {
	IPAddr           string
	Prefix           int
	IPAddrOctets     []int
	SubnetMaskOctets []int
	TotalHosts       int
}

func (c Cidr) GetAddr(name string) []int {
	addr := GetSlice()
	i := 0
	for i < len(c.IPAddrOctets) {
		if name == "broadcast" {
			addr[i] = c.IPAddrOctets[i] | (255 - c.SubnetMaskOctets[i])
		} else if name == "network" {
			addr[i] = c.IPAddrOctets[i] & c.SubnetMaskOctets[i]
		}
		i++
	}
	return addr
}

func (c Cidr) GetTotalHosts() int {
	hostBits := 32 - c.Prefix
	// Subtract both the network and broadcast addresses.
	return PowerOf2(hostBits) - 2
}

func (c Cidr) String() string {
	return fmt.Sprintf(
		"       IP address: %s\n   Network prefix: %d bits\n      Subnet mask: %s\n  Network address: %s\nBroadcast address: %s\n Total # of hosts: %d\n",
		c.IPAddr,
		c.Prefix,
		Stringify(c.SubnetMaskOctets),
		Stringify(c.GetAddr("network")),
		Stringify(c.GetAddr("broadcast")),
		c.GetTotalHosts(),
	)
}

func Atoi(s string) int {
	d, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return d
}

func Btoi(n int) int {
	// Initialize with the value of the left-most bit of the octet.
	// If this function is reached, `n` will be at least 1.
	i := 128
	octet := 0
	for n > 0 {
		octet += i
		i /= 2
		n--
	}
	return octet
}

func GetSlice() []int {
	return []int{0, 0, 0, 0}
}

func GetSubnet(prefix int) []int {
	subnet := GetSlice()
	var bits int
	i := 0
	p := prefix
	for p > 0 {
		if p > 8 {
			bits = 8
			p -= 8
		} else {
			bits = p
			p = 0
		}
		subnet[i] = Btoi(bits)
		i++
	}
	return subnet
}

func ParseArgs() (string, int, []int, error) {
	parts := strings.Split(os.Args[1], "/")
	prefix := Atoi(parts[1])
	if prefix > 32 || prefix < 0 {
		err := fmt.Errorf("%s [ERROR] Prefix cannot be greater than 32 or less than zero.\n", os.Args[0])
		return "", 0, []int{}, err
	}
	i := 0
	ip := strings.SplitN(parts[0], ".", 4)
	if len(ip) != 4 {
		err := fmt.Errorf("%s [ERROR] This does not look like an IP address.  There must be four octets.\n", os.Args[0])
		return "", 0, []int{}, err
	}
	octets := make([]int, 4)
	for i < len(octets) {
		octet := Atoi(ip[i])
		if octet < 0 || octet > 255 {
			err := fmt.Errorf("%s [ERROR] Octet cannot be greater than 255 or less than zero.\n", os.Args[0])
			return "", 0, []int{}, err
		}
		octets[i] = octet
		i++
	}
	return parts[0], prefix, octets, nil
}

func PowerOf2(n int) int {
	if n == 0 {
		return 1
	}
	return PowerOf2(n-1) * 2
}

func Stringify(addr []int) string {
	parts := make([]string, 4)
	for i := 0; i < len(addr); i++ {
		parts[i] = strconv.Itoa(addr[i])
	}
	return strings.Join(parts, ".")
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("%s [ERROR] Not enough arguments.\n", os.Args[0])
		os.Exit(1)
	}

	if !strings.Contains(os.Args[1], "/") {
		fmt.Printf("%s [ERROR] Please provide a network prefex (i.e., /24).\n", os.Args[0])
		os.Exit(1)
	}

	ipAddr, prefix, octets, err := ParseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(Cidr{
		IPAddr:           ipAddr,
		Prefix:           prefix,
		IPAddrOctets:     octets,
		SubnetMaskOctets: GetSubnet(prefix),
	})
}
