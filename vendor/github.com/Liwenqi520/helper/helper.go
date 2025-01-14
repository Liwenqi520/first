package helper

import (
	// "brewing-sass/gopkg/logger"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	mathRand "math/rand"
	netDefault "net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/go-gomail/gomail"
	// "github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/syyongx/php2go"
	"golang.org/x/crypto/bcrypt"
)

//加密
func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

//验证
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}

/* --------------  验证  ----------------*/
// CheckMobile 检查手机号
func CheckMobile(mobile string, lang string) bool {
	var regString string
	if lang == "" || lang == "cn" {
		regString = `^1[3-9][0-9]\d{8}$`
	} else if lang == "jp" {
		regString = `^\d{11}$`
	} else if lang == "en" {
		regString = `^\d{10}$`
	}
	reg := regexp.MustCompile(regString)
	return reg.MatchString(mobile)
}

//CheckTelPhone 检查电话号码
func CheckTelPhone(phone string) bool {
	reg := regexp.MustCompile(`(^(0[0-9]{2,3}-)?([2-9][0-9]{6,7})+(-[0-9]{1,4})?$)|(^((\(\d{3}\))|(\d{3}-))?(1[3578]\d{9})$)|(^(400)-(\d{3})-(\d{4})(.)(\d{1,4})$)|(^(400)-(\d{3})-(\d{4}$))`)
	return reg.MatchString(phone)
}

// CheckEmail 检查电子邮箱
func CheckEmail(email string) bool {
	reg := regexp.MustCompile(`^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`)
	return reg.MatchString(email)
}

//CheckPassWord 检查密码 无特殊字符 无长度要求
func CheckPassWord(password string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9]{1,}$`)
	return reg.MatchString(password)
}

//MatchString 正则检查字符串s是否在正则str里面
func MatchString(str, s string) bool {
	return regexp.MustCompile(str).MatchString(s)
}

//IsMoney 检查是否是Money
func IsMoney(f float64) bool {
	//reg := `^[^-][\d]*\.([\d]{2}|[\d]{1}|[\d]{0})$`
	reg1 := `^[^-][\d]*(\.|)([\d]{0,2})$` //非负数 不超过2位小数
	return regexp.MustCompile(reg1).MatchString(strconv.FormatFloat(f, 'f', -1, 64))
}

// Checkdate 检验年月日是否存在
func Checkdate(month, day, year int) bool {
	if month < 1 || month > 12 || day < 1 || day > 31 || year < 1 || year > 32767 {
		return false
	}
	switch month {
	case 4, 6, 9, 11:
		if day > 30 {
			return false
		}
	case 2:
		// leap year
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			if day > 29 {
				return false
			}
		} else if day > 28 {
			return false
		}
	}
	return true
}

// CheckChineseChar 校验是否中文字符
// str 输入字符串
// checkType 1-纯中文 2-包含中文
func CheckChineseChar(str string, checkType int) bool {
	count := 0
	for _, v := range str {
		if unicode.Is(unicode.Scripts["Han"], v) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(v))) {
			count++
		}
	}
	//
	if checkType == 1 {
		strLen := utf8.RuneCountInString(str)
		if count == strLen {
			return true
		}
	} else if checkType == 2 {
		if count > 0 {
			return true
		}
	}
	return false
}

// CheckURL 校验合法网址
// str 输入字符串
// checkType 1-校验http[s]://www.aaa.com类型的网址 2-校验局域网IP是否合法的网址
func CheckURL(str string, checkType int) bool {
	if !(checkType == 1 || checkType == 2) {
		return false
	}
	regStr := ""
	if checkType == 1 {
		// 校验http[s]://www.aaa.com类型的网址
		regStr = `^((https|http):\/\/)?([\w-]+\.)+[\w-]+(/[\w-./?%&=#]*)?$`
	} else if checkType == 2 {
		// 校验局域网IP是否合法的网址
		regStr = `^([hH][tT]{2}[pP]:\/\/|[hH][tT]{2}[pP][sS]:\/\/)((2[0-4][0-9])|(25[0-5])|(1[0-9]{0,2})|([1-9][0-9])|([1-9]))\.(((2[0-4][0-9])|(25[0-5])|(1[0-9]{0,2})|([1-9][0-9])|([0-9]))\.){2}((2[0-4][0-9])|(25[0-5])|(1[0-9]{0,2})|([1-9][0-9])|([1-9]))$`
	}
	reg := regexp.MustCompile(regStr)
	return reg.MatchString(str)
}

