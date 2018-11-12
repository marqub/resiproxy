package main

import (
	"net/http"

	"github.com/marqub/resiproxy/k8s"
	"github.com/marqub/resiproxy/log"
	"github.com/marqub/resiproxy/rest"
)

func main() {
	router := rest.NewRouter()
	log.Logger().Info("Server started")
	log.Logger().Info("Configuration: ", k8s.Config)
	log.Logger().Fatal(http.ListenAndServe(":8080", router))
}
