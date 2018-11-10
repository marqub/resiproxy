package main

import (
	"net/http"

	"github.com/marqub/resiproxy/log"
	"github.com/marqub/resiproxy/rest"
)

func main() {
	log.Logger().Info("Server started")
	router := rest.NewRouter()
	log.Logger().Fatal(http.ListenAndServe(":8080", router))
	//":"+os.Getenv("PORT")
}
