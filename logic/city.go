package logic

//
//type City struct {
//	Province string `json:"province"`
//	City     string `json:"city"`
//}
//
//type Response struct {
//	Data struct {
//		Code int    `json:"code"`
//		Info []City `json:"info"`
//	} `json:"data"`
//}
//
//type CityData2 struct {
//	Province string   `json:"province"`
//	CityList []string `json:"citylist"`
//}
//
//type cacheItem struct {
//	Timestamp time.Time
//	Data      []byte
//}
//
//var (
//	cityCache sync.Map
//)
//
//func getInfoNodeAPI(username, password string) ([]byte, error) {
//	url := fmt.Sprintf("http://www.91jlsy.com/demoAPI/get_info_node.php?username=%s&password=%s", username, password)
//	response, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer response.Body.Close()
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	return body, nil
//}
//
//func getCachedCityData(username, password string) ([]CityData2, error) {
//	cacheData, found := cityCache.Load(username)
//	if found {
//		cacheItem := cacheData.(*cacheItem)
//		if !isCacheExpired(cacheItem.Timestamp) {
//			var newData []CityData2
//			err := json.Unmarshal(cacheItem.Data, &newData)
//			if err != nil {
//				return nil, err
//			}
//			return newData, nil
//		}
//	}
//	response, err := getInfoNodeAPI(username, password)
//	if err != nil {
//		return nil, err
//	}
//
//	// 解析原始响应数据
//	var originalData Response
//	err = json.Unmarshal(response, &originalData)
//	if err != nil {
//		return nil, err
//	}
//
//	// 创建省份到城市列表的映射
//	cityMap := make(map[string][]string)
//	for _, city := range originalData.Data.Info {
//		cityMap[city.Province] = append(cityMap[city.Province], city.City)
//	}
//
//	// 构建新的响应数据
//	var newData []CityData2
//	for province, cities := range cityMap {
//		newData = append(newData, CityData2{
//			Province: province,
//			CityList: cities,
//		})
//	}
//
//	// 对城市列表进行排序
//	for _, data := range newData {
//		sort.Strings(data.CityList)
//	}
//
//	// 更新缓存
//	cacheItem := &cacheItem{
//		Timestamp: time.Now(),
//		Data:      nil,
//	}
//
//	// 将新的响应数据转换为 JSON 格式并存储到缓存中
//	cacheItem.Data, err = json.Marshal(newData)
//	if err != nil {
//		return nil, err
//	}
//	cityCache.Store(username, cacheItem)
//
//	return newData, nil
//}
//func GetCity(c *gin.Context, username, password string) {
//	cityData, err := getCachedCityData(username, password)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"data": cityData,
//	})
//}
