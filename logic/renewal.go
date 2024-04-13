package logic

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Renewal(c *gin.Context, token string) {
	sk5userid := c.PostForm("sk5userid")
	// 调用API接口
	url := fmt.Sprintf("http://www.91jlsy.com/demoAPI/Order_renewal.php?token=%s&sk5userid=%s", token, sk5userid)

	// 发起 HTTP 请求
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make HTTP request: " + err.Error()})
		logger.Error(err.Error())
		return
	}
	defer response.Body.Close()

	// 处理API的响应
	if response.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API请求失败"})
		logger.Error("API请求失败")
		return
	}
	// 这里假设API返回的是JSON数据
	var data struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}
	res := gin.H{"code": data.Code, "data": data.Data}
	jsonData, err := json.Marshal(res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析失败"})
		logger.Error("Failed to serialize data")
		return
	}
	c.JSON(http.StatusOK, res)
	logger.Debug(string(jsonData))
}
