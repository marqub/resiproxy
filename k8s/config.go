package k8s

// ConfigToxiProxy defines the env variables passed by the configmap to access the ToxiProxy service
type ConfigToxiProxy struct {
	Name      string `required:"true"`
	Port      int    `required:"true"`
	Namespace string `required:"true"`
	Scheme    string `required:"true"`
}

// Config contains the info to access the ToxiProxy service
var Config ConfigToxiProxy
