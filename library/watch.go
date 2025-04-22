package library

import (
	"context"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func SetupWatch[
	ControllerResourceType ControllerResource,
](
	reconciler Reconciler[ControllerResourceType],
	object client.Object,
) func(ctx context.Context, req ctrl.Request) StepResult {
	return func(ctx context.Context, req ctrl.Request) StepResult {
		// Setup watch if not already set
		watchSource := NewWatchKey(object, CacheTypeEnqueueForOwner)
		if !reconciler.IsWatchingSource(watchSource) {
			// Add the watch source to the reconciler
			err := reconciler.GetController().Watch(
				source.Kind(
					reconciler.GetCache(),
					object,
					handler.EnqueueRequestForOwner(reconciler.GetScheme(), reconciler.GetRESTMapper(), reconciler.GetCustomResource()),
				),
			)
			if err != nil {
				return ResultInError(errors.Wrap(err, "failed to add watch source"))
			}

			reconciler.AddWatchSource(watchSource)
		}

		return ResultSuccess()
	}
}
