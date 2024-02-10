package api

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/haven/upload-service/api/handler"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", handler.Home)
	r.POST("/upload", handler.Upload)

}
