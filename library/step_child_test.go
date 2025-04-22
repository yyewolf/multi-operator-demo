package library_test

import (
	"fmt"
	"library"
	"testing"

	appv1 "multi.ch/app/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestInstanceReflect(t *testing.T) {
	var example *appv1.App = &appv1.App{}
	fmt.Printf("example: %T, %p\n", example, example)

	app := library.NewInstanceOf(example) // Retourne un *appv1.App
	fmt.Printf("app: %T, %p\n", app, app)

	if example == app {
		t.Errorf("example and app should not be the same")
	}

	var instanceInterface client.Object = app

	appI := library.NewInstanceOf(instanceInterface) // Retourne un client.Object
	fmt.Printf("appI: %T, %p\n", appI, appI)

	if app == appI {
		t.Errorf("app and appI should not be the same")
	}
}
