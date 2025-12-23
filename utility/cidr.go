package utility

import "net"

// WhetherCidrContainsIp checks if an IP address is within a CIDR range.
//
// This function parses a CIDR notation string and determines whether the
// specified IP address falls within that network range.
//
// Parameters:
//   - cidr: CIDR notation string (e.g., "192.168.1.0/24", "10.0.0.0/8")
//   - ip: IP address string (e.g., "192.168.1.100", "10.0.0.1")
//
// Returns:
//   - bool: true if IP is in CIDR range, false otherwise
//   - error: Error if CIDR parsing fails
//
// Example:
//
//	contains, err := utility.WhetherCidrContainsIp("192.168.1.0/24", "192.168.1.100")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if contains {
//	    fmt.Println("IP is in range")
//	} else {
//	    fmt.Println("IP is not in range")
//	}
//
// Example with error handling:
//
//	cidr := "10.0.0.0/8"
//	ip := "10.5.100.200"
//
//	contains, err := utility.WhetherCidrContainsIp(cidr, ip)
//	if err != nil {
//	    fmt.Printf("Invalid CIDR: %v\n", err)
//	    return
//	}
//
//	fmt.Printf("%s is in %s: %v\n", ip, cidr, contains)
func WhetherCidrContainsIp(cidr, ip string) (bool, error) {
	if _, subnet, err := net.ParseCIDR(cidr); err != nil {
		return false, err
	} else {
		return subnet.Contains(net.ParseIP(ip)), nil
	}
}

// GetAllIpsOfCidr returns all usable IP addresses in a CIDR range.
//
// This function generates a list of all IP addresses within the specified
// CIDR notation, excluding the network and broadcast addresses.
//
// Parameters:
//   - cidr: CIDR notation string (e.g., "192.168.1.0/24")
//
// Returns:
//   - []string: Slice of IP address strings (excluding network and broadcast IPs)
//   - error: Error if CIDR parsing fails
//
// For example, "192.168.1.0/30" contains 4 addresses:
//   - 192.168.1.0 (network address, excluded)
//   - 192.168.1.1 (returned)
//   - 192.168.1.2 (returned)
//   - 192.168.1.3 (broadcast address, excluded)
//
// Warning: Large CIDR ranges (e.g., /8) will generate millions of IPs
// and may consume significant memory. Use with caution.
//
// Example with small range:
//
//	ips, err := utility.GetAllIpsOfCidr("192.168.1.0/30")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Found %d usable IPs:\n", len(ips))
//	for _, ip := range ips {
//	    fmt.Println(ip)
//	}
//	// Output:
//	// Found 2 usable IPs:
//	// 192.168.1.1
//	// 192.168.1.2
//
// Example with /24 subnet:
//
//	ips, err := utility.GetAllIpsOfCidr("10.0.1.0/24")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Subnet has %d usable IPs\n", len(ips))
//	// Output: Subnet has 254 usable IPs
//
// Example for IP allocation:
//
//	cidr := "172.16.0.0/28"
//	allIps, _ := utility.GetAllIpsOfCidr(cidr)
//
//	// Allocate first 5 IPs
//	allocated := allIps[:5]
//	available := allIps[5:]
//
//	fmt.Printf("Allocated: %v\n", allocated)
//	fmt.Printf("Available: %d IPs\n", len(available))
func GetAllIpsOfCidr(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); {
		ips = append(ips, ip.String())

		for i := len(ip) - 1; i >= 0; i-- {
			ip[i]++
			if ip[i] > 0 {
				break
			}
		}
	}

	switch {
	case len(ips) < 2:
		return ips, nil
	default:
		return ips[1 : len(ips)-1], nil
	}
}
