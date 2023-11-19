package testcontroller

import (
	"os"

	"github.com/jasonkeene/test-after-deploy/pkg/reconcilers"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func Run() {
	log.SetLogger(zap.New())

	scheme := runtime.NewScheme()
	err := clientgoscheme.AddToScheme(scheme)
	if err != nil {
		log.Log.Error(err, "can't add scheme")
		return
	}

	namespace := os.Getenv("NAMESPACE")

	mgr, err := ctrl.NewManager(
		ctrl.GetConfigOrDie(),
		ctrl.Options{
			Scheme: scheme,
			// This cache config prevents the controller from making API
			// requests for resources in other namespaces. The filtering in
			// SetupWithManager on namespace is just a safeguard after the API
			// request was already made. By preventing the API requests from
			// being made we limit the permissions scope of the controller.
			Cache: cache.Options{
				DefaultNamespaces: map[string]cache.Config{
					namespace: {},
				},
			},
		},
	)
	if err != nil {
		log.Log.Error(err, "can't create manager")
		return
	}

	r := reconcilers.NewDeployment(
		mgr.GetClient(),
		"server",
		namespace,
		"tests",
		"tests",
	)
	err = r.SetupWithManager(mgr)
	if err != nil {
		log.Log.Error(err, "can't setup reconciler")
		return
	}

	err = mgr.Start(ctrl.SetupSignalHandler())
	if err != nil {
		log.Log.Error(err, "manager exited with error")
		return
	}
}
