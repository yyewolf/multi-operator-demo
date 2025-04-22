package controller

import (
	"context"
	"fmt"
	"library"

	"golang.org/x/exp/maps"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	routev1 "multi.ch/route/api/v1"
)

// RouteReconciler reconciles a Route object
type RouteReconciler struct {
	ctrl.Manager
	client.Client
	library.WatchCache

	RuntimeScheme *runtime.Scheme
	controller    controller.TypedController[reconcile.Request]

	// CR
	route routev1.Route

	// Dependencies
	targets map[routev1.RouteTargetReference]*unstructured.Unstructured

	// Children
	httproute gatewayv1.HTTPRoute
}

var _ library.Reconciler[*routev1.Route] = &RouteReconciler{}

func (reconciler *RouteReconciler) GetController() controller.TypedController[reconcile.Request] {
	return reconciler.controller
}

func (reconciler *RouteReconciler) GetFinalizer() string {
	return "route.multi.ch/finalizer"
}

func (reconciler *RouteReconciler) GetCustomResource() *routev1.Route {
	return &reconciler.route
}

func (reconciler *RouteReconciler) SetCustomResource(route *routev1.Route) {
	reconciler.route = *route
}

func (reconciler *RouteReconciler) GetDependencies(ctx context.Context, req ctrl.Request) (dependencies []library.GenericDependencyResource, err error) {
	reconciler.targets = make(map[routev1.RouteTargetReference]*unstructured.Unstructured)

	for _, target := range reconciler.route.Spec.TargetRefs {
		var output unstructured.Unstructured
		reconciler.targets[*target] = &output

		gv, err := schema.ParseGroupVersion(target.APIVersion)
		if err != nil {
			return nil, err
		}

		gvk := schema.GroupVersionKind{
			Group:   gv.Group,
			Version: gv.Version,
			Kind:    target.Kind,
		}

		dependency := library.NewUntypedDependencyResource(
			gvk,
			library.WithName[*unstructured.Unstructured](target.Name),
			library.WithNamespace[*unstructured.Unstructured](reconciler.route.Namespace),
			library.WithOutput(&output),
		)

		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

func (reconciler *RouteReconciler) GetChildren(ctx context.Context, req ctrl.Request) ([]library.GenericChildResource, error) {
	return []library.GenericChildResource{
		library.NewChildResource(
			&gatewayv1.HTTPRoute{},
			library.WithChildOutput(&reconciler.httproute),
			library.WithChildGenerator(reconciler.httpRouteGenerator),
		),
	}, nil
}

// +kubebuilder:rbac:groups=route.multi.ch,resources=routes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=route.multi.ch,resources=routes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=route.multi.ch,resources=routes/finalizers,verbs=update

// +kubebuilder:rbac:groups=app.multi.ch,resources=apps,verbs=get;list;watch;update;patch

func (reconciler *RouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	stepper := library.NewStepper(logger,
		library.WithStep(library.NewFindControllerResourceStep(reconciler)),
		library.WithStep(library.NewResolveDynamicDependenciesStep(reconciler)),
		// library.WithStep(reconciler.NewTestStep()),
		library.WithStep(library.NewReconcileChildrenStep(reconciler)),
		library.WithStep(library.NewEndStep(reconciler)),
	)

	return stepper.Execute(ctx, req)
}

func (reconciler *RouteReconciler) httpRouteGenerator(ctx context.Context, req ctrl.Request) (*gatewayv1.HTTPRoute, bool, error) {
	var hostnames []gatewayv1.Hostname
	for _, hostname := range reconciler.route.Spec.Hostnames {
		hostnames = append(hostnames, gatewayv1.Hostname(hostname))
	}

	var rules []gatewayv1.HTTPRouteRule
	for i, key := range maps.Keys(reconciler.targets) {
		target := reconciler.targets[key]

		// Get the route contract from the target
		routeContract, err := library.GetContract[routev1.RouteContract](target, "routeContract")
		if err != nil {
			return nil, false, err
		}

		rules = append(rules, gatewayv1.HTTPRouteRule{
			Name: library.Opt(gatewayv1.SectionName(fmt.Sprintf("target.%d", i))),
			Matches: []gatewayv1.HTTPRouteMatch{
				{
					Path: &gatewayv1.HTTPPathMatch{
						Type:  library.Opt(gatewayv1.PathMatchPathPrefix),
						Value: library.Opt("/"),
					},
				},
			},
			BackendRefs: []gatewayv1.HTTPBackendRef{
				{
					BackendRef: gatewayv1.BackendRef{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name: gatewayv1.ObjectName(routeContract.ServiceRef.Name),
							Port: library.Opt(gatewayv1.PortNumber(routeContract.ServiceRef.Port)),
						},
					},
				},
			},
		})
	}

	return &gatewayv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reconciler.route.Name,
			Namespace: reconciler.route.Namespace,
		},
		Spec: gatewayv1.HTTPRouteSpec{
			Hostnames: hostnames,
			Rules:     rules,
		},
	}, false, nil
}

// func (reconciler *RouteReconciler) NewTestStep() library.Step {
// 	return library.Step{
// 		Name: "Test Step",
// 		Step: func(ctx context.Context, req ctrl.Request) library.StepResult {
// 			for _, target := range reconciler.targets {
// 				fmt.Println(target)
// 			}

// 			return library.ResultSuccess()
// 		},
// 	}
// }

// SetupWithManager sets up the controller with the Manager.
func (reconciler *RouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	reconciler.Manager = mgr

	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&routev1.Route{}).
		Named("route").
		Build(reconciler)
	if err != nil {
		return err
	}

	reconciler.controller = controller

	return nil
}
