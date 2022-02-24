package utils

import "net"

// Reserved IP blocks (source: https://wikipedia.org/wiki/Reserved_IP_Addresses)
var reservedIPs = []string{
	"0.0.0.0/8",          // Current network, only valid as source
	"10.0.0.0/8",         // Used for local communications within a private network
	"100.64.0.0/10",      // Shared address space for use with carrier-grade NAT
	"127.0.0.0/8",        // Used for loopback addresses to the local host
	"169.154.0.0/16",     // Used for link-local addresses when no IP is otherwise specificed
	"172.16.0.0/12",      // Used for local communications within a private network
	"192.0.0.0/24",       // IETF Protocol Assignments
	"192.0.2.0/24",       // Assigned as TEST-NET-1, documentation and examples
	"192.88.99.0/24",     // Reserved. Formerly used for IPv6 to IPv4 relay
	"192.168.0.0/16",     // Used for local communications within a private network
	"198.18.0.0/15",      // Used for benchmark testing of inter-network communications
	"198.51.100.0/24",    // Assigned as TEST-NET-2, documentation and examples
	"203.0.113.0/24",     // Assigned as TEST-NET-3, documentation and examples
	"224.0.0.0/4",        // In use for IP multicast. (Former class D network)
	"240.0.0.0/4",        // Reserved for future use. (Former class E network)
	"255.255.255.255/32", // Reserved for the "limited broadcast" destination address
}

var reservedBlocks []*net.IPNet

func init() {
	for _, b := range reservedIPs {
		_, block, err := net.ParseCIDR(b)
		if err != nil {
			panic(err)
		}
		reservedBlocks = append(reservedBlocks, block)
	}
}

// IsReserved 是否为保留IP
func IsReserved(ip net.IP) bool {
	for _, block := range reservedBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// IsPublicIP 是否为公网IP
func IsPublicIP(IP string) bool {
	ip := net.ParseIP(IP)
	return !IsReserved(ip)
}

// IsIP 是否为IP地址
func IsIP(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil
}
