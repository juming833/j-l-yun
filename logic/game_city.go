package logic

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
	"time"
)

var localCache = cache.New(10*time.Minute, 30*time.Minute)

type GameInfo struct {
	Game     string `json:"game"`
	City     string `json:"city"`
	Count    any    `json:"count"`
	Province string `json:"province"`
}

type ResponseData struct {
	Code int        `json:"code"`
	Info []GameInfo `json:"info"`
}

type ResponseData2 struct {
	Code int        `json:"code"`
	Info []CityData `json:"info"`
}

type GameCity struct {
	Data ResponseData `json:"data"`
}

type GameCity2 struct {
	Data ResponseData2 `json:"data"`
}

type CityData struct {
	Game     string         `json:"game"`
	Province string         `json:"province"`
	CityList map[string]any `json:"citylist"`
}

func GetGameCity(c *gin.Context, token string) {
	game := c.Query("game")
	cacheKey := fmt.Sprintf("gamecity-%s", game)
	if cachedData, found := localCache.Get(cacheKey); found {
		c.JSON(http.StatusOK, cachedData)
		return
	}
	url := fmt.Sprintf("http://www.91jlsy.com/demoAPI/surpluslist.php?token=%s&game=%s", token, game)
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make HTTP request: " + err.Error()})
		logger.Error(err.Error())
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body: " + err.Error()})
		logger.Error(err.Error())
		return
	}

	var responseData GameCity
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON unmarshal error: " + err.Error()})
		logger.Error(err.Error())
		return
	}

	// Group city data by province
	cityData := make(map[string]CityData)
	for _, info := range responseData.Data.Info {
		province := info.Province
		if data, ok := cityData[province]; ok {
			data.CityList[info.City] = info.Count
			cityData[province] = data
		} else {
			data := CityData{
				Game:     info.Game,
				Province: province,
				CityList: map[string]any{
					info.City: info.Count,
				},
			}
			cityData[province] = data
		}
	}

	// Construct the final result
	result := ResponseData2{
		Code: responseData.Data.Code,
		Info: make([]CityData, 0, len(cityData)),
	}
	for _, data := range cityData {
		result.Info = append(result.Info, data)
	}

	finalResult := GameCity2{Data: result}
	jsonData, err := json.Marshal(finalResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize data: " + err.Error()})
		logger.Error("Failed to serialize data: " + err.Error())
		return
	}

	// Cache the final result
	localCache.Set(cacheKey, finalResult, time.Duration(Data.CacheTime)*time.Second)

	// Return data
	c.JSON(http.StatusOK, finalResult)
	logger.Debug(string(jsonData))
}
