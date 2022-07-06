package web

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (w *ScanInfosService) Router(router *gin.Engine) *gin.Engine {
	router.Use(w.RequestLoggerMiddleware())
	router.NoRoute(w.NotFoundHandler())
	api := router.Group("/api/v1")
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Accept", "content-type", "User-Agent", "Accept-Language", "Referer", "DNT", "Connection", "Pragma", "Cache-Control", "TE"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api.POST("/scaninfos", w.StoreScanInfosHandler())
	api.GET("/scaninfos/:id", w.GetScanInfosHandler())
	api.GET("/scaninfos", w.GetAllScanInfosHandler())
	api.PUT("/scaninfos", w.UpdateScanInfosHandler())
	api.DELETE("/scaninfos/:id", w.DeleteScanInfosHandler())
	return router
}
