package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPublicUsageRoutes exposes only privacy-reduced public usage summaries.
func RegisterPublicUsageRoutes(v1 *gin.RouterGroup, h *handler.Handlers, panelRateLimiter *middleware.PanelRateLimiter) {
	public := v1.Group("/public")
	public.Use(panelRateLimiter.PublicIP())
	public.GET("/recent-usage-activity", h.Usage.PublicActivityTrend)
}
