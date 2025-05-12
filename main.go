package main

import (
	"github.com/lost-woods/go-util/k8sutil"
	"github.com/lost-woods/go-util/osutil"
	"github.com/lost-woods/resource-patcher/appconfig"
	"github.com/lost-woods/resource-patcher/lib"
)

var (
	log    = osutil.GetLogger()
	config = appconfig.GetConfig()
)

func main() {
	// Define the Kubernetes client
	client := k8sutil.CreateKubernetesClient(k8sutil.CreateKubeconfig())

	// Watch and Reconcile
	log.Infof("Watching %s '%s'...", config.Resource.Type, config.Resource.Name)
	informer := k8sutil.WatchResource(client, config.Resource.Type, config.Resource.Namespace, config.Resource.Name)
	lib.AddEventHandlers(informer, client)

	watcher := k8sutil.GenericWatcher{Name: "Resource Watcher"}
	watcher.Start(informer)

	osutil.WaitForCtrlC()
	log.Info("Shutting down...")
	watcher.Stop()
}
