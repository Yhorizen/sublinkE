package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"sublink/dto"
	"sublink/models"
	"time"
)

func GenerateAPIKey(c *gin.Context) {
	var userAccessKey dto.UserAccessKey
	if err := c.BindJSON(&userAccessKey); err != nil {
		c.JSON(500, gin.H{"msg": "参数错误"})
		return
	}
	// 归属校验: 仅管理员可为他人生成，普通用户只能为自己生成
	current, err := currentUser(c)
	if err != nil {
		c.JSON(401, gin.H{"msg": "未登录"})
		return
	}
	if current.Role != "admin" && userAccessKey.UserName != current.Username {
		c.JSON(403, gin.H{"msg": "无权为其他用户生成API Key"})
		return
	}
	user := &models.User{Username: userAccessKey.UserName}
	err = user.Find()
	if err != nil {
		c.JSON(400, gin.H{"msg": "用户不存在"})
		return
	}

	var accessKey models.AccessKey
	accessKey.ExpiredAt = userAccessKey.ExpiredAt
	accessKey.Description = userAccessKey.Description
	accessKey.UserID = user.ID
	accessKey.CreatedAt = time.Now()
	accessKey.Username = user.Username

	apiKey, err := accessKey.GenerateAPIKey()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"msg": "生成API Key失败"})
		return
	}
	err = accessKey.Generate()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"msg": "生成API Key失败"})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": map[string]string{
			"accessKey": apiKey,
		},
		"msg": "API Key生成成功",
	})
}

func DeleteAPIKey(c *gin.Context) {
	apiKeyIDParam := c.Param("apiKeyId")
	if apiKeyIDParam == "" {
		c.JSON(400, gin.H{"msg": "缺少API Key ID"})
		return
	}
	apiKeyID, err := strconv.Atoi(apiKeyIDParam)
	if err != nil {
		c.JSON(500, gin.H{"msg": "删除API Key失败"})
		return
	}
	// 归属校验: 仅管理员可删除他人Key，普通用户只能删除自己的
	current, err := currentUser(c)
	if err != nil {
		c.JSON(401, gin.H{"msg": "未登录"})
		return
	}
	var fullKey models.AccessKey
	if err := models.DB.First(&fullKey, apiKeyID).Error; err != nil {
		c.JSON(400, gin.H{"msg": "API Key不存在"})
		return
	}
	if current.Role != "admin" && fullKey.Username != current.Username {
		c.JSON(403, gin.H{"msg": "无权删除他人API Key"})
		return
	}

	var accessKey models.AccessKey
	accessKey.ID = apiKeyID
	err = accessKey.Delete()
	if err != nil {
		c.JSON(500, gin.H{"msg": "删除API Key失败"})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "删除API Key成功",
	})
}

func GetAPIKey(c *gin.Context) {
	userIDParam := c.Param("userId")
	if userIDParam == "" {
		c.JSON(400, gin.H{"msg": "缺少User ID"})
		return
	}
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(500, gin.H{"msg": "删除API Key失败"})
		return
	}
	// 归属校验: 仅管理员可查看他人Key，普通用户只能查看自己的
	current, err := currentUser(c)
	if err != nil {
		c.JSON(401, gin.H{"msg": "未登录"})
		return
	}
	if current.Role != "admin" && current.ID != userID {
		c.JSON(403, gin.H{"msg": "无权查看他人API Key"})
		return
	}
	apiKeys, err := models.FindValidAccessKeys(userID)
	if err != nil {
		c.JSON(500, gin.H{"msg": "查询API Key失败"})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": apiKeys,
		"msg":  "查询API Key成功",
	})
}
