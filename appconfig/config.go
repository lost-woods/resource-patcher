package appconfig

import (
	"strings"

	"github.com/lost-woods/go-util/k8sutil"
	"github.com/lost-woods/go-util/osutil"
)

var (
	config *AppConfig
)

func GetConfig() *AppConfig {
	if config == nil {
		loadConfig()
	}

	return config
}

func loadConfig() {
	config = &AppConfig{
		Resource: ResourceConfig{
			Type:      k8sutil.ParseResourceType(osutil.GetEnvStrReq("RESOURCE_TYPE")),
			Namespace: strings.ToLower(osutil.GetEnvStr("RESOURCE_NAMESPACE")), // can be empty
			Name:      strings.ToLower(osutil.GetEnvStrReq("RESOURCE_NAME")),
			Patch:     osutil.GetEnvStrReq("RESOURCE_PATCH"),
		},
	}
}
