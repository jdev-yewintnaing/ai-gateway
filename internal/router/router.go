package router

import (
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/yewintnaing/ai-gateway/internal/config"
)

type Router struct {
	routes []config.Route
}

func NewRouter(routes []config.Route) *Router {
	return &Router{routes: routes}
}

func (r *Router) Route(useCase string) config.Route {
	for _, route := range r.routes {
		if route.Match.UseCase == useCase {
			return route
		}
	}

	for _, route := range r.routes {
		if route.Name == "default" {
			return route
		}
	}

	// Fallback to minimal default
	return config.Route{
		Name:    "default",
		Primary: config.Target{Provider: "openai", Model: "gpt-4o-mini"},
		Retries: 1,
	}
}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Network errors
	if _, ok := err.(net.Error); ok {
		return true
	}
	if os.IsTimeout(err) {
		return true
	}
	if err == syscall.ECONNREFUSED || err == syscall.ECONNRESET {
		return true
	}

	return false
}

func StatusCodeIsRetryable(code int) bool {
	return code == http.StatusTooManyRequests || (code >= 500 && code <= 599)
}
