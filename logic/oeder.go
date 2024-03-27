package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func BuyOrder(c *gin.Context, username, password string) {
	daili := c.Query("daili")
	nodetime := c.Query("nodetime")
	gamename := c.Query("gamename")
	city := c.Query("city")
	// 构建URL
	url := fmt.Sprintf("http://www.91jlsy.com/demoAPI/buy_order.php?username=%s&password=%s&daili=%s&nodetime=%s&gamename=%s&city=%s",
		username, password, daili, nodetime, gamename, city)

	// 发送HTTP GET请求
	response, err := http.Get(url)
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
	c.Data(http.StatusOK, "application/json", body)
}
