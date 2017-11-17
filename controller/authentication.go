package controller

import (
	"github.com/appleboy/gin-jwt"
	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/gin-gonic/gin"
)

type Authentication struct {
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (a Authentication) InitRoutes() {

	a.Router.POST("/api/v1/auth", a.AuthMiddlewareFactory.AuthMiddleware.LoginHandler)

	auth := a.Router.Group("/api/v1/auth")
	auth.Use(a.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())
	{
		auth.GET("/", a.GetAuthentication)
		auth.PUT("/", a.AuthMiddlewareFactory.AuthMiddleware.RefreshHandler)
	}
}

func (a Authentication) GetAuthentication(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"claims": claims,
	})
}
