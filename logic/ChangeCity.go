package logic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type BatchChangeResponse struct {
	Code int `json:"code"`
	Data struct {
		Info []struct {
			Game         string `json:"game"`
			City         string `json:"city"`
			Port         string `json:"port"`
			GAddress     string `json:"Gaddress"`
			NewSk5UserID string `json:"newsk5userid"`
			OldSk5UserID string `json:"oldsk5userid"`
		} `json:"info"`
	} `json:"data"`
}

func ChangeCity(c *gin.Context, token string) {
	sk5UserID := c.PostForm("sk5userid")
	gamename := c.PostForm("gamename")
	city := c.PostForm("city")
	encodedGamename := url.QueryEscape(gamename)
	encodedCity := url.QueryEscape(city)

	url := fmt.Sprintf("http://www.91jlsy.com/demoAPI/Batch_change.php?token=%s&sk5userid=%s&game=%s&city=%s", token, sk5UserID, encodedGamename, encodedCity)

	// Send HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make HTTP request: " + err.Error()})
		logger.Error(err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API请求失败"})
		logger.Error("API请求失败")
		return
	}

	var batchChangeResponse BatchChangeResponse
	err = json.NewDecoder(response.Body).Decode(&batchChangeResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error(err.Error())
		return
	}
	res := gin.H{
		"0":    batchChangeResponse.Data.Info[0].OldSk5UserID,
		"data": batchChangeResponse.Data,
		"code": batchChangeResponse.Code,
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize data"})
		logger.Error("Failed to serialize data")
		return
	}
	logger.Debug(string(jsonData))
	c.JSON(http.StatusOK, res)
}
