package rest

import (
	"github.com/gorilla/mux"
)

//NewRouter returns a given router
func NewRouter() *mux.Router {
	// Generate router
	router := mux.NewRouter().StrictSlash(true)
	// Handle CreateProxy
	router.Methods("POST").Path("/proxies").Name("CreateProxy").HandlerFunc(CreateProxy)
	//Health check
	router.Methods("GET").Path("/status").Name("Healthcheck").HandlerFunc(Healthcheck)
	// Proxy all other requests
	router.Methods([]string{"POST", "GET", "PATCH", "PUT", "DELETE"}...).Path("/{any:.+}").Name("Default").HandlerFunc(ProxyRequest)
	return router
}
