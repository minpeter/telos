package api

import (
	"github.com/gin-gonic/gin"
	"github.com/minpeter/rctf-backend/api/admin"
	"github.com/minpeter/rctf-backend/api/auth"
	"github.com/minpeter/rctf-backend/api/challs"
	"github.com/minpeter/rctf-backend/api/integrations"
	"github.com/minpeter/rctf-backend/api/leaderboard"
	"github.com/minpeter/rctf-backend/api/users"
)

func NewRouter() *gin.Engine {
	app := gin.Default()

	router := app.Group("/api")

	// Admin-related routes
	admin.Routes(router.Group("/admin"))

	// Authentication-related routes
	auth.Routes(router.Group("/auth"))

	// Challenge-specific routes
	challs.Routes(router.Group("/challs"))

	// Integrations-related routes
	integrations.Routes(router.Group("/integrations"))

	// Leaderboard-related routes
	leaderboard.Routes(router.Group("/leaderboard"))

	// User-related routes
	users.Routes(router.Group("/users"))

	return app
}