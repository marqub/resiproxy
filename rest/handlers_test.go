package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var recorder *httptest.ResponseRecorder

func setupTestCase(t *testing.T) func(t *testing.T) {
	recorder = httptest.NewRecorder()
	return func(t *testing.T) {
	}
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
	}
}
func TestCreateProxy(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)
	type args struct {
		w     *httptest.ResponseRecorder
		proxy *Proxy
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "ERROR - No InClusterConfig",
			args: args{
				proxy: &Proxy{
					Name:     "TEST",
					Listen:   ":9090",
					Upstream: "http://google.com",
					Enabled:  true,
					Toxics:   "",
				},
				w: recorder,
			},
			want: []byte(`{"code":500,"text":"unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.args.proxy)
			req, _ := http.NewRequest("POST", "/proxies", bytes.NewBuffer(body))
			CreateProxy(tt.args.w, req)
			if out := tt.args.w.Body.Bytes(); !reflect.DeepEqual(out, tt.want) {
				t.Errorf("CreateProxy() = %v, want %v", string(out), tt.want)
			}
		})
	}
}

func Test_returnError(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)
	type args struct {
		code    int
		message string
		w       *httptest.ResponseRecorder
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "OK",
			args: args{
				code:    500,
				message: "Nope!",
				w:       recorder,
			},
			want: []byte(`{"code":500,"text":"Nope!"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnError(tt.args.code, tt.args.message, tt.args.w)

			if out := tt.args.w.Body.Bytes(); !reflect.DeepEqual(out, tt.want) {
				t.Errorf("returnError() = %v, want %v", string(out), tt.want)
			}
		})
	}
}

func TestProxyRequest(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "OK",
			args: args{
				w: recorder,
				r: httptest.NewRequest("GET", "/", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ProxyRequest(tt.args.w, tt.args.r)
		})
	}
}

func Test_serveReverseProxy(t *testing.T) {
	teardown := setupTestCase(t)
	defer teardown(t)
	type args struct {
		url *url.URL
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serveReverseProxy(tt.args.url, tt.args.res, tt.args.req)
		})
	}
}
