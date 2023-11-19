package reconcilers

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/polymorphichelpers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type Deployment struct {
	client       client.Client
	statusViewer polymorphichelpers.DeploymentStatusViewer
	name         string
	namespace    string
	jobName      string
	jobImage     string
}

func NewDeployment(
	client client.Client,
	name string,
	namespace string,
	jobName string,
	jobImage string,
) *Deployment {
	return &Deployment{
		client:    client,
		name:      name,
		namespace: namespace,
		jobName:   jobName,
		jobImage:  jobImage,
	}
}

func (r *Deployment) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var dep appsv1.Deployment
	err := r.client.Get(ctx, req.NamespacedName, &dep)
	if err != nil {
		return ctrl.Result{}, err
	}

	um, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&dep)
	if err != nil {
		return ctrl.Result{}, err
	}
	var ud unstructured.Unstructured
	ud.SetUnstructuredContent(um)

	_, ok, err := r.statusViewer.Status(&ud, 0)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !ok {
		return ctrl.Result{}, nil
	}
	jobName := fmt.Sprintf("%s-%d", r.jobName, dep.Status.ObservedGeneration)
	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: r.namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      jobName,
					Namespace: r.namespace,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            r.jobName,
							Image:           r.jobImage,
							ImagePullPolicy: v1.PullIfNotPresent,
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			TTLSecondsAfterFinished: int32Ptr(60 * 60),
		},
	}

	log.Log.Info("Attempting to create job", "jobName", jobName)
	err = r.client.Create(context.Background(), &job)
	if err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *Deployment) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(predicate.NewPredicateFuncs(func(o client.Object) bool {
			return o.GetNamespace() == r.namespace &&
				o.GetName() == r.name
		})).
		Complete(r)
}

func int32Ptr(i int32) *int32 {
	return &i
}
