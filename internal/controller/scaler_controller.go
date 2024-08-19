/*
Copyright 2024.

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

package controller

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	appsv1 "k8s.io/api/apps/v1"     //
	"k8s.io/apimachinery/pkg/types" //

	scalerv1alpha1 "github.com/subir0071/scaler-operator/api/v1alpha1"
)

// ScalerReconciler reconciles a Scaler object
type ScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=scaler.startup.com,resources=scalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scaler.startup.com,resources=scalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scaler.startup.com,resources=scalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Starting controller", "Time ::", time.Now())
	// TODO(user): your logic here
	scaler := &scalerv1alpha1.Scaler{}
	err := r.Get(ctx, req.NamespacedName, scaler)
	if err != nil {
		logger.Error(err, "Error encountered while getting details of Scaler Custom Resource")
		return ctrl.Result{}, err
	}
	dep := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: scaler.Spec.Deployment, Namespace: scaler.Namespace}, dep)
	if err != nil {
		logger.Error(err, "Error encountered while getting details of Deployment Resource")
		return ctrl.Result{}, err
	}
	dep.Spec.Replicas = &scaler.Spec.Replica
	err = r.Update(ctx, dep)
	if err != nil {
		logger.Error(err, "Deployment could not be scaled via the Scaler Controller")
		return ctrl.Result{}, err
	}
	logger.Info("Ending controller", "Time ::", time.Now())
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalerv1alpha1.Scaler{}).
		Complete(r)
}