// CheckCertificate 校验证件
// str 输入字符串
// checkType 1-身份证 2-护照 3-港澳通行证
func CheckCertificate(str string, checkType int) bool {
	if !(checkType == 1 || checkType == 2 || checkType == 3) {
		return false
	}
	regStr := ""
	if checkType == 1 { // 身份证
		regStr = `^(^\d{18}$|^\d{17}(\d|X|x))$`
	} else if checkType == 2 { // 护照
		regStr = `^1[45][0-9]{7}$|([P|p|S|s]\d{7}$)|([S|s|G|g]\d{8}$)|([Gg|Tt|Ss|Ll|Qq|Dd|Aa|Ff]\d{8}$)|([H|h|M|m]\d{8,10})$`
	} else if checkType == 3 { // 港澳通行证
		regStr = `^([A-Z]\d{6,10}(\w1)?)$`
	}
	reg := regexp.MustCompile(regStr)
	return reg.MatchString(str)
}

/*------------------------转换-----------------------------*/
func StructToMap(content interface{}) (data map[string]interface{}, err error) {
	var name map[string]interface{}
	if marshalContent, err := json.Marshal(content); err != nil {
		return nil, err
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber() // 设置将float64转为一个number
		if err := d.Decode(&name); err != nil {
			return nil, err
		} else {
			for k, v := range name {
				name[k] = v
			}
		}
	}
	return name, nil
}

/*-------------------------加解密--------------------------------*/
//Md5 MD5加密
func Md5(str string) string {
	plain := md5.New()
	plain.Write([]byte(str))
	return hex.EncodeToString(plain.Sum(nil))
}

//Base64Encode base64加密
func Base64Encode(str string) string {
	encodeString := base64.StdEncoding.EncodeToString([]byte(str))
	return encodeString
}

//Base64Decode base64解密
func Base64Decode(str string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatalln(err)
	}
	return string(decodeBytes)
}

//Base64EncodeInURL base64加密 如果要用在url中，需要使用URLEncoding
func Base64EncodeInURL(str string) string {
	uEnc := base64.URLEncoding.EncodeToString([]byte(str))
	return uEnc
}

//Base64DecodeWithInURL base64解密
func Base64DecodeWithInURL(str string) string {
	uDec, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		log.Fatalln(err)
	}
	return string(uDec)
}

// // GenerateRSAKey 生成RSA私钥和公钥，保存到文件中
// func GenerateRSAKey(bits int, savePathPublic string, savePathPrivate string) {
// 	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
// 	//Reader是一个全局、共享的密码用强随机数生成器
// 	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//保存私钥
// 	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
// 	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
// 	//使用pem格式对x509输出的内容进行编码
// 	//创建文件保存私钥
// 	privateFile, err := os.Create(path.Join(savePathPrivate, "private.pem"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer func() { _ = privateFile.Close() }()
// 	//构建一个pem.Block结构体对象
// 	privateBlock := pem.Block{Type: "PRIVATE KEY", Bytes: X509PrivateKey}
// 	//将数据保存到文件
// 	pem.Encode(privateFile, &privateBlock)
// 	//保存公钥
// 	//获取公钥的数据
// 	publicKey := privateKey.PublicKey
// 	//X509对公钥编码
// 	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//pem格式编码
// 	//创建用于保存公钥的文件
// 	publicFile, err := os.Create(path.Join(savePathPublic, "public.pem"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer func() { _ = publicFile.Close() }()
// 	//创建一个pem.Block结构体对象
// 	publicBlock := pem.Block{Type: "PUBLIC KEY", Bytes: X509PublicKey}
// 	//保存到文件
// 	pem.Encode(publicFile, &publicBlock)
// }
// // RSAEncrypt RSA加密
// func RSAEncrypt(plainText []byte, path string) []byte {
// 	//打开文件
// 	file, err := os.Open(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer func() { _ = file.Close() }()
// 	//读取文件的内容
// 	info, _ := file.Stat()
// 	buf := make([]byte, info.Size())
// 	_, err = file.Read(buf)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//pem解码
// 	block, _ := pem.Decode(buf)
// 	//x509解码
// 	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//类型断言
// 	publicKey := publicKeyInterface.(*rsa.PublicKey)
// 	//对明文进行加密
// 	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//返回密文
// 	return cipherText
// }
// // RSADecrypt RSA解密
// func RSADecrypt(cipherText []byte, path string) []byte {
// 	//打开文件
// 	file, err := os.Open(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer func() { _ = file.Close() }()
// 	//获取文件内容
// 	info, _ := file.Stat()
// 	buf := make([]byte, info.Size())
// 	_, err = file.Read(buf)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//pem解码
// 	block, _ := pem.Decode(buf)
// 	//X509解码
// 	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//对密文进行解密
// 	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
// 	//返回明文
// 	return plainText
// }
// pKCS7Padding pKCS7Padding
func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// pKCS7UnPadding pKCS7UnPadding
func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS7Padding(origData, blockSize) // CBC加解密。明文的长度不一定总是128的倍数，采用PKCS7填充方式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt AES解密
// key 对称秘钥长度必须是16的倍数
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}

