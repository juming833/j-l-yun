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
	expirationTime := timestamp.Add(time.Duration(Data.CacheTime) * time.Second)
	return time.Now().After(expirationTime)
}

func getInfoGameAPI(token string) ([]byte, error) {
	url := fmt.Sprintf("http://www.91jlsy.com/demoAPI/get_info_game.php?token=%s", token)
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

func GetGame(c *gin.Context, token string) {
	cacheData, found := gameCache.Load(token)
	if found {
		cacheItem := cacheData.(*CacheItem)
		if !isCacheExpired(cacheItem.Timestamp) {
			c.Data(http.StatusOK, "application/json", cacheItem.Data)
			logger.Debug("使用缓存：" + string(cacheItem.Data))
			return
		}
	}
	response, err := getInfoGameAPI(token)
	if err != nil {
		logger.Error("获取信息失败：" + err.Error())
		return
	}
	cacheItem := &CacheItem{
		Data:      response,
		Timestamp: time.Now(),
	}
	gameCache.Store(token, cacheItem)
	c.Data(http.StatusOK, "application/json", response)
	logger.Debug(string(response))
}
