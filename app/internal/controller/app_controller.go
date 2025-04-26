package controller

import (
	"bytes"
	"context"
	"library"
	"reflect"
	"text/template"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1 "multi.ch/app/api/v1"
	"multi.ch/app/internal/supervisord"
	routev1 "multi.ch/route/api/v1"
)

const (
	Selector = "app.multi.ch/app-name"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	ctrl.Manager
	client.Client
	library.WatchCache

	RuntimeScheme *runtime.Scheme
	controller    controller.TypedController[reconcile.Request]

	app appv1.App

	// Children
	configMap  corev1.ConfigMap
	deployment appsv1.Deployment
	service    corev1.Service
}

var _ library.Reconciler[*appv1.App] = &AppReconciler{}

func (reconciler *AppReconciler) GetController() controller.TypedController[reconcile.Request] {
	return reconciler.controller
}

func (reconciler *AppReconciler) GetFinalizer() string {
	return "app.multi.ch/finalizer"
}

func (reconciler *AppReconciler) GetCustomResource() *appv1.App {
	return &reconciler.app
}

func (reconciler *AppReconciler) SetCustomResource(app *appv1.App) {
	reconciler.app = *app
}

func (reconciler *AppReconciler) GetDependencies(ctx context.Context, req ctrl.Request) ([]library.GenericDependencyResource, error) {
	return []library.GenericDependencyResource{
		library.NewDependencyResource(
			&corev1.ConfigMap{},
			library.WithName[*corev1.ConfigMap]("test"),
			library.WithNamespace[*corev1.ConfigMap]("default"),
		),
	}, nil
}

func (reconciler *AppReconciler) GetChildren(ctx context.Context, req ctrl.Request) ([]library.GenericChildResource, error) {
	return []library.GenericChildResource{
		library.NewChildResource(
			&corev1.ConfigMap{},
			library.WithChildOutput(&reconciler.configMap),
			library.WithChildGenerator(reconciler.configMapGenerator),
		),
		library.NewChildResource(
			&appsv1.Deployment{},
			library.WithChildOutput(&reconciler.deployment),
			library.WithChildGenerator(reconciler.deploymentGenerator),
		),
		library.NewChildResource(
			&corev1.Service{},
			library.WithChildOutput(&reconciler.service),
			library.WithChildGenerator(reconciler.serviceGenerator),
		),
	}, nil
}

// +kubebuilder:rbac:groups=app.multi.ch,resources=apps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.multi.ch,resources=apps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.multi.ch,resources=apps/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (reconciler *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	stepper := library.NewStepper(logger,
		library.WithStep(library.NewFindControllerResourceStep(reconciler)),
		library.WithStep(library.NewResolveDynamicDependenciesStep(reconciler)),
		library.WithStep(library.NewReconcileChildrenStep(reconciler)),
		library.WithStep(reconciler.NewFillContractStep()),
		library.WithStep(library.NewEndStep(reconciler)),
	)

	return stepper.Execute(ctx, req)
}

// SetupWithManager sets up the controller with the Manager.
func (reconciler *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	reconciler.Manager = mgr

	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&appv1.App{}).
		Named("app").
		Build(reconciler)
	if err != nil {
		return err
	}

	reconciler.controller = controller

	return err
}

type WorkloadConfigurationTemplateData struct {
	Command string
}

func (reconciler *AppReconciler) configMapGenerator(ctx context.Context, req ctrl.Request) (*corev1.ConfigMap, bool, error) {
	supervisordConfiguration, err := supervisord.GetStaticFile("supervisord.conf")
	if err != nil {
		return nil, false, err
	}

	initConfiguration, err := supervisord.GetStaticFile("init.conf")
	if err != nil {
		return nil, false, err
	}

	workloadConfiguration, err := supervisord.GetStaticFile("workload.conf")
	if err != nil {
		return nil, false, err
	}

	template, err := template.New("workload").Parse(string(workloadConfiguration))
	if err != nil {
		return nil, false, err
	}

	workloadConfigurationData := WorkloadConfigurationTemplateData{
		Command: reconciler.app.Spec.Command,
	}

	var output bytes.Buffer
	err = template.Execute(&output, workloadConfigurationData)
	if err != nil {
		return nil, false, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reconciler.app.Name,
			Namespace: req.Namespace,
		},
		Data: map[string]string{
			"supervisord.conf": string(supervisordConfiguration),
			"workload.conf":    output.String(),
			"init.conf":        string(initConfiguration),
		},
	}, false, nil
}

func (reconciler *AppReconciler) serviceGenerator(ctx context.Context, req ctrl.Request) (*corev1.Service, bool, error) {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reconciler.app.Name,
			Namespace: req.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				Selector: reconciler.app.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "agent",
					Protocol:   corev1.ProtocolTCP,
					Port:       1080,
					TargetPort: intstr.FromInt(1080),
				},
				{
					Name:       "workload",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(int(reconciler.app.Spec.Port)),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}, false, nil
}

func (reconciler *AppReconciler) deploymentGenerator(ctx context.Context, req ctrl.Request) (*appsv1.Deployment, bool, error) {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reconciler.app.Name,
			Namespace: req.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					Selector: reconciler.app.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						Selector: reconciler.app.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "supervisord",
							Image: "workload",
							Command: []string{
								"/usr/local/bin/supervisord",
								"-c",
								"/etc/supervisor.d/supervisord.conf",
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "supervisor-config",
									MountPath: "/etc/supervisor.d",
								},
								{
									Name:      "logs",
									MountPath: "/var/log/supervisor",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "supervisor-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: reconciler.app.Name,
									},
									DefaultMode: library.Opt(int32(0555)),
								},
							},
						},
						{
							Name: "logs",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}, false, nil
}

func (reconciler *AppReconciler) NewFillContractStep() library.Step {
	return library.Step{
		Name: "Fill Contract",
		Step: func(ctx context.Context, req ctrl.Request) library.StepResult {
			newContract := routev1.RouteContract{
				ServiceRef: &routev1.RouteContractLocalServiceRef{
					Name: reconciler.service.Name,
					Port: 80,
				},
			}

			if reflect.DeepEqual(reconciler.app.Status.RouteContractInjector.RouteContract, newContract) {
				return library.ResultSuccess()
			}

			reconciler.app.Status.RouteContractInjector.RouteContract = newContract

			if err := reconciler.Status().Update(ctx, &reconciler.app); err != nil {
				return library.ResultInError(err)
			}

			return library.ResultSuccess()
		},
	}
}