/*--------------------------邮箱----------------------*/
type EmailConfig struct {
	ServerHost string `mapstructure:"server_host" json:"server_host"`     // ServerHost 邮箱服务器地址，如腾讯企业邮箱为smtp.exmail.qq.com
	ServerPort int    `mapstructure:"server_port" json:"server_port"`     // ServerPort 邮箱服务器端口，如腾讯企业邮箱为465
	FromEmail  string `mapstructure:"from_email" json:"from_email"`       // FromEmail　发件人邮箱地址
	FromPasswd string `mapstructure:"from_password" json:"from_password"` //发件人邮箱密码（注意，这里是明文形式)
}
type EmailInfo struct {
	Recipient []string //收件人邮箱
	CC        []string //抄送
}

var emailMessage *gomail.Message

/**
 * @Description: 发送邮件
 * @Param : subject[主题]、body[内容]、emailInfo[发邮箱需要的信息(参考EmailInfo)]
 * @Return:
 */
func SendEmail(subject, body string, emailInfo *EmailInfo, emailConf EmailConfig) (err error) {
	if len(emailInfo.Recipient) == 0 {
		log.Print("收件人列表为空")
		return
	}
	emailMessage = gomail.NewMessage()
	//设置收件人
	emailMessage.SetHeader("To", emailInfo.Recipient...)
	//设置抄送列表
	if len(emailInfo.CC) != 0 {
		emailMessage.SetHeader("Cc", emailInfo.CC...)
	}
	// 第三个参数为发件人别名，如"dcj"，可以为空（此时则为邮箱名称）
	emailMessage.SetAddressHeader("From", emailConf.FromEmail, "dcj")
	//主题
	emailMessage.SetHeader("Subject", subject)
	//正文
	emailMessage.SetBody("text/html", body)
	d := gomail.NewDialer(emailConf.ServerHost, emailConf.ServerPort,
		emailConf.FromEmail, emailConf.FromPasswd)
	err = d.DialAndSend(emailMessage)
	if err != nil {
		log.Println("发送邮件失败: ", err)
	} else {
		log.Println("已成功发送邮件到指定邮箱")
	}
	return
}

/*-----------------文件处理--------------*/
// FileExists file_exists()
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

//IsModeDir 是否为文件夹
func IsModeDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	fm := fd.Mode()
	return fm.IsDir(), nil
}

//ExistsDirAndMkdir 验证文件夹是否存在并创建
func ExistsDirAndMkdir(pathDir string) bool {
	exist, err := PathExists(pathDir)
	if err != nil {
		return false
	}
	if exist {
		return true //如果有 直接返回true
		//has dir!
	}
	//no dir!
	// 创建文件夹
	err1 := os.MkdirAll(pathDir, os.ModePerm)
	return err1 == nil
}

//PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//GetCwd 获取到当前目录的完整路径
func GetCwd() (string, error) {
	dir, err := os.Getwd()
	return dir, err
}

//获取目录下的文件或目录，不包含子级
func GetFileAndDir(filePath string) (FileDirList []string, err error) {
	if filePath == "" {
		filePath = "./"
	}
	pwdClean := filepath.Clean(filePath)
	fmt.Println("pwd:", filePath)
	fmt.Println("pwdClean:", pwdClean)
	//获取当前目录下的所有文件或目录信息
	filepathNames, err := filepath.Glob(filepath.Join(pwdClean, "*"))
	if err != nil {
		log.Fatal(err)
	}
	for i := range filepathNames {
		fmt.Println(filepathNames[i]) //打印path
		FileDirList = append(FileDirList, filepathNames[i])
	}
	return
}

//获取目录下的文件或目录，包含子级
func GetFileAndDirRecursion(filePath string) (FileDirList []string, err error) {
	if filePath == "" {
		filePath = "./"
	}
	pwdClean := filepath.Clean(filePath)
	fmt.Println("pwd:", filePath)
	fmt.Println("pwdClean:", pwdClean)
	//获取当前目录下的所有文件或目录信息
	filepath.Walk(pwdClean, func(path string, info os.FileInfo, err error) error {
		FileDirList = append(FileDirList, info.Name())
		return nil
	})
	return
}

