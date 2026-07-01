package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Hello is a simple test/debug endpoint.
// GET /api/hello -> {"message": "Hello, World!"}
func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}