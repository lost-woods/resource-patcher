package lib

import (
	"github.com/lost-woods/go-util/k8sutil"
	"github.com/lost-woods/go-util/osutil"
	"github.com/lost-woods/go-util/treeutil"
	"github.com/lost-woods/resource-patcher/appconfig"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
)

var (
	log    = osutil.GetLogger()
	config = appconfig.GetConfig()
)

func patchResource(client *dynamic.DynamicClient) {
	log.Infof("Patching %s '%s' with the following data:\n%s", config.Resource.Type, config.Resource.Name, config.Resource.Patch)
	k8sutil.PatchResource(client, config.Resource.Type, config.Resource.Namespace, config.Resource.Name, []byte(config.Resource.Patch))
}

func AddEventHandlers(informer cache.SharedIndexInformer, client *dynamic.DynamicClient) {
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Infof("Resource %s '%s' created or appeared, reconciling...", config.Resource.Type.Resource, config.Resource.Name)
			patchResource(client)
		},

		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			_newObj, okNew := newObj.(*unstructured.Unstructured)
			if !okNew {
				log.Fatalf("Failed to cast '%s' into a kubernetes object", config.Resource.Type.Resource)
			}

			// Check patches
			shouldPatch := false
			jsonPatches := k8sutil.ParseJsonPatch(config.Resource.Patch)
			for _, patch := range jsonPatches {
				val, ok := treeutil.GetValueAtPath(_newObj.UnstructuredContent(), patch.Path)
				if !ok || !treeutil.DeepEqualNormalized(val, patch.Value) {
					shouldPatch = true
				}
			}

			if shouldPatch {
				log.Infof("Resource %s '%s' changed, reconciling...", config.Resource.Type.Resource, config.Resource.Name)
				patchResource(client)
			}
		},
	})
}
