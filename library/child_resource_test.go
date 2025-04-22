package library_test

import (
	"context"
	"library"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1 "multi.ch/app/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestChildResource(t *testing.T) {
	var original *appv1.App = &appv1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "original",
			Namespace: "default",
		},
	}

	var list []library.GenericChildResource = []library.GenericChildResource{
		library.NewChildResource(
			&appv1.App{},
			library.WithChildGenerator(func(ctx context.Context, req ctrl.Request) (*appv1.App, bool, error) {
				return original, false, nil
			}),
			library.WithChildStatusGetter(func(a *appv1.App) *library.Status {
				return a.GetStatus()
			}),
			library.WithChildOutput(original),
		),
	}

	child := list[0]

	var replaced *appv1.App = &appv1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "replaced",
			Namespace: "default",
		},
	}

	child.Set(replaced)

	if original.Name != replaced.Name {
		t.Errorf("original and replaced should be the same")
	}
}
