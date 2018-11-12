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
		panic(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := json.Unmarshal(body, &proxy); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusPreconditionFailed, Text: "Invalid Proxy"}); err != nil {
			panic(err)
		}
		return
	}
	log.Logger().Info("Create proxy :", string(body))
	err = k8s.CreateK8sMapping(proxy.Listen)
	if err != nil {
		log.Logger().Info("K8s mappings can not be created: ", err)
		return
	}
	ProxyRequest(w, r)
}

// ProxyRequest to toxiproxy service
func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	toxiproxyURL := url.URL{Scheme: k8s.Config.Scheme, Host: fmt.Sprintf("%s.%s:%d", k8s.Config.Name, k8s.Config.Namespace, k8s.Config.Port)}
	log.Logger().Info("Proxy request to: ", toxiproxyURL.String())
	serveReverseProxy(toxiproxyURL.String(), w, r)
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

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
