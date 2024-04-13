package logic

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func writeSk5UserIDToFile(infos []map[string]string) error {
	fileName := "sk5userid.txt"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("无法打开文件 %s: %s", fileName, err)
	}
	defer file.Close()
	for _, info := range infos {
		sk5UserID := info["sk5userid"]
		_, err = file.WriteString(sk5UserID + "\n")
		if err != nil {
			return fmt.Errorf("无法写入sk5userid到文件 %s: %s", fileName, err)
		}
	}

	return nil
}
func checkSk5UserID(sk5userid string) error {
	fileName := "sk5userid.txt"
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("无法打开文件 %s: %s", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == sk5userid {
			return nil // 找到匹配的sk5userid
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("从文件 %s 读取行时出错: %s", fileName, err)
	}
	return fmt.Errorf("文件中找不到匹配的sk5userid")
}

// 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RandomDigits 生成int随机数
func RandomDigits(n int) string {
	var digits = []rune("0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = digits[rand.Intn(len(digits))]
	}
	return string(s)
}

// 生成随机端口
func randomPort() string {
	return fmt.Sprintf("%d", rand.Intn(65535-1024)+1024)
}

// 中文月份到数字的映射
var monthMap = map[string]int{
	"一月": 1,
	"二月": 2,
	"三月": 3,
	"四月": 4,
	"五月": 5,
	"半年": 6,
	"一年": 12,
	"二年": 24,
	"三年": 36,
}

// 生成随机订单ID
func randomOrderID() string {
	timestamp := time.Now().Format("20060102150405")
	randomNum := fmt.Sprintf("%12d", rand.Intn(999999999999))
	uniqueID := timestamp + randomNum
	return uniqueID
}
func BuyOrder1(c *gin.Context, token, username, password string) {
	nodetime := c.PostForm("nodetime")
	city := c.PostForm("city")
	numStr := c.PostForm("n")
	gamename := c.PostForm("gamename")
	// 将n转换为整数
	num, err := strconv.Atoi(numStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的num值"})
		return
	}
	// 解析nodetime
	monthsToAdd, ok := monthMap[nodetime]
	if !ok {
		// 如果nodetime不是有效的中文月份，返回错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的nodetime值"})
		return
	}
	// 计算todata时间
	now := time.Now()
	todata := now.AddDate(0, monthsToAdd, 0)
	// 根据num生成info数组
	var infos []map[string]string
	for i := 0; i < num; i++ {
		info := map[string]string{
			"game":      gamename,
			"city":      city,
			"sk5port":   randomPort(),
			"httpport":  randomPort(),
			"sk5user":   randomString(12),
			"sk5pass":   randomString(4),
			"sk5userid": randomString(7) + "-" + RandomDigits(20),
			"Gaddress":  fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255)),
		}
		infos = append(infos, info)
	}

	// 生成最终的响应内容
	responseData := gin.H{
		"data": gin.H{
			"info": infos,
		},
		"nowtime": now.Format("2006-01-02 15:04:05"),
		"todata":  todata.Format("2006-01-02 15:04:05"),
		"code":    200,
		"orderid": randomOrderID(),
	}
	err = writeSk5UserIDToFile(infos)
	if err != nil {
		log.Println("无法写入sk5userid到文件:", err)
		// 可以根据需要返回适当的错误响应
	}

	c.JSON(http.StatusOK, responseData)
}

func Renewal1(c *gin.Context, token string) {
	sk5userid := c.PostForm("sk5userid")
	err := checkSk5UserID(sk5userid)
	if err != nil {
		response := gin.H{
			"code":  http.StatusBadRequest,
			"error": "无效的sk5userid",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	res := gin.H{
		"code": 200,
		"data": "续费成功!",
	}
	c.JSON(http.StatusOK, res)
}
func Unsubscribe1(c *gin.Context, token string) {
	sk5userid := c.PostForm("sk5userid")
	err := checkSk5UserID(sk5userid)
	if err != nil {
		response := gin.H{
			"code":  http.StatusBadRequest,
			"error": "无效的sk5userid",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	res := gin.H{
		"code": 200,
		"data": "退订提交成功,等待审核!",
	}
	c.JSON(http.StatusOK, res)
}
func Change1(c *gin.Context, token string) {
	sk5userid := c.PostForm("sk5userid")
	err := checkSk5UserID(sk5userid)
	if err != nil {
		response := gin.H{
			"code":  http.StatusBadRequest,
			"error": "无效的sk5userid",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	res := gin.H{
		"code": 200,
		"data": RandomDigits(4),
	}
	c.JSON(http.StatusOK, res)
}
func ChangeCity1(c *gin.Context, token string) {
	game := c.PostForm("gamename")
	city := c.PostForm("city")
	sk5userid := c.PostForm("sk5userid")
	err := checkSk5UserID(sk5userid)
	if err != nil {
		response := gin.H{
			"code":  http.StatusBadRequest,
			"error": "无效的sk5userid",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// 生成随机端口和 IP 地址
	rand.Seed(time.Now().UnixNano())
	port := strconv.Itoa(rand.Intn(65535))
	ip := strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256)) + "." + strconv.Itoa(rand.Intn(256))
	newsk5userid := randomString(7) + "-" + RandomDigits(20) // 假设 UserID 是随机 20 个字符

	// 构建响应数据
	response := map[string]interface{}{
		"0":    sk5userid,
		"code": 200,
		"data": map[string]interface{}{
			"info": []map[string]string{
				{
					"game":         game,
					"city":         city,
					"port":         port,
					"Gaddress":     ip,
					"newsk5userid": newsk5userid,
					"oldsk5userid": sk5userid,
				},
			},
		},
	}
	// 发送 JSON 响应
	c.JSON(http.StatusOK, response)
}
