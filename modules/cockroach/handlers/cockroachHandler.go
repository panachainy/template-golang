package handlers

import "github.com/gin-gonic/gin"

type CockroachHandler interface {
	DetectCockroach(c *gin.Context)
}
