package main

import (
	"91jlsy/api/logic"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	err := logic.LoadConfig()
	if err != nil {
		log.Fatal("无法加载配置文件:", err)
	}
	logic.InitLogger()
	defer logic.CloseLogger()
	logFile, err := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err) // 无法打开或创建文件时退出
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {

		}
	}(logFile)
	// 设置Gin的日志输出到文件
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	r := gin.Default()
	r.Use(VerifySign)
	r.GET("/api/getGame", func(c *gin.Context) {
		logic.GetGame(c, logic.Data.Token)
	})
	r.GET("/api/GetGameCity", func(c *gin.Context) {
		logic.GetGameCity(c, logic.Data.Token)
	})
	{
		r.POST("/api/postOrder", func(c *gin.Context) {
			if logic.Data.Test == true {
				logic.BuyOrder1(c, logic.Data.Token, logic.Data.Username, logic.Data.Password)
			} else {
				logic.BuyOrder(c, logic.Data.Token, logic.Data.Username, logic.Data.Password)
			}
		})
		r.POST("/api/Renewal", func(c *gin.Context) {
			if logic.Data.Test == true {
				logic.Renewal1(c, logic.Data.Token)
			} else {
				logic.Renewal(c, logic.Data.Token)
			}
		})
		r.POST("/api/Change", func(c *gin.Context) {
			if logic.Data.Test == true {
				logic.Change1(c, logic.Data.Token)
			} else {
				logic.Change(c, logic.Data.Token)
			}
		})
		r.POST("/api/Unsubscribe", func(c *gin.Context) {
			if logic.Data.Test == true {
				logic.Unsubscribe1(c, logic.Data.Token)
			} else {
				logic.Unsubscribe(c, logic.Data.Token)
			}
		})
		r.POST("/api/ChangeCity", func(c *gin.Context) {
			if logic.Data.Test == true {
				logic.ChangeCity1(c, logic.Data.Token)
			} else {
				logic.ChangeCity(c, logic.Data.Token)
			}
		})
	}
	port := logic.Data.Port
	if err := r.Run(":" + port); err != nil {
		panic("gin 启动失败")
	}
}

func VerifySign(c *gin.Context) {
	adminid := c.DefaultQuery("adminid", "")
	timestamp := c.DefaultQuery("ti", "")
	nonce := c.DefaultQuery("nonce", "")
	sign := c.DefaultQuery("sign", "")
	apiKey := logic.Data.ApiKey

	if c.Request.Method == http.MethodPost {
		adminid = c.DefaultPostForm("adminid", "")
		timestamp = c.DefaultPostForm("ti", "")
		nonce = c.DefaultPostForm("nonce", "")
		sign = c.DefaultPostForm("sign", "")
	}
	currentTime := time.Now().Unix()
	requestTime, _ := strconv.ParseInt(timestamp, 10, 64)
	if math.Abs(float64(currentTime-requestTime)) > 30 { // 允许30秒的时间差
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