/*-------------------html处理---------------------*/
//TrimSpecialTags 去除危险标签
func TrimSpecialTags(str string) (safeStr string) {
	//去除iframe
	re, _ := regexp.Compile(`<iframe[\S\s]+?</iframe>`)
	src := re.ReplaceAllString(str, "")
	return src
}

//GetPlainTextFromHTML html去除标签获取纯文本，支持截断
func GetPlainTextFromHTML(str string, splitLength int) (result string) {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile(`<[\S\s]+?>`)
	str = re.ReplaceAllStringFunc(str, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile(`<style[\S\s]+?</style>`)
	str = re.ReplaceAllString(str, "")
	//去除SCRIPT
	re, _ = regexp.Compile(`<script[\S\s]+?</script>`)
	str = re.ReplaceAllString(str, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile(`<[\S\s]+?>`)
	str = re.ReplaceAllString(str, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile(`\s{2,}`)
	str = re.ReplaceAllString(str, "\n")
	//去除空格
	str = strings.TrimSpace(str)
	if splitLength != 0 {
		maxSplitLength := int(math.Min(float64(splitLength), float64(utf8.RuneCountInString(str))))
		return string([]rune(str)[:maxSplitLength])
	}
	return strings.TrimSpace(str)
}

/*------------------坐标定位------------------*/
// WGS84坐标系：即地球坐标系，国际上通用的坐标系。
// GCJ02坐标系：即火星坐标系，WGS84坐标系经加密后的坐标系。Google Maps，高德在用。
// BD09坐标系：即百度坐标系，GCJ02坐标系经加密后的坐标系。
const (
	X_PI   = math.Pi * 3000.0 / 180.0
	OFFSET = 0.00669342162296594323
	AXIS   = 6378245.0
)

//BD09toGCJ02 百度坐标系->火星坐标系
func BD09toGCJ02(lon, lat float64) (float64, float64) {
	x := lon - 0.0065
	y := lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*X_PI)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*X_PI)
	gLon := z * math.Cos(theta)
	gLat := z * math.Sin(theta)
	return gLon, gLat
}

//GCJ02toBD09 火星坐标系->百度坐标系
func GCJ02toBD09(lon, lat float64) (float64, float64) {
	z := math.Sqrt(lon*lon+lat*lat) + 0.00002*math.Sin(lat*X_PI)
	theta := math.Atan2(lat, lon) + 0.000003*math.Cos(lon*X_PI)
	bdLon := z*math.Cos(theta) + 0.0065
	bdLat := z*math.Sin(theta) + 0.006
	return bdLon, bdLat
}

//WGS84toGCJ02 WGS84坐标系->火星坐标系
func WGS84toGCJ02(lon, lat float64) (float64, float64) {
	if isOutOFChina(lon, lat) {
		return lon, lat
	}
	mgLon, mgLat := delta(lon, lat)
	return mgLon, mgLat
}

//GCJ02toWGS84 火星坐标系->WGS84坐标系
func GCJ02toWGS84(lon, lat float64) (float64, float64) {
	if isOutOFChina(lon, lat) {
		return lon, lat
	}
	mgLon, mgLat := delta(lon, lat)
	return lon*2 - mgLon, lat*2 - mgLat
}

//BD09toWGS84 百度坐标系->WGS84坐标系
func BD09toWGS84(lon, lat float64) (float64, float64) {
	lon, lat = BD09toGCJ02(lon, lat)
	return GCJ02toWGS84(lon, lat)
}

//WGS84toBD09 WGS84坐标系->百度坐标系
func WGS84toBD09(lon, lat float64) (float64, float64) {
	lon, lat = WGS84toGCJ02(lon, lat)
	return GCJ02toBD09(lon, lat)
}
func delta(lon, lat float64) (float64, float64) {
	dlat := transformlat(lon-105.0, lat-35.0)
	dlon := transformlng(lon-105.0, lat-35.0)
	radlat := lat / 180.0 * math.Pi
	magic := math.Sin(radlat)
	magic = 1 - OFFSET*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((AXIS * (1 - OFFSET)) / (magic * sqrtmagic) * math.Pi)
	dlon = (dlon * 180.0) / (AXIS / sqrtmagic * math.Cos(radlat) * math.Pi)
	mgLat := lat + dlat
	mgLon := lon + dlon
	return mgLon, mgLat
}
func transformlat(lon, lat float64) float64 {
	var ret = -100.0 + 2.0*lon + 3.0*lat + 0.2*lat*lat + 0.1*lon*lat + 0.2*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*math.Pi) + 40.0*math.Sin(lat/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*math.Pi) + 320*math.Sin(lat*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}
func transformlng(lon, lat float64) float64 {
	var ret = 300.0 + lon + 2.0*lat + 0.1*lon*lon + 0.1*lon*lat + 0.1*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lon*math.Pi) + 40.0*math.Sin(lon/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lon/12.0*math.Pi) + 300.0*math.Sin(lon/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}
func isOutOFChina(lon, lat float64) bool {
	return !(lon > 73.66 && lon < 135.05 && lat > 3.86 && lat < 53.55)
}

/*------------数学运算------------*/
//EarthDistance 计算两地距离，返回值 X km
func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6378.137
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}

