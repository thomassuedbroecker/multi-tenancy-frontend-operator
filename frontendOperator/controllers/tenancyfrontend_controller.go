/*
Copyright 2022.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	// Add to read error from Kubernetes
	"k8s.io/apimachinery/pkg/api/errors"

	multitenancyv1alpha1 "github.com/thomassuedbroecker/multi-tenancy-frontend-operator/api/v1alpha1"
)

// TenancyFrontendReconciler reconciles a TenancyFrontend object
type TenancyFrontendReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=multitenancy.example.net,resources=tenancyfrontends,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multitenancy.example.net,resources=tenancyfrontends/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=multitenancy.example.net,resources=tenancyfrontends/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TenancyFrontend object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *TenancyFrontendReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// "Verify if a CRD of TenancyFrontend exists"
	logger.Info("Verify if a CRD of TenancyFrontend exists")
	tenancyfrontend := &multitenancyv1alpha1.TenancyFrontend{}
	err := r.Get(ctx, req.NamespacedName, tenancyfrontend)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("TenancyFrontend resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get TenancyFrontend")
		return ctrl.Result{}, err
	} else {
		return ctrl.Result{}, nil
	}

}

// SetupWithManager sets up the controller with the Manager.
func (r *TenancyFrontendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multitenancyv1alpha1.TenancyFrontend{}).
		Complete(r)
}
