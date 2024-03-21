package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.POST("/session", rt.wrap(rt.handleLogin))
	rt.router.PATCH("/users/:username", rt.wrap(rt.setMyUsername))
	rt.router.GET("/users/:username", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:username/mystream/", rt.wrap(rt.getStream))
	rt.router.POST("/users/:username/photos/", rt.wrap(rt.uplaodPhoto))
	rt.router.DELETE("/users/:username/photos/:photoid", rt.wrap(rt.deletePhoto))
	rt.router.POST("/users/:username/photos/:photoid/likes/", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:username/photos/:photoid/likes/:likingusername", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/users/:username/photos/:photoid/comments/", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:username/photos/:photoid/comments/:commentid", rt.wrap(rt.uncommentPhoto))
	rt.router.POST("/users/:username/following/", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/following/:followingusername", rt.wrap(rt.unfollowUser))
	rt.router.POST("/users/:username/banned", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/banned/:bannedusername", rt.wrap(rt.unbanUser))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
