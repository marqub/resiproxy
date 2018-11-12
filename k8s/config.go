package k8s

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/marqub/resiproxy/log"
)

// ConfigToxiProxy defines the env variables passed by the configmap to access the ToxiProxy service
type ConfigToxiProxy struct {
	Name      string `required:"true"`
	Port      int    `required:"true"`
	Namespace string `required:"true"`
	Scheme    string `required:"true"`
}

// Config contains the info to access the ToxiProxy service
var Config ConfigToxiProxy

func init() {
	err := envconfig.Process("toxiproxy", &Config)
	if err != nil {
		log.Logger().Fatal(err.Error())
	}
}
