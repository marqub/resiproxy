package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/marqub/resiproxy/k8s"
	"github.com/marqub/resiproxy/log"
	"github.com/marqub/resiproxy/rest"
	"net/http"
)

func init() {
	err := envconfig.Process("toxiproxy", &k8s.Config)
	if err != nil {
		log.Logger().Fatal(err.Error())
	}
}
func main() {
	router := rest.NewRouter()
	log.Logger().Info("Server started")
	log.Logger().Info("Configuration: ", k8s.Config)
	log.Logger().Fatal(http.ListenAndServe(":8080", router))
}
