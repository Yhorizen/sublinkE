package routers

import (
	"time"

	"sublink/api"
	"sublink/middlewares"

	"github.com/gin-gonic/gin"
)

func User(r *gin.Engine) {
	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/login", middlewares.RateLimit(10, 15*time.Minute), api.UserLogin)
		authGroup.DELETE("/logout", api.UserOut)
		authGroup.GET("/captcha", api.GetCaptcha)
		authGroup.POST("/register", middlewares.RateLimit(5, time.Hour), api.UserRegister)
	}
	userGroup := r.Group("/api/v1/users")
	{
		userGroup.GET("/me", api.UserMe)
		userGroup.GET("/page", api.UserPages)
		userGroup.POST("/update", api.UserSet)
		userGroup.GET("/pull-logs", api.UserPullLogs)

	}
}
