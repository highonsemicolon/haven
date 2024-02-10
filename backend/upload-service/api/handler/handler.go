package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome to haven's upload service."})
}

func Upload(c *gin.Context) {
	var requestData struct {
		Url string `json:"q"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s url not found.", requestData.Url)})
	}
}
