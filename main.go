package main

import (
	"91jlsy/api/logic"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	err := logic.LoadConfig()
	if err != nil {
		log.Fatal("无法加载配置文件:", err)
	}

	r := gin.Default()
	//r.Use(VerifySign)
	//r.GET("/api/getCity", func(c *gin.Context) {
	//	logic.GetCity(c, logic.Data.Username, logic.Data.Password)
	//})
	r.GET("/api/getGame", func(c *gin.Context) {
		logic.GetGame(c, logic.Data.Username, logic.Data.Password)
	})
	r.POST("/api/postOrder", func(c *gin.Context) {
		logic.BuyOrder(c, logic.Data.Username, logic.Data.Password)
	})
	r.GET("/api/GetGameCity", func(c *gin.Context) {
		logic.GetGameCity(c, logic.Data.Username, logic.Data.Password)
	})
	if err := r.Run(":8080"); err != nil {
		panic("gin 启动失败")
	}
}

func VerifySign(c *gin.Context) {
	adminid := c.Query("adminid")
	timestamp := c.Query("ti")
	nonce := c.Query("nonce")
	sign := c.Query("sign")
	apiKey := logic.Data.ApiKey
	// 将时间戳、nonce和API密钥进行拼接
	signStr := timestamp + nonce + apiKey
	// 使用MD5哈希算法计算签名的摘要
	signBytes := md5.Sum([]byte(signStr))
	// 将签名的摘要转换为16进制字符串
	expectedSign := hex.EncodeToString(signBytes[:])
	// 验证签名是否匹配
	if sign != expectedSign || adminid != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid sign",
		})
		c.Abort()
		return
	}
	c.Next()
}
