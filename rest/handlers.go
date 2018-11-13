package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// CreateProxy by delegating the call to the ToxiProxy service and try to open the k8s ports before
func CreateProxy(w http.ResponseWriter, r *http.Request) {

	var proxy Proxy
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Logger().Info("Request can't be read: ", err)
		returnError(http.StatusInternalServerError, "Invalid request", w)
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := json.Unmarshal(body, &proxy); err != nil {
		// unprocessable entity
		returnError(http.StatusPreconditionFailed, "Invalid Proxy", w)
		return
	}
	log.Logger().Info("Create proxy :", string(body))
	err = k8s.CreateK8sMapping(proxy.Listen)
	if err != nil {
		returnError(http.StatusInternalServerError, "Internal error", w)
		log.Logger().Info("K8s mappings can not be created: ", err)
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
