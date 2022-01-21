package infra

import (
	"context"
	"fmt"
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type ClusterService struct {
	cluster string
	client  *hcloud.Client
}

type LoadBalancer struct {
	ID        int
	Name      string
	Location  string
	Type      string
	Algorithm string
	IPV4      net.IP
	IPV6      net.IP
	Port      int
}

func NewClusterService(cluster string, client *hcloud.Client) *ClusterService {
	return &ClusterService{
		cluster: cluster,
		client:  client,
	}
}

func (cs *ClusterService) GetLoadBalancer(ctx context.Context, id int) (LoadBalancer, error) {
	lb, _, err := cs.client.LoadBalancer.GetByID(ctx, id)
	if err != nil {
		return LoadBalancer{}, err
	}
	return LoadBalancer{
		ID:        lb.ID,
		Name:      lb.Name,
		Location:  lb.Location.Name,
		Type:      lb.LoadBalancerType.Name,
		Algorithm: string(lb.Algorithm.Type),
		IPV4:      lb.PublicNet.IPv4.IP,
		IPV6:      lb.PublicNet.IPv6.IP,
		Port:      lb.Services[0].ListenPort,
	}, nil
}

func (cs *ClusterService) CreateLoadBalancer(ctx context.Context, lbType string, location string) (LoadBalancer, error) {
	opts := hcloud.LoadBalancerCreateOpts{
		Name: fmt.Sprintf("lb-%s", cs.cluster),
		LoadBalancerType: &hcloud.LoadBalancerType{
			Name: lbType,
		},
		Algorithm: &hcloud.LoadBalancerAlgorithm{
			Type: hcloud.LoadBalancerAlgorithmTypeRoundRobin,
		},
		Location: &hcloud.Location{
			Name: location,
		},
		Labels: map[string]string{
			"cluster": cs.cluster,
		},
	}
	result, _, err := cs.client.LoadBalancer.Create(ctx, opts)
	if err != nil {
		return LoadBalancer{}, err
	}
	return LoadBalancer{
		ID:        result.LoadBalancer.ID,
		Name:      result.LoadBalancer.Name,
		Location:  result.LoadBalancer.Location.Name,
		Type:      result.LoadBalancer.LoadBalancerType.Name,
		Algorithm: string(result.LoadBalancer.Algorithm.Type),
		IPV4:      result.LoadBalancer.PublicNet.IPv4.IP,
		IPV6:      result.LoadBalancer.PublicNet.IPv6.IP,
		Port:      result.LoadBalancer.Services[0].ListenPort,
	}, nil
}
