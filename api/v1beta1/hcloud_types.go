package v1beta1

import "net"

type LoadBalancer struct {
	// ID of the load balancer, generated by the cloud provider
	// +optional
	ID int `json:"ID,omitempty"`

	// Algorithm used
	// +optional
	Algorithm string `json:"algorithm,omitempty"`

	// Type of the load balancer
	// +optional
	Type string `json:"type,omitempty"`

	// Public IPv4 address
	// +optional
	IPv4 net.IP `json:"ipv4,omitempty"`

	// Public IPv6 address
	// +optional
	IPv6 net.IP `json:"ipv6,omitempty"`

	// Port where the load balancer will bind to
	// +optional
	Port int `json:"port,omitempty"`
}
