package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

func BuyOrder(c *gin.Context, username, password string) {
	daili := c.PostForm("daili")
	nodetime := c.PostForm("nodetime")
	gamename := c.PostForm("gamename")
	city := c.PostForm("city")
	encodedNodetime := url.QueryEscape(nodetime)
	encodedGamename := url.QueryEscape(gamename)
	encodedCity := url.QueryEscape(city)

	// 构建URL
	url := fmt.Sprintf("http://www.91jlsy.com/demoAPI/buy_order.php?username=%s&password=%s&daili=%s&nodetime=%s&gamename=%s&city=%s",
		username, password, daili, encodedNodetime, encodedGamename, encodedCity)

	// 创建新的请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置请求头
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "www.91jlsy.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", fmt.Sprintf("username=%s; password=%s", username, password))

	// 发送HTTP请求
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "text/html", body)
}
