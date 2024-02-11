package api

import "github.com/gin-gonic/gin"

func home(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome to haven's deployment service."})
}
