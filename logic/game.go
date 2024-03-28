package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// @Summary 获取游戏列表
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} Article "成功"
// @Failure 400 {object} string "请求错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/getGame [get]

var gameCache sync.Map

type CacheItem struct {
	Data      []byte
	Timestamp time.Time
}

func isCacheExpired(timestamp time.Time) bool {
	expirationTime := timestamp.Add(10 * time.Minute)
	return time.Now().After(expirationTime)
}

func getInfoGameAPI(username, password string) ([]byte, error) {
	url := fmt.Sprintf("http://www.91jlsy.com/demoAPI/get_info_game.php?username=%s&password=%s", username, password)
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetGame(c *gin.Context, username, password string) {
	cacheData, found := gameCache.Load(username)
	if found {
		cacheItem := cacheData.(*CacheItem)
		if !isCacheExpired(cacheItem.Timestamp) {
			c.Data(http.StatusOK, "application/json", cacheItem.Data)
			return
		}
	}

	response, err := getInfoGameAPI(username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cacheItem := &CacheItem{
		Data:      response,
		Timestamp: time.Now(),
	}
	gameCache.Store(username, cacheItem)

	c.Data(http.StatusOK, "application/json", response)
}
