package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type testURLDeps struct {
	router *gin.Engine
}

func newTestURLHandler(deps testURLDeps) {
	deps.router.GET("/", testURL)
}

func testURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
