package infra

import (
	"context"
	"fmt"

	infrav1 "github.com/Freggy/cluster-api-provider-hcloud/api/v1beta1"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type ClusterService struct {
	cluster string
	client  *hcloud.Client
}

func NewClusterService(cluster string, client *hcloud.Client) *ClusterService {
	return &ClusterService{
		cluster: cluster,
		client:  client,
	}
}

func (cs *ClusterService) GetLoadBalancer(ctx context.Context, infraLB infrav1.LoadBalancer) (*hcloud.LoadBalancer, error) {
	lb, _, err := cs.client.LoadBalancer.GetByID(ctx, infraLB.ID)
	if err != nil {
		return nil, err
	}
	return lb, nil
}

func (cs *ClusterService) CreateLoadBalancer(ctx context.Context, lb infrav1.LoadBalancer) (*hcloud.LoadBalancer, error) {
	applyLBDefaults(&lb)
	opts := hcloud.LoadBalancerCreateOpts{
		Name: fmt.Sprintf("lb-%s", cs.cluster),
		LoadBalancerType: &hcloud.LoadBalancerType{
			Name: lb.Type,
		},
		Algorithm: &hcloud.LoadBalancerAlgorithm{
			Type: hcloud.LoadBalancerAlgorithmType(lb.Algorithm),
		},
		Location: &hcloud.Location{
			Name: lb.Location,
		},
		Labels: map[string]string{
			"cluster": cs.cluster,
		},
	}
	result, _, err := cs.client.LoadBalancer.Create(ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.LoadBalancer, nil
}

func applyLBDefaults(lb *infrav1.LoadBalancer) {
	if lb.Type == "" {
		lb.Type = ""
	}

	if lb.Location == "" {
		lb.Location = ""
	}

	if lb.Algorithm == "" {
		lb.Algorithm = "round_robin"
	}
}