//RandInt64 取两个数之间的随机数
func RandInt64(min, max int64) int64 {
	//仅支持非负整数，且两端不能同时等于0，否则返回0
	if min < 0 || max < 0 || (min == 0 && max == 0) {
		return 0
	}
	if min >= max {
		return mathRand.Int63n(min-max) + max
	}
	return mathRand.Int63n(max-min) + min
}

//浮点运算
func BasicOperation(express string) (result float64, err error) {
	//	将中缀表达式转换成后缀表达式（逆波兰式），postfixExpress：后缀表达式
	express = strings.ReplaceAll(express, " ", "")
	postfixExpress, errRet := transPostfixExpress(express)
	if errRet != nil {
		fmt.Println(errRet)
		return
	}
	//	后缀表达式求值
	result, errRet = calc(postfixExpress)
	if errRet != nil {
		fmt.Println("error:", errRet)
		return
	}
	return result, nil
}
func priority(s byte) int {
	switch s {
	case '+':
		return 1
	case '-':
		return 1
	case '*':
		return 2
	case '/':
		return 2
	}
	return 0
}

//	将中缀表达式转换成后缀表达式（逆波兰式），postfixExpress：后缀表达式
func transPostfixExpress(express string) (postfixExpress []string, err error) {
	var (
		opStack Stack //	运算符堆栈
		i       int
	)
LABEL:
	for i < len(express) { //	从左至右扫描中缀表达式
		switch {
		//	1. 若读取的是操作数，则将该操作数存入后缀表达式。
		case (express[i] >= '0' && express[i] <= '9') || express[i] == '.':
			var number []byte //	如数字123，由'1'、'2'、'3'组成
			for ; i < len(express); i++ {
				if (express[i] < '0' || express[i] > '9') && express[i] != '.' {
					break
				}
				number = append(number, express[i])
			}
			postfixExpress = append(postfixExpress, string(number))
		//	2. 若读取的是运算符：
		//	(1) 该运算符为左括号"("，则直接压入运算符堆栈。
		case express[i] == '(':
			opStack.Push(fmt.Sprintf("%c", express[i]))
			i++
		//	(2) 该运算符为右括号")"，则输出运算符堆栈中的运算符到后缀表达式，直到遇到左括号为止。
		case express[i] == ')':
			for !opStack.IsEmpty() {
				data, _ := opStack.Pop()
				if data[0] == '(' {
					break
				}
				postfixExpress = append(postfixExpress, data)
			}
			i++
		//	(3) 该运算符为非括号运算符:
		case express[i] == '+' || express[i] == '-' || express[i] == '*' || express[i] == '/':
			//	(a)若运算符堆栈为空,则直接压入运算符堆栈。
			if opStack.IsEmpty() {
				opStack.Push(fmt.Sprintf("%c", express[i]))
				i++
				continue LABEL
			}
			data, _ := opStack.Top()
			//	(b)若运算符堆栈栈顶的运算符为括号，则直接压入运算符堆栈。(只可能为左括号这种情况)
			if data[0] == '(' {
				opStack.Push(fmt.Sprintf("%c", express[i]))
				i++
				continue LABEL
			}
			//	(c)若比运算符堆栈栈顶的运算符优先级低或相等，则输出栈顶运算符到后缀表达式,直到栈为空或者找到优先级高于当前运算符。并将当前运算符压入运算符堆栈。
			if priority(express[i]) <= priority(data[0]) {
				tmp := priority(express[i])
				for !opStack.IsEmpty() && tmp <= priority(data[0]) {
					postfixExpress = append(postfixExpress, data)
					opStack.Pop()
					data, _ = opStack.Top()
				}
				opStack.Push(fmt.Sprintf("%c", express[i]))
				i++
				continue LABEL
			}
			//	(d)若比运算符堆栈栈顶的运算符优先级高，则直接压入运算符堆栈。
			opStack.Push(fmt.Sprintf("%c", express[i]))
			i++
		default:
			err = fmt.Errorf("invalid express:%v", express[i])
			return
		}
	}
	//	3. 扫描结束，将运算符堆栈中的运算符依次弹出，存入后缀表达式。
	for !opStack.IsEmpty() {
		data, _ := opStack.Pop()
		if data[0] == '#' {
			break
		}
		postfixExpress = append(postfixExpress, data)
	}
	fmt.Println("postfixExpress:", postfixExpress)
	return
}

