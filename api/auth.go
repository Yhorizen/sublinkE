package api

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"sublink/middlewares"
	"sublink/models"
	"sublink/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取token
func GetToken(username string, tokenVersion int) (string, error) {
	// 过期天数取配置(默认14天)
	config := models.ReadConfig()
	expireDays := config.ExpireDays
	if expireDays <= 0 {
		expireDays = 14
	}
	c := &middlewares.JwtClaims{
		Username:     username,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * time.Duration(expireDays))),
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			Subject:   username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(middlewares.Secret)
}

// 获取滑块验证码
func GetCaptcha(c *gin.Context) {
	id, bs4, err := utils.GetCaptcha()
	if err != nil {
		log.Println("获取验证码失败")
		c.JSON(400, gin.H{
			"msg": "获取验证码失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"captchaKey":    id,
			"captchaBase64": bs4,
			"captchaType":   "slider",
			"bgWidth":       utils.SliderBgWidth,
			"bgHeight":      utils.SliderBgHeight,
			"pieceSize":     utils.SliderPieceSize,
		},
		"msg": "获取验证码成功",
	})

}

// 用户登录
func UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	captchaCode := c.PostForm("captchaCode")
	captchaKey := c.PostForm("captchaKey")
	// 验证验证码(位置 + 拖动轨迹)
	if !utils.VerifyCaptcha(captchaKey, captchaCode) || !utils.VerifyTrajectory(c.PostForm("trajectory")) {
		log.Println("验证码错误")
		c.JSON(400, gin.H{
			"msg": "验证码错误",
		})
		return
	}
	user := &models.User{Username: username, Password: password}
	err := user.Verify()
	if err != nil {
		log.Println("账号或者密码错误")
		c.JSON(400, gin.H{
			"msg": "账号或者密码错误",
		})
		return
	}
	// 生成token
	token, err := GetToken(username, user.TokenVersion)
	if err != nil {
		log.Println("获取token失败")
		c.JSON(400, gin.H{
			"msg": "获取token失败",
		})
		return
	}
	// 登录成功返回token
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"accessToken":  token,
			"tokenType":    "Bearer",
			"refreshToken": nil,
			"expires":      nil,
		},
		"msg": "登录成功",
	})
}
func UserOut(c *gin.Context) {
	// 拿到jwt中的username
	username, _ := c.Get("username")
	u, _ := username.(string)
	if u == "" {
		c.JSON(400, gin.H{"msg": "未登录"})
		return
	}
	// 自增token版本号，使已签发的JWT立即失效(服务端登出吊销)
	if err := models.IncrementTokenVersion(u); err != nil {
		log.Println("登出吊销token失败:", err)
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "退出成功",
	})
}
