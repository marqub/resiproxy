package rest

import (
	"github.com/gorilla/mux"
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name       string
		wantRoutes []string
	}{
		{
			name:       "OK",
			wantRoutes: []string{"CreateProxy", "Default"},
		},
	}
	for _, tt := range tests {
		teardown := setupTestCase(t)
		defer teardown(t)
		t.Run(tt.name, func(t *testing.T) {
			got := NewRouter()
			foundRoutes := []string{}
			err := got.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
				foundRoutes = append(foundRoutes, route.GetName())
				return nil
			})
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tt.wantRoutes, foundRoutes) {
				t.Errorf("Logger() = %v, want %v", foundRoutes, tt.wantRoutes)

			}
		})
	}
}
