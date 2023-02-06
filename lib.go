package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func atoi(s string) (int, error) {
	d, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return d, nil
}

func btoi(n int) int {
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

func getSlice() []int {
	return []int{0, 0, 0, 0}
}

func parseArgs(args []string) (string, int, []int, error) {
	parts := strings.Split(args[0], "/")
	prefix, err := atoi(parts[1])
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
		octet, err := atoi(ip[i])
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

func powerOf2(n int) int {
	if n == 0 {
		return 1
	}
	return powerOf2(n-1) * 2
}

func stringify(addr []int) string {
	parts := make([]string, 4)
	for i := 0; i < len(addr); i++ {
		parts[i] = strconv.Itoa(addr[i])
	}
	return strings.Join(parts, ".")
}
