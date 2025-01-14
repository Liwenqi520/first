package logic

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	AppInit "first/init"
	"first/init/myredis"
	"first/init/mysql"
	"first/internal/defines"
	"first/internal/model"
	"first/internal/schema"
	"fmt"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Liwenqi520/errorx"
	"github.com/Liwenqi520/helper"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"github.com/skip2/go-qrcode"
	"github.com/tealeg/xlsx"
	"gopkg.in/gomail.v2"
)

type Common struct{}

type Data struct {
	CompanyID  string      `json:"company_id"`
	OperatorID string      `json:"operator_id"`
	Data       string      `json:"data"`
	Before     interface{} `json:"before"`
	After      interface{} `json:"after"`
	DataName   string      `json:"data_name"`
	IP         string      `json:"ip"`
}

type DeviceAdminConfig struct {
	IsGoDeviceAdmin   string `json:"is_go_device_admin"`
	DeviceAdminTcpURL string `json:"device_admin_tcp_url"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Type              string `json:"type"`
	Timeout           int    `json:"timeout"`
}

type CaiyunConfig struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

type PushTestDeviceDataConfig struct {
	Domain    string `json:"domain"`
	CompanyID string `json:"company_id"`
}

type SmsConfig struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	RegionId        string `json:"region_id"`
	SignName        string `json:"sign_name"`
	TemplateCode    string `json:"template_code"`
}

type EmailConfig struct {
	SmtpHost       string `json:"smtp_host"`
	SmtpPort       string `json:"smtp_port"`
	SenderEmail    string `json:"sender_email"`
	SenderPassword string `json:"sender_password"`
}

// StrToBase64Str 字符串转base64编码
func (c Common) StrToBase64Str(qrContent string) (base64ImgStr string) {
	// 创建一个二维码
	qrCode, err := qrcode.New(qrContent, qrcode.Highest)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建一个缓冲区用于存储二维码图片的字节数据
	var buf bytes.Buffer
	qrImage := qrCode.Image(256)

	// 将二维码图片写入缓冲区（这里使用PNG格式）
	if err := png.Encode(&buf, qrImage); err != nil {
		fmt.Println(err)
		return
	}
	// 将缓冲区中的数据转换为Base64编码的字符串
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	base64ImgStr = "data:image/png;base64," + base64Str
	return
}

func (c Common) RandomString(slice []string) string {
	if len(slice) == 0 {
		return "" // 如果切片为空，返回空字符串
	}

	// 使用当前时间的纳秒数作为随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成一个0到len(slice)-1之间的随机数
	randomIndex := rand.Intn(len(slice))

	// 返回切片中对应索引的元素
	return slice[randomIndex]
}

func (c Common) GetSyncPluginSwitch() (Value string) {
	db := mysql.NewDB()
	var config model.SystemConfig
	db.Model(&model.SystemConfig{}).Where("key = ? and deleted_time = 0", "is_sync_plugin_to_iot").First(&config)
	Value = config.Value
	return
}

func (c Common) GetEmailConfig() (EmailConfig EmailConfig) {
	emailKey := defines.EmailConfigKey
	rds := myredis.RedisCli()
	Config, err := rds.HGetAll(context.Background(), emailKey).Result()
	if err != nil || len(Config) == 0 {
		var config model.SystemConfig
		db := mysql.NewDB()
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "email_config").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		json.Unmarshal([]byte(Value), &EmailConfig)
		// 保存
		if _, err := rds.HMSet(context.Background(), emailKey,
			"smtp_host", EmailConfig.SmtpHost,
			"smtp_port", EmailConfig.SmtpPort,
			"sender_email", EmailConfig.SenderEmail,
			"sender_password", EmailConfig.SenderPassword,
		).Result(); err != nil {
			return
		}
	} else {
		EmailConfig.SmtpHost = Config["smtp_host"]
		EmailConfig.SmtpPort = Config["smtp_port"]
		EmailConfig.SenderEmail = Config["sender_email"]
		EmailConfig.SenderPassword = Config["sender_password"]
	}
	return
}

func (c Common) GetSmsConfig() (SmsConfig SmsConfig) {
	smskey := defines.SmsConfigKey
	rds := myredis.RedisCli()
	Config, err := rds.HGetAll(context.Background(), smskey).Result()
	if err != nil || len(Config) == 0 {
		var config model.SystemConfig
		db := mysql.NewDB()
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "sms_config").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		json.Unmarshal([]byte(Value), &SmsConfig)
		// 保存
		if _, err := rds.HMSet(context.Background(), smskey,
			"access_key_id", SmsConfig.AccessKeyId,
			"access_key_secret", SmsConfig.AccessKeySecret,
			"region_id", SmsConfig.RegionId,
			"sign_name", SmsConfig.SignName,
			"template_code", SmsConfig.TemplateCode,
		).Result(); err != nil {
			return
		}
	} else {
		SmsConfig.AccessKeyId = Config["access_key_id"]
		SmsConfig.AccessKeySecret = Config["access_key_secret"]
		SmsConfig.RegionId = Config["region_id"]
		SmsConfig.SignName = Config["sign_name"]
		SmsConfig.TemplateCode = Config["template_code"]
	}
	return
}

func (c Common) GetCaiyunConfig() (Caiyun CaiyunConfig) {
	caiyunKey := defines.CaiyunConfigKey
	rds := myredis.RedisCli()
	Config, err := rds.HGetAll(context.Background(), caiyunKey).Result()
	if err != nil || len(Config) == 0 {
		var config model.SystemConfig
		db := mysql.NewDB()
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "caiyun").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		json.Unmarshal([]byte(Value), &Caiyun)
		// 保存
		if _, err := rds.HMSet(context.Background(), caiyunKey,
			"url", Caiyun.Url,
			"key", Caiyun.Key,
		).Result(); err != nil {
			return
		}
	} else {
		Caiyun.Url = Config["url"]
		Caiyun.Key = Config["key"]
	}
	return
}

// GetOfflineHourConfig 获取离线时间配置
func (c Common) GetOfflineHourConfig() (OfflineHour int64) {
	OfflineHourKey := defines.OfflineHourKey
	rds := myredis.RedisCli()
	OfflineHour, err := rds.Get(context.Background(), OfflineHourKey).Int64()
	if err != nil || OfflineHour == 0 {
		db := mysql.NewDB()
		var config model.SystemConfig
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "offline_hour").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		// 保存
		if _, err := rds.Set(context.Background(), OfflineHourKey, Value, 7200*time.Second).Result(); err != nil {
			return
		}
		temp, _ := strconv.Atoi(Value)
		OfflineHour = int64(temp)
	}
	return
}

// IsCheckDeviceConnectedStatus 检查设备是否在线
func (c Common) IsCheckDeviceConnectedStatus() (IsCheck int64) {
	db := mysql.NewDB()
	var config model.SystemConfig
	db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "is_check_device_connected_status").First(&config)
	temp, _ := strconv.Atoi(config.Value)
	IsCheck = int64(temp)
	return
}

// GetAesKeyConfig 获取加密字符串
func (c Common) GetAesKeyConfig() (Key string) {
	aesKey := defines.AesKey
	rds := myredis.RedisCli()
	Key, err := rds.Get(context.Background(), aesKey).Result()
	if err != nil || Key == "" {
		db := mysql.NewDB()
		var config model.SystemConfig
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "aes_key").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		// 保存
		if _, err := rds.Set(context.Background(), aesKey, Value, 7200*time.Second).Result(); err != nil {
			return
		}
		Key = Value
	}
	return
}

// GetAdminEmail 获取管理员邮箱地址
func (c Common) GetAdminEmail() (Email string) {
	AdminEmailKey := defines.AdminEmailKey
	rds := myredis.RedisCli()
	Email, err := rds.Get(context.Background(), AdminEmailKey).Result()
	if Email == "" || err != nil {
		db := mysql.NewDB()
		var config model.SystemConfig
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "admin_email").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		// 保存
		if _, err := rds.Set(context.Background(), AdminEmailKey, Value, 7200*time.Second).Result(); err != nil {
			return
		}
		Email = Value
	}
	return
}

// GetIsOnline 查询是否是线上
func (c Common) GetIsOnline() (Online int64) {
	OnlineKey := defines.IsOnline
	rds := myredis.RedisCli()
	Online, err := rds.Get(context.Background(), OnlineKey).Int64()
	if err != nil || Online == 0 {
		db := mysql.NewDB()
		var config model.SystemConfig
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "is_online").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		// 保存
		if _, err := rds.Set(context.Background(), OnlineKey, Value, 7200*time.Second).Result(); err != nil {
			return
		}
		temp, _ := strconv.Atoi(Value)
		Online = int64(temp)
	}
	return
}

func (c Common) SecondToDHMS(Second int64) (int64, int64, int64, int64) {
	d := time.Duration(Second) * time.Second
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return int64(h / 24), int64(h % 24), int64(m), int64(s)
}

// GetIsPushTestDeviceDataConfig 获取是否推送测试数据开关配置
func (c Common) GetIsPushTestDeviceDataConfig() (IsPushTestData int64) {
	IsPushTestDataKey := defines.IsPushTestData
	rds := myredis.RedisCli()
	IsPushTestData, err := rds.Get(context.Background(), IsPushTestDataKey).Int64()
	if err != nil || IsPushTestData == 0 {
		db := mysql.NewDB()
		var config model.SystemConfig
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "is_push_test_data").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		// 保存
		if _, err := rds.Set(context.Background(), IsPushTestDataKey, Value, 7200*time.Second).Result(); err != nil {
			return
		}
		temp, _ := strconv.Atoi(Value)
		IsPushTestData = int64(temp)
	}
	return
}

func (c Common) GetGaodeConfig() (Key string) {
	GaodeKey := defines.GaodeKey
	rds := myredis.RedisCli()
	Key, err := rds.Get(context.Background(), GaodeKey).Result()
	if err != nil || Key == "" {
		db := mysql.NewDB()
		var config model.SystemConfig
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "gaode").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		// 保存
		if _, err := rds.Set(context.Background(), GaodeKey, Value, 7200*time.Second).Result(); err != nil {
			return
		}
		Key = Value
	}
	return
}

// GetNbiDeviceConfig 获取设备中台配置
func (c Common) GetNbiDeviceConfig() (Config DeviceAdminConfig) {
	deviceAdminKey := defines.DeviceAdminConfig
	rds := myredis.RedisCli()
	DeviceAdminMap, err := rds.HGetAll(context.Background(), deviceAdminKey).Result()
	if err != nil || len(DeviceAdminMap) == 0 {
		db := mysql.NewDB()
		var config model.SystemConfig
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "nbi_device_service").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		json.Unmarshal([]byte(Value), &Config)
		// 保存
		if _, err := rds.HMSet(context.Background(), deviceAdminKey,
			"is_go_device_admin", Config.IsGoDeviceAdmin,
			"device_admin_tcp_url", Config.DeviceAdminTcpURL,
			"username", Config.Username,
			"password", Config.Password,
			"type", Config.Type,
			"timeout", Config.Timeout,
		).Result(); err != nil {
			return
		}
	} else {
		Config.IsGoDeviceAdmin = DeviceAdminMap["is_go_device_admin"]
		Config.DeviceAdminTcpURL = DeviceAdminMap["device_admin_tcp_url"]
		Config.Username = DeviceAdminMap["username"]
		Config.Password = DeviceAdminMap["password"]
		Config.Type = DeviceAdminMap["type"]
		Config.Timeout, _ = strconv.Atoi(DeviceAdminMap["timeout"])
	}
	return
}

// GetPushTestDeviceDataConfig 获取推送模拟数据配置
func (c Common) GetPushTestDeviceDataConfig() (Config PushTestDeviceDataConfig) {
	PushTestDataKey := defines.PushTestDataKey
	rds := myredis.RedisCli()
	PushDeviceDataMap, err := rds.HGetAll(context.Background(), PushTestDataKey).Result()
	if err != nil || len(PushDeviceDataMap) == 0 {
		db := mysql.NewDB()
		var config model.SystemConfig
		db.Model(&model.SystemConfig{}).Where("`key` = ? and deleted_time = 0", "push_test_device_data_config").First(&config)
		if config.ID == "" {
			return
		}
		Value := config.Value
		json.Unmarshal([]byte(Value), &Config)
		// 保存
		if _, err := rds.HMSet(context.Background(), PushTestDataKey,
			"domain", Config.Domain,
			"company_id", Config.CompanyID,
		).Result(); err != nil {
			return
		}
	} else {
		Config.Domain = PushDeviceDataMap["domain"]
		Config.CompanyID = PushDeviceDataMap["company_id"]
	}
	return
}

// 根据省市区获取经纬度
func (c Common) GetLocationByArea(province, city, district string) (Location string, err error) {
	Gaode := c.GetGaodeConfig()
	if Gaode == "" {
		err = errorx.New(400018, "没有配置高德地图key")
		return
	}
	RequestUrl := "https://restapi.amap.com/v3/geocode/geo?address=?address=" + province + city + district + "&&output=JSON&key=" + Gaode
	res, err := helper.HTTPGet(context.Background(), RequestUrl, nil, 0)
	type Geocodes struct {
		Location string `json:"location"`
	}
	type GaodeRes struct {
		Geocodes []Geocodes `json:"geocodes"`
	}
	var GaodeResData GaodeRes
	err = json.Unmarshal(res, &GaodeResData)
	Location = string(GaodeResData.Geocodes[0].Location)
	return
}

func (c Common) ReverseStrings(originalSlice []string) []string {
	// 如果切片长度为0或1，则不需要翻转
	if len(originalSlice) <= 1 {
		return originalSlice
	}

	// 创建一个空的切片来存储反转后的结果
	reversedSlice := make([]string, len(originalSlice))

	// 遍历原始切片，并将元素以相反的顺序添加到新切片中
	for i, v := range originalSlice {
		reversedIndex := len(originalSlice) - 1 - i // 计算反转后的索引
		reversedSlice[reversedIndex] = v
	}
	return reversedSlice
}

// DoExport 执行导出
func (c Common) DoExport(co *gin.Context, Header []string, data [][]string, FileName string) {
	// 创建一个新的工作簿
	file := xlsx.NewFile()
	// 添加一个工作表
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		panic(err)
	}

	// 添加表头
	row := sheet.AddRow()
	for _, v := range Header {
		row.AddCell().Value = v
	}
	for _, d := range data {
		row := sheet.AddRow()
		for _, v := range d {
			cell := row.AddCell()
			cell.Value = v
		}
	}

	// 设置HTTP响应头，指定文件名和Content-Type
	co.Header("Content-Disposition", "attachment; filename="+FileName)
	co.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	co.Header("FileName", FileName)

	// 写入HTTP响应体
	err = file.Write(co.Writer)
	if err != nil {
		panic(err)
	}
}

// DoExports 执行导出两个sheet
func (c Common) DoExports(co *gin.Context, Header map[string][]string, data map[string][][]string, FileName map[string]string) {
	// 创建一个新的工作簿
	file := xlsx.NewFile()
	// 添加一个工作表
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		panic(err)
	}
	sheet1, err := file.AddSheet("Sheet2")
	if err != nil {
		panic(err)
	}
	// 添加表头
	row := sheet.AddRow()
	for _, v := range Header["metric"] {
		row.AddCell().Value = v
	}
	for _, d := range data["metric"] {
		row := sheet.AddRow()
		for _, v := range d {
			cell := row.AddCell()
			cell.Value = v
		}
	}
	// 添加表头
	row1 := sheet1.AddRow()
	for _, v := range Header["check"] {
		row1.AddCell().Value = v
	}
	for _, d := range data["check"] {
		row1 := sheet1.AddRow()
		for _, v := range d {
			cell1 := row1.AddCell()
			cell1.Value = v
		}
	}

	// 设置HTTP响应头，指定文件名和Content-Type
	co.Header("Content-Disposition", "attachment; filename="+FileName["check"])
	co.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// 写入HTTP响应体
	err = file.Write(co.Writer)
	if err != nil {
		panic(err)
	}
}

func (c Common) Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func (c Common) InArray(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func (c Common) InArrayInt(need int64, needArr []int64) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

// validateLatLng 经纬度校验
func (c Common) ValidateLatLng(lat, lng string) bool {
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return false
	}
	lngFloat, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		return false
	}
	if latFloat < -90 || latFloat > 90 {
		return false
	}
	if lngFloat < -180 || lngFloat > 180 {
		return false
	}
	return true
}

// ValidateCommaPosition 校验字符串中是否包含逗号
func (c Common) ValidateCommaPosition(s string) bool {
	index := strings.Index(s, ",")
	// 如果字符串中不包含逗号，则Index函数返回-1
	if index == -1 {
		return false
	}
	// 如果逗号在字符串的首尾位置，则认为是无效的
	if index == 0 || index == len(s)-1 {
		return false
	}
	// 如果有多个逗号也返回错误
	if strings.Count(s, ",") > 1 {
		return false
	}
	return true
}

// HTTPGet 发送get 请求 timeout 单位 秒
func (c Common) HTTPGet(url string, header map[string]string, timeout time.Duration) (res []byte, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "application/json")
	for key, value := range header {
		request.Header.Set(key, value)
	}
	client := &http.Client{}
	client.Timeout = time.Second * timeout
	resp, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer func() { _ = resp.Body.Close() }()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	return respData, nil
}

func (c Common) PostJSON(url string, data interface{}, header map[string]string, timeout time.Duration) (res []byte, err error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return res, err
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "application/json")
	for key, value := range header {
		request.Header.Set(key, value)
	}
	client := &http.Client{}
	client.Timeout = time.Second * timeout
	resp, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer func() { _ = resp.Body.Close() }()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	return respData, nil
}

// HasDuplicate 判断是否有重复值
func (c Common) HasDuplicate(s []string) bool {
	m := make(map[string]bool)
	for _, v := range s {
		if m[v] {
			return true
		}
		m[v] = true
	}
	return false
}

// HasIntersection 判断两个[]string是否有交集
func (c Common) HasIntersection(s1, s2 []string) bool {
	m := make(map[string]bool)
	for _, v := range s1 {
		m[v] = true
	}
	for _, v := range s2 {
		if m[v] {
			return true
		}
	}
	return false
}

// hasDuplicateNumbers
func (c Common) HasDuplicateNumbers(Strings []string) bool {
	seen := make(map[string]bool)
	for _, str := range Strings {
		nums := strings.Split(str, ",")
		for _, num := range nums {
			if seen[num] {
				return true
			}
			seen[num] = true
		}
	}
	return false
}

// CreateDirIfNotExist 路径不存在则创建
func (c Common) CreateDirIfNotExist(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

// MergeStringSlices 多个切片合并
func (c Common) MergeStringSlices(slices ...[]string) []string {
	var merged []string
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	return merged
}

// 均匀获取[]int64中的n个点
func (c Common) GetAvgElements(input []int64, Num int) []int64 {
	n := len(input)
	step := n / Num

	result := make([]int64, 0, Num)

	for i := 0; i < n; i += step {
		result = append(result, input[i])
		if len(result) == Num {
			break
		}
	}
	return result
}

// 均匀获取[]int64中的n个点，保留收尾两个点
func (c Common) GetEvenlyPoints(numbers []int64, n int) []int64 {
	length := len(numbers)
	if n >= length {
		return numbers
	}
	result := make([]int64, n)
	interval := (length - 2) / (n - 1) // 计算均匀间隔
	result[0] = numbers[0]             // 保留开头点
	result[n-1] = numbers[length-1]    // 保留结尾点
	for i := 1; i < n-1; i++ {
		index := i * interval
		result[i] = numbers[index]
	}
	return result
}

// RemoveDuplicates 去重
func (c Common) RemoveDuplicates(slice []string) []string {
	encountered := map[string]bool{} // 用于记录已经遇到的字符串
	result := []string{}             // 用于存储去重后的字符串
	for _, str := range slice {
		if !encountered[str] { // 如果字符串还没有被遇到过
			encountered[str] = true      // 将其标记为已遇到
			result = append(result, str) // 添加到结果列表中
		}
	}
	return result
}

func (c Common) RemoveIntDuplicates(slice []int8) []int8 {
	encountered := map[int8]bool{} // 用于记录已经遇到的字符串
	result := []int8{}             // 用于存储去重后的字符串
	for _, str := range slice {
		if !encountered[str] { // 如果字符串还没有被遇到过
			encountered[str] = true      // 将其标记为已遇到
			result = append(result, str) // 添加到结果列表中
		}
	}
	return result
}

// StringSliceDifference 取差集
func (c Common) StringSliceDifference(slice1, slice2 []string) []string {
	set := make(map[string]bool)

	// 将第一个切片的元素添加到集合中
	for _, str := range slice1 {
		set[str] = true
	}

	// 从集合中删除第二个切片中的元素
	for _, str := range slice2 {
		delete(set, str)
	}

	// 构建差集切片
	difference := make([]string, 0, len(set))
	for str := range set {
		difference = append(difference, str)
	}

	return difference
}

// Intersect 取交集
func (c Common) Intersect(a, b []string) []string {
	m := make(map[string]bool)
	var result []string
	// 将第一个切片的元素添加到 map 中
	for _, value := range a {
		m[value] = true
	}
	// 遍历第二个切片，如果元素存在于 map 中，则为交集元素
	for _, value := range b {
		if m[value] {
			result = append(result, value)
		}
	}
	return result
}

// RemoveString 删除某个元素
func (c Common) RemoveString(slice []string, value string) []string {
	result := []string{}
	for _, item := range slice {
		if item != value {
			result = append(result, item)
		}
	}
	return result
}

// 校验ip
func (c Common) IsValidIP(ip string) bool {
	// 使用正则表达式验证IP的合法性
	pattern := `^([0-9]{1,3}\.){3}[0-9]{1,3}$`
	match, _ := regexp.MatchString(pattern, ip)
	if !match {
		return false
	}
	// 使用net.ParseIP验证IP的合法性
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return false
	}

	return true
}

// Contain 判断元素是否在切片中
func (c Common) Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, nil
}

type AddressComponent struct {
	City         []string `json:"city"`
	Province     string   `json:"province"`
	Adcode       string   `json:"adcode"`
	District     string   `json:"district"`
	Towncode     string   `json:"towncode"`
	StreetNumber struct {
		Number    string `json:"number"`
		Location  string `json:"location"`
		Direction string `json:"direction"`
		Distance  string `json:"distance"`
		Street    string `json:"street"`
	} `json:"streetNumber"`
	Township      string `json:"township"`
	BusinessAreas []struct {
		Location string `json:"location"`
		Name     string `json:"name"`
		ID       string `json:"id"`
	} `json:"businessAreas"`
}

type Regeocode struct {
	AddressComponent AddressComponent `json:"addressComponent"`
}

type GeocodeResult struct {
	Status    string    `json:"status"`
	Regeocode Regeocode `json:"regeocode"`
}

// LatToAddress 经纬度转地址
func (c Common) LatToAddress(LatLng string) (Address string, err error) {
	Gaode := c.GetGaodeConfig()
	if Gaode == "" {
		err = errorx.New(400018, "没有配置高德地图key")
		return
	}
	apiUrl := "https://restapi.amap.com/v3/geocode/regeo"
	params := map[string]string{
		"key":      Gaode,
		"location": LatLng,
	}
	url := apiUrl + "?" + c.UrlParams(params)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var result GeocodeResult
	err = json.Unmarshal(body, &result)

	if result.Status == "1" {
		province := result.Regeocode.AddressComponent.Province
		city := result.Regeocode.AddressComponent.District
		district := result.Regeocode.AddressComponent.Township
		Address = province + "/" + city + "/" + district
	} else {
		err = errorx.New(500124, "Failed to get address")
		return
	}
	return
}

func (c Common) UrlParams(params map[string]string) string {
	pairs := []string{}
	for k, v := range params {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, "&")
}

func (c Common) CapitalizeFirstLetter(s string) string {
	if len(s) >= 1 {
		return strings.ToUpper(string(s[0])) + s[1:]
	}
	return s
}

// SendQQEmail 发送QQ邮件
func (c Common) SendQQEmail(Type, Title, WarningContent string, Emails []string, CompanyID string) {
	// 邮件配置
	EmailConfig := AppInit.AppConfig.EmailConfig
	smtpHost := EmailConfig.SmtpHost
	smtpPort := EmailConfig.SmtpPort
	senderEmail := EmailConfig.SenderEmail
	senderPassword := EmailConfig.SenderPassword

	// 收件人列表
	recipientEmails := Emails

	// 创建邮件内容
	subject := Title
	body := WarningContent
	if Type == "" {
		Type = "text/plain"
	}

	// 初始化邮件对象
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("Subject", subject)
	m.SetBody(Type, body)

	// 设置 SMTP 邮件服务器配置
	do := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	// 批量发送邮件
	for _, recipient := range recipientEmails {
		status := 1
		m.SetHeader("To", recipient)
		// 发送邮件
		err := do.DialAndSend(m)
		if err != nil {
			status = 2
			continue
		}
		var log schema.SendMsgLogSaveParams
		log.CompanyID = CompanyID
		log.SendType = 1
		log.Target = recipient
		log.Content = WarningContent
		log.Status = int8(status)
		res, _ := json.Marshal(err)
		log.Res = string(res)
		c.SendMsgLogSave(log)
	}
}

func (c Common) SendEmail(Title string, Emails []string, CompanyID string) {
	// 邮件配置
	EmailConfig := c.GetEmailConfig()
	// EmailConfig := AppInit.AppConfig.EmailConfig
	smtpHost := EmailConfig.SmtpHost
	smtpPort, _ := strconv.Atoi(EmailConfig.SmtpPort)
	senderEmail := EmailConfig.SenderEmail
	senderPassword := EmailConfig.SenderPassword

	// 收件人列表
	recipientEmails := Emails
	Code, _ := c.GenerateNRandomDecimal(6)
	// 创建邮件内容
	subject := Title
	// body := "登录验证码为【<span style=color:'red;'>" + Code + "</span>】"
	body := "登录验证码为【" + Code + "】"

	// 初始化邮件对象
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// 设置 SMTP 邮件服务器配置
	do := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)
	rdb := myredis.RedisCli()
	// 批量发送邮件
	for _, recipient := range recipientEmails {
		status := 1
		m.SetHeader("To", recipient)
		// 发送邮件
		err := do.DialAndSend(m)
		if err != nil {
			status = 2
			continue
		}
		var log schema.SendMsgLogSaveParams
		log.CompanyID = CompanyID
		log.SendType = 1
		log.Target = recipient
		log.Content = body
		log.Status = int8(status)
		res, _ := json.Marshal(err)
		log.Res = string(res)
		c.SendMsgLogSave(log)
		// 将邮箱和邮箱验证码存入redis
		EmailKey := defines.EmailLogin + recipient
		rdb.Set(context.Background(), EmailKey, Code, 300*time.Second).Result() // 5分钟有效
	}
}

// SendSms 短信发送
func (c Common) SendSms(WarningContent string, Telephones []string, CompanyID string) {
	// 阿里云短信服务配置
	SmsConfig := c.GetSmsConfig()
	domain := "dysmsapi.aliyuncs.com"
	// SmsConfig := AppInit.AppConfig.SmsConfig
	accessKeyId := SmsConfig.AccessKeyId
	accessKeySecret := SmsConfig.AccessKeySecret
	regionId := SmsConfig.RegionId
	signName := SmsConfig.SignName
	templateCode := SmsConfig.TemplateCode

	// 创建短信客户端
	client, err := dysmsapi.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("创建短信客户端时出错:", err)
		return
	}
	// 创建发送短信的请求
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.Domain = domain
	request.SignName = signName
	request.TemplateCode = templateCode
	// 设置短信模板参数
	// request.TemplateParam = `{"sn":"354c060361af0104", "name":"高温预警", "code":"66.5"}`
	Code, _ := c.GenerateNRandomDecimal(6)
	request.TemplateParam = `{"code":"` + Code + `"}`
	rdb := myredis.RedisCli()
	// 发送短信
	for _, phoneNumber := range Telephones {
		status := 1
		request.PhoneNumbers = phoneNumber
		response, err := client.SendSms(request)
		if err != nil {
			status = 2
			fmt.Println("发送短信时出错:", err)
			return
		}
		var log schema.SendMsgLogSaveParams
		log.CompanyID = CompanyID
		log.SendType = 2
		log.Target = phoneNumber
		log.Content = Code
		log.Status = int8(status)
		res, _ := json.Marshal(response)
		log.Res = string(res)
		c.SendMsgLogSave(log)
		// 将邮箱和邮箱验证码存入redis
		PhoneKey := defines.PhoneLogin + phoneNumber
		rdb.Set(context.Background(), PhoneKey, Code, 300*time.Second).Result() // 5分钟有效
	}
}

// SendMsgLogSave 告警消息推送记录
func (c Common) SendMsgLogSave(Params schema.SendMsgLogSaveParams) {
	db := mysql.NewDB()
	var InsertParams model.SendMsgLog
	InsertParams.ID, _ = helper.UUID{}.GetUniqueKey()
	InsertParams.CompanyID = Params.CompanyID
	InsertParams.SendType = Params.SendType
	InsertParams.Target = Params.Target
	InsertParams.Content = Params.Content
	InsertParams.Status = Params.Status
	InsertParams.Res = Params.Res
	InsertParams.CreatedTime = time.Now().Unix()
	db.Model(&model.SendMsgLog{}).Create(&InsertParams)
}

// ZhToEn 汉字转拼音
func (c Common) ZhToEn(Str string) (PStr string) {
	// pinyin := pinyin.NewArgs()
	for _, ch := range Str {
		// Lazy := pinyin.LazyPinyin(string(ch), pinyin.NewArgs())
		arr := pinyin.LazyConvert(string(ch), nil)
		PStr += arr[0]
	}
	return
}

// 生成n位十进制随机数的函数
func (c Common) GenerateNRandomDecimal(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("number of digits must be greater than 0")
	}
	// 设置随机数种子，以当前时间为例
	rand.Seed(time.Now().UnixNano())
	// 最大的n位数（比如，对于3位数，最大的数是999）
	maxNum := int(c.Pow10(n) - 1)
	// 生成随机数
	randomNum := rand.Intn(maxNum + 1)
	// 转换为字符串
	randomStr := strconv.Itoa(randomNum)

	// 如果生成的随机数位数不足n位，前面补0
	for len(randomStr) < n {
		randomStr = "0" + randomStr
	}
	return randomStr, nil
}

func (c Common) Pow10(n int) int64 {
	result := int64(1)
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

func (c Common) GetServerIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()

	return ip, nil
}

func (c Common) FormatTimeStr(timeStr string) (timestamp int64, err error) {
	// 定义多个可能的时间格式
	formats := []string{
		"20060102",
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
		"2006-1-2",
		"2006-1-2 15:04:05",
		"2006/1/2",
		"2006/1/2 15:04:05",
		"06-01-02",
		"06/01/02",
		"06-01-02 15:04:05",
		"06/01/02 15:04:05",
	}
	var t time.Time
	// 尝试使用多个格式解析时间字符串
	for _, format := range formats {
		t, err = time.ParseInLocation(format, timeStr, time.Local)
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Println("解析时间出错:", err)
		return
	}
	// 将时间转换为时间戳（秒）
	timestamp = t.Unix()
	return
}

// AjaxReturn ajax return
func (c Common) AjaxReturn(g *gin.Context, code int64, data interface{}) {
	g.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  code,
	})
}

func (c Common) GetStartOfDay(timestamp int64) int64 {
	// 将时间戳转换为时间对象
	t := time.Unix(timestamp, 0)

	// 获取当天的起始时间
	startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	// 将起始时间转换为时间戳
	startOfDayTimestamp := startOfDay.Unix()

	return startOfDayTimestamp
}

func (c Common) GetTimePoints(start_time, end_time, time_point int64) []int64 {
	// 加载本地时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	start := time.Unix(start_time, 0).In(loc)
	end := time.Unix(end_time, 0).In(loc)

	// 存储符合条件的时间点
	var timePoints []int64

	// 循环遍历每一天
	for current := start; current.Before(end) || current.Equal(end); current = current.AddDate(0, 0, 1) {
		// 将时间点设置为指定的日期、小时和分钟
		desiredTime := time.Date(current.Year(), current.Month(), current.Day(), int(time_point), 0, 0, 0, loc)

		// 检查所需时间是否在日期范围内
		if desiredTime.After(current) && desiredTime.Before(current.AddDate(0, 0, 1)) {
			timePoints = append(timePoints, desiredTime.Unix())
		}
	}
	return timePoints
}

func (c Common) RemoveLeadingZeros(str string) string {
	num, _ := strconv.Atoi(str)
	return strconv.Itoa(num)
}

func (c Common) ConvertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unable to convert value to float64")
	}
}

func (c Common) ContainsUpperCaseOrDigit(str string) bool {
	for _, char := range str {
		if unicode.IsUpper(char) || unicode.IsDigit(char) {
			return true
		}
	}
	return false
}