//	后缀表达式求值
func calc(postfixExpress []string) (result float64, err error) {
	var (
		num1 string
		num2 string
		s    Stack //	操作栈，用于存入操作数，运算符
	)
	//	从左至右扫描后缀表达式
	for i := 0; i < len(postfixExpress); i++ {
		var cur = postfixExpress[i]
		//	1. 若读取的是运算符
		if cur[0] == '+' || cur[0] == '-' || cur[0] == '*' || cur[0] == '/' {
			//	从操作栈中弹出两个数进行运算
			num1, err = s.Pop()
			if err != nil {
				return
			}
			num2, err = s.Pop()
			if err != nil {
				return
			}
			//	先弹出的数为B，后弹出的数为A
			B, _ := strconv.ParseFloat(num1, 64)
			A, _ := strconv.ParseFloat(num2, 64)
			var res float64
			switch cur[0] {
			case '+':
				res = A + B
			case '-':
				res = A - B
			case '*':
				res = A * B
			case '/':
				res = A / B
			default:
				err = fmt.Errorf("invalid operation")
				return
			}
			//	将中间结果压栈
			s.Push(fmt.Sprintf("%f", res))
		} else {
			//	1. 若读取的是操作数，直接压栈
			s.Push(cur)
		}
	}
	//	计算结束，栈顶保存最后结果
	resultStr, err := s.Top()
	if err != nil {
		return
	}
	result, err = strconv.ParseFloat(resultStr, 64)
	return
}

