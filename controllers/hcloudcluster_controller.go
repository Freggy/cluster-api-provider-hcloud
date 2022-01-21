/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	infrav1 "github.com/Freggy/cluster-api-provider-hcloud/api/v1beta1"
	"github.com/Freggy/cluster-api-provider-hcloud/infra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// HCloudClusterReconciler reconciles a HCloudCluster object
type HCloudClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hcloudclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hcloudclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hcloudclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HCloudCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *HCloudClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ret ctrl.Result, reterr error) {
	logger := log.FromContext(ctx)

	var hcCluster infrav1.HCloudCluster
	if err := r.Get(ctx, req.NamespacedName, &hcCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	capiCluster, err := util.GetOwnerCluster(ctx, r.Client, hcCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}

	if capiCluster == nil {
		logger.Info("cluster object currently has no OwnerRef")
		return ctrl.Result{}, nil
	}

	clustersvc := infra.NewClusterService(capiCluster.GenerateName, nil)

	if hcCluster.Status.LoadBalancer.ID == 0 {
		lb, err := clustersvc.CreateLoadBalancer(ctx, hcCluster.Spec.LoadBalancer)
		if err != nil {
			return ctrl.Result{}, err
		}
		hcCluster.Status.LoadBalancer = infrav1.LoadBalancer{
			ID:        lb.ID,
			Algorithm: string(lb.Algorithm.Type),
			Location:  lb.Location.Name,
			Type:      lb.LoadBalancerType.Name,
			IPv4:      lb.PublicNet.IPv4.IP,
			IPv6:      lb.PublicNet.IPv6.IP,
			Port:      lb.Services[0].ListenPort,
		}
		hcCluster.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
			AdvertiseAddress: lb.PublicNet.IPv4.IP.String(),
			BindPort:         int32(lb.Services[0].ListenPort),
		}
	}

	defer func() {
		if err := r.Update(ctx, &hcCluster); err != nil {
			reterr = err
		}
	}()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HCloudClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.HCloudCluster{}).
		Complete(r)
}
