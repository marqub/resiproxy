package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"CreateProxy",
		[]string{"POST"},
		"/proxies",
		CreateProxy},
	Route{
		"Default",
		[]string{"POST", "GET", "PATCH", "PUT", "DELETE"},
		"/{any:.+}",
		ProxyRequest},
}