/*-----------------网络请求----------------*/
// PostJSON 发送post 请求 timeout 单位 秒
func PostJSON(ctx context.Context, url string, data interface{}, header map[string]string, timeout time.Duration) (res []byte, err error) {
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
func PostFormUrlencoded(ctx context.Context, url string, data string, header map[string]string, timeout time.Duration) (res []byte, err error) {
	request, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return res, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range header {
		request.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	return respData, nil
}

// HTTPGet 发送get 请求 timeout 单位 秒
func HTTPGet(ctx context.Context, url string, header map[string]string, timeout time.Duration) (res []byte, err error) {
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

// URLEncode urlencode()
func URLEncode(str string) string {
	return url.QueryEscape(str)
}

// URLDecode urldecode()
func URLDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// HTTPBuildQuery http_build_query()
func HTTPBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

//GetIP 获取IP地址
func GetIP(req *http.Request) (IP string) {
	XForwardedFor := "X-Forwarded-For"
	XRealIP := "X-Real-IP"
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = netDefault.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

/*-----------------------SMS---------------------*/
type AliSms struct {
	Key          string `mapstructure:"key" json:"key"`
	Secret       string `mapstructure:"secret" json:"secret"`
	SignName     string `mapstructure:"sign_name" json:"sign_name"`
	TemplateCode string `mapstructure:"template_code" json:"template_code"`
}

//SendAliMessage 发送阿里云短信消息
//phoneNumber 接收人手机号码
//templateType 短信模板类型 1-找回(修改)密码
//templateParam 模板参数
func SendAliMessage(phoneNumber string, templateType int, templateParam []byte, AliConfig AliSms) (err error) {
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", AliConfig.Key, AliConfig.Secret)
	if err != nil {
		fmt.Println(err)
		//阿里云初始化失败的错误处理
		return
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phoneNumber    //手机号
	request.QueryParams["SignName"] = AliConfig.SignName //阿里云验证过的项目名 自己设置
	var templateCode string
	if templateType == 1 {
		templateCode = AliConfig.TemplateCode
	}
	request.QueryParams["TemplateCode"] = templateCode //阿里云的短信模板号 自己设置
	//request.QueryParams["TemplateParam"] = "{\"code\":" + "123456" + ",\"product\":" + "nbi" + "}" //短信模板中的验证码内容 自己生成
	request.QueryParams["TemplateParam"] = string(templateParam) //短信模板中的验证码内容 自己生成
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		fmt.Println(err)
		//短息发送失败的错误处理
		return
	}
	//json数据解析
	var message struct {
		Message   string
		RequestId string
		BizId     string
		Code      string
	} //阿里云返回的json信息对应的类
	_ = json.Unmarshal(response.GetHttpContentBytes(), &message)
	if message.Message != "OK" {
		fmt.Println(message)
		//阿里云操作失败的错误处理
		return errors.New("发送短信失败:" + message.Message)
	}
	return nil
}

/*--------------------字符串处理------------------*/
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//RandomString 随机字符串
func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[mathRand.Intn(len(letterBytes))]
	}
	return string(b)
}

// StrRev 字符串反转
func StrRev(str string) string {
	return php2go.Strrev(str)
}

/*-------------------性能---------------------*/
//VirtualMemory 获取虚拟内存
func VirtualMemory() (info *mem.VirtualMemoryStat) {
	info, _ = mem.VirtualMemory()
	return
}

//IOCounters 获取网络io
func IOCounters() (IoInfo []net.IOCountersStat) {
	IoInfo, _ = net.IOCounters(false)
	return
}

//CPUPercent 获取cpu占用比例
func CPUPercent() (cpuInfo []float64) {
	cpuInfo, _ = cpu.Percent(time.Second, false)
	return
}

/*-----------------时间处理---------------------*/
// MakeTimeType 时间返回格式
type MakeTimeType struct {
	Year    string `json:"year"`
	Quarter string `json:"quarter"`
	Month   string `json:"month"`
	Day     string `json:"day"`
}

//FormDateToUnit 日期格式转时间戳
func FormDateToUnit(formDate string) (int64, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tInfo, _ := time.ParseInLocation("20060102", formDate, loc)
	return tInfo.Unix(), nil
}

//FormDateToUnitAndLayout 日期格式转时间戳
func FormDateToUnitAndLayout(formDate string, layout string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tInfo, _ := time.ParseInLocation(layout, formDate, loc)
	return tInfo.Unix()
}

// GetMonth 根据时间戳获取月份【1-12】
func GetMonth(timeStamp int64) (month int) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	date := time.Unix(timeStamp, 0).In(loc).Format("2006-01-02")
	month, _ = strconv.Atoi(strings.Split(date, "-")[1])
	return
}

//MapDate 根据两个时间戳返回两个时间的日期切片 eg: 20200316
func MapDate(startTime, endTime int64) (list []string, err error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	//获取开始时间的凌晨零点
	t1, err := time.ParseInLocation("20060102150405", time.Unix(startTime, 0).In(loc).Format("20060102")+"000000", loc)
	if err != nil {
		return
	}
	//时区
	//var zone = time.FixedZone("CST",8*3600)//东八
	ft := t1.Unix()
	//日期
	for i := ft; i <= endTime; i = i + 86400 { //<= 如果结束时间相等（结束时间是第二天的凌晨）那一天还算
		temp := time.Unix(i, 0).In(loc).Format("20060102")
		list = append(list, temp)
	}
	return
}

//GetLatestSeasonPeriod 获取最近几个季度区间【start,end】
func GetLatestSeasonPeriod(count int) (List []map[string]time.Time) {
	ThisMonth := int(time.Now().Month())
	cutMonth := ThisMonth % 3     //距离最近的季度最后一个月相差几个月
	for i := 1; i <= count; i++ { //一共要读取几个季度片段
		season := make(map[string]time.Time)
		SeasonEndMonth := time.Now().AddDate(0, -cutMonth-(i-1)*3, 0) //每个季度片段的结束时间
		seasonEndMonthInt, _ := strconv.ParseInt(SeasonEndMonth.Format("01"), 10, 64)
		SeasonEnd := GetMonthStartAndEnd(SeasonEndMonth.Year(), int(seasonEndMonthInt))
		season["end"] = SeasonEnd["end"]
		SeasonStartMonth := time.Now().AddDate(0, -cutMonth-(i-1)*3-2, 0) //每个季度片段的开始时间
		seasonStartMonthInt, _ := strconv.ParseInt(SeasonStartMonth.Format("01"), 10, 64)
		SeasonStart := GetMonthStartAndEnd(SeasonStartMonth.Year(), int(seasonStartMonthInt))
		season["start"] = SeasonStart["start"]
		List = append(List, season)
	}
	return
}

//GetMonthStartAndEnd 获取月份的第一天和最后一天
func GetMonthStartAndEnd(myYear int, myMonth int) map[string]time.Time {
	// 数字月份必须前置补零
	var myMonthString string
	if myMonth < 10 {
		myMonthString = "0" + strconv.Itoa(myMonth)
	} else {
		myMonthString = strconv.Itoa(myMonth)
	}
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, strconv.Itoa(myYear)+"-"+myMonthString+"-01 00:00:00", loc)
	newMonth := theTime.Month()
	t1 := time.Date(myYear, newMonth, 1, 0, 0, 0, 0, time.Local)
	t2 := time.Date(myYear, newMonth+1, 0, 23, 59, 59, 0, time.Local)
	result := map[string]time.Time{"start": t1, "end": t2}
	return result
}

