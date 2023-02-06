package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

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
	addr := GetSlice()
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

func (c Cidr) GetTotalHosts() int {
	hostBits := 32 - c.Prefix
	// Subtract both the network and broadcast addresses.
	return PowerOf2(hostBits) - 2
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

func Atoi(s string) (int, error) {
	d, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return d, nil
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

func ParseArgs(args []string) (string, int, []int, error) {
	parts := strings.Split(args[0], "/")
	prefix, err := Atoi(parts[1])
	if err != nil {
		err := fmt.Errorf("%s [ERROR] %s\n", err, os.Args[0])
		return "", 0, []int{}, err
	}

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
		octet, err := Atoi(ip[i])
		if err != nil {
			err := fmt.Errorf("%s [ERROR] %s\n", err, os.Args[0])
			return "", 0, []int{}, err
		}
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
	_json := flag.Bool("json", false, "json format")
	_yaml := flag.Bool("yaml", false, "yaml format")
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		log.Fatalf("%s [ERROR] Not enough arguments.\n", os.Args[0])
	}

	if !strings.Contains(args[0], "/") {
		log.Fatalf("%s [ERROR] Please provide a network prefix (i.e., /24).\n", os.Args[0])
	}

	ipAddr, prefix, octets, err := ParseArgs(args)
	if err != nil {
		log.Fatalln(err)
	}

	c := Cidr{
		IPAddr:           ipAddr,
		Prefix:           prefix,
		HostIdentifier:   32 - prefix,
		ipAddrOctets:     octets,
		subnetMaskOctets: GetSubnet(prefix),
	}

	c.SubnetMask = Stringify(c.subnetMaskOctets)
	c.NetworkAddr = Stringify(c.GetAddr("network"))
	c.BroadcastAddr = Stringify(c.GetAddr("broadcast"))
	c.TotalHosts = c.GetTotalHosts()

	if *_json {
		jsonEncoder := json.NewEncoder(os.Stdout)
		jsonEncoder.SetIndent("", "    ")
		if err := jsonEncoder.Encode(c); err != nil {
			log.Fatalf("%s [ERROR] %v", os.Args[0], err)
		}
	}

	if *_yaml {
		d, err := yaml.Marshal(&c)
		if err != nil {
			log.Fatalf("%s [ERROR] %v", os.Args[0], err)
		}
		fmt.Printf("--- cidr dump:\n%s", string(d))
	}

	if !*_yaml && !*_json {
		fmt.Println(c)
	}
}
