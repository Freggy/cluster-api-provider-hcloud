package v1beta1

import "net"

type LoadBalancer struct {
	ID        int
	Algorithm string
	Location  string
	Type      string
	IPv4      net.IP
	IPv6      net.IP
	Port      int
}