// MakeYearTimeSlice 按年份拆分时间切片
// start 2006-01-02 15:04:05
//end 2006-01-02 15:04:05
func MakeYearTimeSlice(start string, end string) (arr []string) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	for {
		if startTime.After(endTime) {
			//跳出
			break
		}
		arr = append(arr, startTime.Format("2006"))
		start = startTime.AddDate(1, 0, 0).Format("2006") //+1年
		startTime, _ = time.ParseInLocation("2006", start, time.Local)
	}
	return arr
}

//MakeQuarterTimeSlice 按季度拆分时间切片
//start 2006-01-02 15:04:05
func MakeQuarterTimeSlice(start string, end string) (arr []MakeTimeType) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	for {
		if startTime.After(endTime) {
			//跳出
			break
		}
		item := MakeTimeType{}
		item.Year = startTime.Format("2006")
		item.Month = startTime.Format("01")
		item.Day = startTime.Format("02")
		switch item.Month {
		case "01", "02", "03":
			item.Quarter = "1"
		case "04", "05", "06":
			item.Quarter = "2"
		case "07", "08", "09":
			item.Quarter = "3"
		case "10", "11", "12":
			item.Quarter = "4"
		}
		arr = append(arr, item)
		start = startTime.AddDate(0, 3, 0).Format("2006-01-02 15:04:05")
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	}
	return arr
}

//MakeMonthTimeSlice 按月份获取时间切片
//start 2006-01-02 15:04:05
func MakeMonthTimeSlice(start string, end string) (arr []MakeTimeType) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	for {
		if startTime.After(endTime) {
			//跳出
			break
		}
		item := MakeTimeType{}
		item.Year = startTime.Format("2006")
		item.Month = startTime.Format("01")
		item.Day = startTime.Format("02")
		switch item.Month {
		case "01", "02", "03":
			item.Quarter = "1"
		case "04", "05", "06":
			item.Quarter = "2"
		case "07", "08", "09":
			item.Quarter = "3"
		case "10", "11", "12":
			item.Quarter = "4"
		}
		arr = append(arr, item)
		start = startTime.AddDate(0, 1, 0).Format("2006-01-02 15:04:05")
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	}
	return arr
}

//MakeDayTimeSlice 按日期获取时间切片
//start 2006-01-02 15:04:05
func MakeDayTimeSlice(start string, end string) (arr []MakeTimeType) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	for {
		if startTime.After(endTime) {
			break
		}
		item := MakeTimeType{}
		item.Year = startTime.Format("2006")
		item.Month = startTime.Format("01")
		item.Day = startTime.Format("02")
		switch item.Month {
		case "01", "02", "03":
			item.Quarter = "1"
		case "04", "05", "06":
			item.Quarter = "2"
		case "07", "08", "09":
			item.Quarter = "3"
		case "10", "11", "12":
			item.Quarter = "4"
		}
		arr = append(arr, item)
		start = startTime.AddDate(0, 0, 1).Format("2006-01-02 15:04:05")
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	}
	return arr
}

//SplitDayTimeSlice 时间切片
//start 2006-01-02 15:04:05
func SplitDayTimeSlice(start time.Time, end time.Time) (arr []time.Time) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", start.Format("2006-01-02 15:04:05"), time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end.Format("2006-01-02 15:04:05"), time.Local)
	for {
		if startTime.After(endTime) {
			break
		}
		arr = append(arr, startTime)
		start = startTime.AddDate(0, 0, 1)
		startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", start.Format("2006-01-02 15:04:05"), time.Local)
	}
	return arr
}

/*------------------------堆栈------------------*/
type Stack struct {
	data [1024]string
	top  int
}

func (s *Stack) IsEmpty() bool {
	return s.top == 0
}
func (s *Stack) Top() (ret string, err error) {
	if s.top == 0 {
		err = fmt.Errorf("stack is empty")
		return
	}
	ret = s.data[s.top-1]
	return
}
func (s *Stack) Push(str string) {
	s.data[s.top] = str
	s.top++
}
func (s *Stack) Pop() (ret string, err error) {
	if s.top == 0 {
		err = fmt.Errorf("stack is empty")
		return
	}
	s.top--
	ret = s.data[s.top]
	return
}

/*---------------------会话控制----------------------*/
