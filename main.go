package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

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

	ipAddr, prefix, octets, err := parseArgs(args)
	if err != nil {
		log.Fatalln(err)
	}

	c := Cidr{
		IPAddr:         ipAddr,
		Prefix:         prefix,
		HostIdentifier: 32 - prefix,
		ipAddrOctets:   octets,
	}

	c.subnetMaskOctets = c.GetSubnet()
	c.SubnetMask = stringify(c.subnetMaskOctets)
	c.NetworkAddr = stringify(c.GetAddr("network"))
	c.BroadcastAddr = stringify(c.GetAddr("broadcast"))
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
