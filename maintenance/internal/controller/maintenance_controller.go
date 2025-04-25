/*
Copyright 2025.

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
	"library"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	envoyapiv1alpha1 "github.com/envoyproxy/gateway/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	maintenancev1 "multi.ch/maintenance/api/v1"
	routev1 "multi.ch/route/api/v1"
)

// MaintenanceReconciler reconciles a Maintenance object
type MaintenanceReconciler struct {
	ctrl.Manager
	client.Client
	library.WatchCache

	RuntimeScheme *runtime.Scheme
	controller    controller.TypedController[reconcile.Request]

	maintenance maintenancev1.Maintenance

	backend envoyapiv1alpha1.Backend
}

var _ library.Reconciler[*maintenancev1.Maintenance] = &MaintenanceReconciler{}

func (reconciler *MaintenanceReconciler) GetController() controller.TypedController[reconcile.Request] {
	return reconciler.controller
}

func (reconciler *MaintenanceReconciler) GetFinalizer() string {
	return "maintenance.multi.ch/finalizer"
}

func (reconciler *MaintenanceReconciler) GetCustomResource() *maintenancev1.Maintenance {
	return &reconciler.maintenance
}

func (reconciler *MaintenanceReconciler) SetCustomResource(maintenance *maintenancev1.Maintenance) {
	reconciler.maintenance = *maintenance
}

func (reconciler *MaintenanceReconciler) GetChildren(ctx context.Context, req ctrl.Request) ([]library.GenericChildResource, error) {
	return []library.GenericChildResource{
		library.NewChildResource(
			&envoyapiv1alpha1.Backend{},
			library.WithChildOutput(&reconciler.backend),
			library.WithChildGenerator(reconciler.backendGenerator),
		),
	}, nil
}

// +kubebuilder:rbac:groups=maintenance.multi.ch,resources=maintenances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=maintenance.multi.ch,resources=maintenances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=maintenance.multi.ch,resources=maintenances/finalizers,verbs=update

// +kubebuilder:rbac:groups=gateway.envoyproxy.io,resources=backends,verbs=get;list;watch;create;update;patch;delete

func (reconciler *MaintenanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	stepper := library.NewStepper(logger,
		library.WithStep(library.NewFindControllerResourceStep(reconciler)),
		library.WithStep(library.NewReconcileChildrenStep(reconciler)),
		library.WithStep(reconciler.NewFillContractStep()),
		library.WithStep(library.NewEndStep(reconciler)),
	)

	return stepper.Execute(ctx, req)
}

// SetupWithManager sets up the controller with the Manager.
func (reconciler *MaintenanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	reconciler.Manager = mgr

	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&maintenancev1.Maintenance{}).
		Named("maintenance").
		Build(reconciler)
	if err != nil {
		return err
	}

	reconciler.controller = controller

	return err
}

func (reconciler *MaintenanceReconciler) backendGenerator(ctx context.Context, req ctrl.Request) (*envoyapiv1alpha1.Backend, bool, error) {
	return &envoyapiv1alpha1.Backend{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reconciler.maintenance.Name,
			Namespace: reconciler.maintenance.Namespace,
		},
		Spec: envoyapiv1alpha1.BackendSpec{
			Endpoints: []envoyapiv1alpha1.BackendEndpoint{
				{
					FQDN: &envoyapiv1alpha1.FQDNEndpoint{
						Hostname: "maintenance.maintenance.svc",
						Port:     80,
					},
				},
			},
		},
	}, false, nil
}

func (reconciler *MaintenanceReconciler) NewFillContractStep() library.Step {
	return library.Step{
		Name: "Fill Contract",
		Step: func(ctx context.Context, req ctrl.Request) library.StepResult {
			newContract := routev1.RouteContract{
				BackendRef: &routev1.RouteContractLocalBackendRef{
					Name: reconciler.maintenance.Name,
					Port: 80,
				},
			}

			if reflect.DeepEqual(reconciler.maintenance.Status.RouteContractInjector.RouteContract, newContract) {
				return library.ResultSuccess()
			}

			reconciler.maintenance.Status.RouteContractInjector.RouteContract = newContract

			if err := reconciler.Status().Update(ctx, &reconciler.maintenance); err != nil {
				return library.ResultInError(err)
			}

			return library.ResultSuccess()
		},
	}
}
