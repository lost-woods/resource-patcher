package appconfig

import "k8s.io/apimachinery/pkg/runtime/schema"

type AppConfig struct {
	Resource ResourceConfig
}

type ResourceConfig struct {
	Type      schema.GroupVersionResource
	Namespace string
	Name      string
	Patch     string
}
