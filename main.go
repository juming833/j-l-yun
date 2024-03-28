package main

import (
	"91jlsy/api/logic"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	err := logic.LoadConfig()
	if err != nil {
		log.Fatal("无法加载配置文件:", err)
	}

	r := gin.Default()
	//r.Use(VerifySign)
	r.GET("/api/getCity", func(c *gin.Context) {
		logic.GetCity(c, logic.Data.Username, logic.Data.Password)
	})
	r.GET("/api/getGame", func(c *gin.Context) {
		logic.GetGame(c, logic.Data.Username, logic.Data.Password)
	})
	r.POST("/api/postOrder", func(c *gin.Context) {
		logic.BuyOrder(c, logic.Data.Username, logic.Data.Password)
	})
	r.GET("/api/GetGameCity", func(c *gin.Context) {
		logic.GetGameCity(c, logic.Data.Username, logic.Data.Password)
	})
	port := logic.Data.Port
	if err := r.Run(":" + port); err != nil {
		panic("gin 启动失败")
	}
}

func VerifySign(c *gin.Context) {
	adminid := c.Query("adminid")
	timestamp := c.Query("ti")
	nonce := c.Query("nonce")
	sign := c.Query("sign")
	apiKey := logic.Data.ApiKey

	currentTime := time.Now().Unix()
	requestTime, _ := strconv.ParseInt(timestamp, 10, 64)
	if math.Abs(float64(currentTime-requestTime)) > 300 { // 允许30秒的时间差
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid timestamp",
		})
		c.Abort()
		return
	}

	//验证nonce的唯一性
	if !isNonceUnique(nonce) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Duplicate nonce",
		})
		c.Abort()
		return
	}
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

var nonceCache sync.Map

func isNonceUnique(nonce string) bool {
	// 检查nonce是否已经存在于缓存中
	_, loaded := nonceCache.Load(nonce)
	if loaded {
		// nonce已经存在，表示已被使用过
		return false
	}

	// 将nonce存储到缓存中，并设置过期时间为24小时
	expiration := time.Now().Add(24 * time.Hour)
	nonceCache.Store(nonce, expiration)

	// 启动一个goroutine来定期清理过期的nonce
	go cleanExpiredNonces()

	return true
}

func cleanExpiredNonces() {
	// 每隔一段时间（例如1小时）清理一次过期的nonce
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		nonceCache.Range(func(key, value interface{}) bool {
			expiration := value.(time.Time)
			if now.After(expiration) {
				// nonce已过期，从缓存中删除
				nonceCache.Delete(key)
			}
			return true
		})
	}
}
