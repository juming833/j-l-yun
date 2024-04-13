package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func BuyOrder(c *gin.Context, token string) {
	nodetime := c.PostForm("nodetime")
	gamename := c.PostForm("gamename")
	city := c.PostForm("city")
	num := c.PostForm("n")
	user := c.PostForm("user")
	proxyuser := c.PostForm("proxyuser")
	proxypassword := c.PostForm("proxypassword")
	//encodedNodetime := url.QueryEscape(nodetime)
	//encodedGamename := url.QueryEscape(gamename)
	//encodedCity := url.QueryEscape(city)
	Baseurl := fmt.Sprintf("http://www.91jlsy.com/demoAPI/buy_order.php?token=%s&nodetime=%s&gamename=%s&city=%s&n=%s", token, nodetime, gamename, city, num)
	params := url.Values{}
	if user == "zdy" {
		params.Add("user", user)
		params.Add("proxyuser", proxyuser)
		params.Add("proxypassword", proxypassword)
	}
	Baseurl += "&" + params.Encode()
	// 创建新的请求
	req, err := http.NewRequest("GET", Baseurl, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}

	// 设置请求头
	//req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	//req.Header.Add("Accept", "*/*")
	//req.Header.Add("Host", "www.91jlsy.com")
	//req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Cookie", fmt.Sprintf("username=%s; password=%s", username, password))

	// 发送HTTP请求
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}
	defer response.Body.Close()
	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}
	c.Data(http.StatusOK, "text/html", body)
	logJSON(json.RawMessage(body))
}
func logJSON(data interface{}) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false) // 确保HTML相关的字符不会被转义
	// 首先，正常编码JSON数据
	if err := encoder.Encode(data); err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}
	// 将编码的JSON字符串解码到一个临时的interface{}变量中
	var rawData interface{}
	if err := json.Unmarshal(buffer.Bytes(), &rawData); err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	// 清空buffer以用于重新编码
	buffer.Reset()
	// 再次编码，通常这一步不会对非ASCII字符进行转义
	if err := encoder.Encode(rawData); err != nil {
		log.Println("Error re-encoding JSON:", err)
		return
	}
	jsonString := strings.ReplaceAll(buffer.String(), "\n", "")
	// 输出或记录处理后的JSON
	logger.Debug(jsonString)
}
