package v1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	routev1 "multi.ch/route/api/v1"
)

// nolint:unused
// log is for logging in this package.
var routelog = logf.Log.WithName("route-resource")

// SetupRouteWebhookWithManager registers the webhook for Route in the manager.
func SetupRouteWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&routev1.Route{}).
		WithValidator(&RouteCustomValidator{}).
		WithDefaulter(&RouteCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-route-multi-ch-v1-route,mutating=true,failurePolicy=fail,sideEffects=None,groups=route.multi.ch,resources=routes,verbs=create;update,versions=v1,name=mroute-v1.kb.io,admissionReviewVersions=v1

// RouteCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Route when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type RouteCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &RouteCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Route.
func (d *RouteCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	route, ok := obj.(*routev1.Route)

	if !ok {
		return fmt.Errorf("expected an Route object but got %T", obj)
	}
	routelog.Info("Defaulting for Route", "name", route.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-route-multi-ch-v1-route,mutating=false,failurePolicy=fail,sideEffects=None,groups=route.multi.ch,resources=routes,verbs=create;update,versions=v1,name=vroute-v1.kb.io,admissionReviewVersions=v1

// RouteCustomValidator struct is responsible for validating the Route resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type RouteCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &RouteCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Route.
func (v *RouteCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	route, ok := obj.(*routev1.Route)
	if !ok {
		return nil, fmt.Errorf("expected a Route object but got %T", obj)
	}
	routelog.Info("Validation for Route upon creation", "name", route.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Route.
func (v *RouteCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	route, ok := newObj.(*routev1.Route)
	if !ok {
		return nil, fmt.Errorf("expected a Route object for the newObj but got %T", newObj)
	}
	routelog.Info("Validation for Route upon update", "name", route.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Route.
func (v *RouteCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	route, ok := obj.(*routev1.Route)
	if !ok {
		return nil, fmt.Errorf("expected a Route object but got %T", obj)
	}
	routelog.Info("Validation for Route upon deletion", "name", route.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
