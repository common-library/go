package utility

import "net"

func WhetherCidrContainsIp(cidr, ip string) (bool, error) {
	if _, subnet, err := net.ParseCIDR(cidr); err != nil {
		return false, err
	} else {
		return subnet.Contains(net.ParseIP(ip)), nil
	}
}

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
