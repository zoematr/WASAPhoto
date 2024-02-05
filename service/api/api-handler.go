package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/session", rt.wrap(rt.handleLogin))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
