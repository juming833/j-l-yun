package logic

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

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

func GetGameCity(c *gin.Context, username, password string) {
	game := c.Query("game")
	// 构造 URL
	url := fmt.Sprintf("http://www.91jlsy.com/demoAPI/surpluslist.php?username=%s&password=%s&game=%s", username, password, game)

	// 发起 HTTP 请求
	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	// 读取响应数据
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 解析响应数据
	var responseData GameCity
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 根据省份分组城市数据
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

	// 构造最终结果
	result := ResponseData2{
		Code: responseData.Data.Code,
		Info: make([]CityData, 0, len(cityData)),
	}
	for _, data := range cityData {
		result.Info = append(result.Info, CityData{
			Game:     data.Game,
			Province: data.Province,
			CityList: data.CityList,
		})
	}
	// 返回数据
	c.JSON(http.StatusOK, GameCity2{Data: result})
}
