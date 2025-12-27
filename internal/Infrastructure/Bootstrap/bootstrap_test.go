package bootstrap

import (
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBootstrapGin(t *testing.T) {
	container := dependency_injection.GetContainer()
	r := BootstrapGin(*container)
	assert.NotNil(t, r)

	expectedRoutes := [][]string{
		{"GET", "/api/v1/posts"},
		{"GET", "/api/v1/posts/:id"},
		{"PUT", "/api/v1/posts/:id"},
		{"POST", "/api/v1/posts"},
		{"DELETE", "/api/v1/posts/:id"},
		{"GET", "/auth/:provider/callback"},
		{"GET", "/auth/:provider"},
		{"GET", "/auth/logout/:provider"},
	}
	for _, route := range r.Routes() {
		found := false
		for _, expectedRoute := range expectedRoutes {
			if route.Method == expectedRoute[0] && strings.HasPrefix(route.Path, expectedRoute[1]) {
				found = true
				break
			}
		}
		assert.True(t, found, "Route %s %s not found", route.Method, route.Path)
	}
}
