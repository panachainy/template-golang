package handlers

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	Login(c *gin.Context)
	AuthCallback(c *gin.Context)
	Logout(c *gin.Context)
	Information(c *gin.Context)
	Routes(routerGroup *gin.RouterGroup)
}
