package authentication

import (
	"fmt"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddlewareFactory struct {
	ProfileService *service.Profile `inject:""`
	AuthMiddleware *jwt.GinJWTMiddleware
}

func (a *AuthMiddlewareFactory) Init() {

	a.AuthMiddleware = &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    8760 * time.Hour, // One year
		MaxRefresh: 8760 * time.Hour, // One year
		Authenticator: func(username string, password string, c *gin.Context) (string, bool) {
			return username, a.ProfileService.Authenticate(username, password)
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		PayloadFunc: func(username string) map[string]interface{} {
			fmt.Println(username)
			user, _ := a.ProfileService.GetProfileByUsername(username)
			return map[string]interface{}{
				"configurationId": user.ConfigurationID,
			}
		},
	}
}
