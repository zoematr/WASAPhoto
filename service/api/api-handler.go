package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.POST("/session", rt.wrap(rt.handleLogin))
	rt.router.GET("/users/:username/mystream/", rt.wrap(rt.GetStream))
	rt.router.PATCH("/users/:username", rt.wrap(rt.setMyUsername))
	rt.router.GET("/users/:username", rt.wrap(rt.getUserProfile))
	rt.router.POST("/users/:username/following/", rt.wrap(rt.followUser))
	rt.router.POST("/users/:username/following/:followingusername", rt.wrap(rt.unfollowUser))
	rt.router.POST("/users/:username/banned", rt.wrap(rt.banUser))
	rt.router.POST("/users/:username/banned/:bannedusername", rt.wrap(rt.unbanUser))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
