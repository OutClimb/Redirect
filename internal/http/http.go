package http

import (
	"os"
	"strings"

	"github.com/OutClimb/Redirect/internal/app"
	"github.com/gin-gonic/gin"
)

type httpLayer struct {
	engine *gin.Engine
	app    app.AppLayer
}

func New(appLayer app.AppLayer) *httpLayer {
	h := &httpLayer{
		engine: gin.New(),
		app:    appLayer,
	}

	if proxies, proxiesExist := os.LookupEnv("TRUSTED_PROXIES"); proxiesExist {
		h.engine.SetTrustedProxies(strings.Split(proxies, ","))
	}

	h.setupRedirectRoutes()
	h.setupApiRoutes()

	return h
}

func (h *httpLayer) setupRedirectRoutes() {
	// If no route is matched, redirect to the main page
	h.engine.NoRoute(h.findRedirect)
}

func (h *httpLayer) setupApiRoutes() {
	api := h.engine.Group("/api/v1")
	{
		// Health Check
		api.GET("/ping", h.GetPing)

		// Login Route
		api.POST("/token", h.createToken)

		adminReset := api.Group("/").Use(AuthMiddleware(h, "admin", "reset"))
		{
			adminReset.PUT("/password", h.updatePassword)
		}

		// Admin Authenticated Routes
		adminApi := api.Group("/").Use(AuthMiddleware(h, "admin", "api"))
		{
			adminApi.GET("/self", h.getSelf)

			adminApi.GET("/redirect", h.getRedirects)
			adminApi.GET("/redirect/:id", h.getRedirect)
			adminApi.POST("/redirect", h.createRedirect)
			adminApi.PUT("/redirect/:id", h.updateRedirect)
			adminApi.DELETE("/redirect/:id", h.deleteRedirect)
		}
	}
}

func (h *httpLayer) Run() {
	if addr, addrExists := os.LookupEnv("LISTEN_ADDR"); !addrExists {
		h.engine.Run("0.0.0.0:8080")
	} else {
		h.engine.Run(addr)
	}
}
