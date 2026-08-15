package middlewares

import (
	"github.com/gin-gonic/gin"
)

// CorsMiddleware CORS中间件
// 认证基于 Authorization 头(非Cookie)，无需携带凭据；
// 仅对携带 Origin 的跨域请求回显该 Origin，避免 "*"+Credentials 的非法组合。
// 注: 生产环境前端同源部署(内嵌静态资源)，本中间件未在 main.go 注册，属防御性代码。
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-API-Key")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Language, Content-Type")
			// 基于Authorization头认证，不开启凭据，防止反射Origin+Cookie的CSRF风险
		}

		// 处理OPTIONS预检请求
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
