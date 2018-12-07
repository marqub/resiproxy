package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/marqub/resiproxy/k8s"
	"github.com/marqub/resiproxy/log"
)

type Proxy struct {
	Name     string `json:"name"`
	Listen   string `json:"listen"`
	Upstream string `json:"upstream"`
	Enabled  bool   `json:"enabled"`
	Toxics   string `json:"-"`
}
type jsonStatus struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

// CreateProxy by delegating the call to the ToxiProxy service and try to open the k8s ports before
func CreateProxy(w http.ResponseWriter, r *http.Request) {

	var proxy Proxy

	// Check that this is a POST first
	if r.Method != "POST" {
		log.Logger().Error("invalid method ")
		returnError(http.StatusMethodNotAllowed, "Invalid Method for this URL", w)
	}

	// Read the in POST body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Logger().Errorf("request can't be read: %v ", err)
		returnError(http.StatusInternalServerError, "Invalid request", w)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	// Unmarshall the data
	if err := json.Unmarshal(data, &proxy); err != nil {
		// unprocessable entity
		log.Logger().Errorf("unprocessable entity: %v", err)
		returnError(http.StatusPreconditionFailed, "Invalid Proxy", w)
		return
	}

	// Create the proxy
	log.Logger().Infof("Create proxy : %v ", string(data))
	err = k8s.CreateK8sMapping(proxy.Listen)
	if err != nil {
		log.Logger().Errorf("k8s mappings can not be created: %v", err)
		returnError(http.StatusInternalServerError, err.Error(), w)
		return
	}
	ProxyRequest(w, r)
}

func returnError(code int, message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(422)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: code, Text: message}); err != nil {
		log.Logger().Error(err)
	}
}

// ProxyRequest to toxiproxy service
func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	toxiproxyURL := url.URL{Scheme: k8s.Config.Scheme, Host: fmt.Sprintf("%s.%s:%d", k8s.Config.Name, k8s.Config.Namespace, k8s.Config.Port)}
	log.Logger().Info("Proxy request to: ", toxiproxyURL.String())
	serveReverseProxy(&toxiproxyURL, w, r)
}

// Serve a reverse proxy for a given url
func serveReverseProxy(url *url.URL, res http.ResponseWriter, req *http.Request) {
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Healthcheck the microservice
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(jsonStatus{Status: "OK", Name: "ResiProxy"}); err != nil {
		returnError(http.StatusInternalServerError, "Invalid response", w)
	}
}
