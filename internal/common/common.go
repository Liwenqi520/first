package common

import (
	"context"
	"first/init/myredis"
	"first/internal/defines"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Liwenqi520/errorx"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserInfo struct {
	ID             string `json:"id"`
	CompanyID      string `json:"company_id"`
	RoleID         string `json:"role_id"`
	RoleType       int8   `json:"role_type"`
	UserName       string `json:"user_name"`
	NickName       string `json:"nick_name"`
	Image          string `json:"image"`
	Mobile         string `json:"mobile"`
	Email          string `json:"email"`
	Status         int8   `json:"status"`
	IsAdmin        int8   `json:"is_admin"`
	AreaPermission string `json:"area_permission"` // 区域权限ID集合
	// FactorgGroupManageButton int64  `json:"factory_group_manage_button"` // 产区管理按钮
	// FactorgGroupSwitch       int64  `json:"factory_group_switch"`        // 产区开关
	// WaterQualitySwitch       int64  `json:"water_quality_switch"`        // 水质开关
	Authorization string `json:"authorization"`
	ShareKey      string `json:"share_key"`
	HavePlatforms string `json:"have_platforms"`
}

func SessionGet(c *gin.Context) (userInfo UserInfo, err error) {
	Authorization := c.GetHeader("Authorization")
	key := defines.SESSION_TOKEN + Authorization
	rds := myredis.RedisCli()
	userID, err := rds.Get(context.Background(), key).Result()
	if err != nil {
		return userInfo, err
	}
	userInfoKey := defines.SESSION_USER_INFO + userID
	userInfoMap, err := rds.HGetAll(context.Background(), userInfoKey).Result()
	if err != nil {
		return userInfo, err
	}
	userInfo.ID = userID
	userInfo.CompanyID = userInfoMap["company_id"]
	userInfo.RoleID = userInfoMap["role_id"]
	RoleType, _ := strconv.Atoi(userInfoMap["role_type"])
	userInfo.RoleType = int8(RoleType)
	userInfo.UserName = userInfoMap["user_name"]
	userInfo.NickName = userInfoMap["nick_name"]
	userInfo.Image = userInfoMap["image"]
	userInfo.Mobile = userInfoMap["mobile"]
	userInfo.Email = userInfoMap["email"]
	isAdmin, _ := strconv.Atoi(userInfoMap["is_admin"])
	userInfo.IsAdmin = int8(isAdmin)
	status, _ := strconv.Atoi(userInfoMap["status"])
	userInfo.Status = int8(status)
	userInfo.Authorization = Authorization
	userInfo.AreaPermission = userInfoMap["area_permission"]
	userInfo.HavePlatforms = userInfoMap["have_platforms"]
	return userInfo, err
}

func GetErrCode(Params interface{}, err error) error {
	var ErrorCode int64
	t := reflect.TypeOf(Params).Elem()
	binding := make(map[string]string)
	errcode := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		binding[t.Field(i).Name] = t.Field(i).Tag.Get("binding")
		errcode[t.Field(i).Name] = t.Field(i).Tag.Get("errcode")
	}
	for _, errVal := range err.(validator.ValidationErrors) {
		fmt.Println(errVal.ActualTag())
		FieldName := strings.Split(errVal.Namespace(), ".")[1]
		if bindingString, ok := binding[FieldName]; ok {
			bindingList := strings.Split(bindingString, ",")
			for k, v := range bindingList {
				if strings.Split(v, "=")[0] == errVal.ActualTag() {
					errcodeList := strings.Split(errcode[FieldName], ",")
					if k >= len(errcodeList) {
						ErrorCode = 100000
					} else {
						ErrorCode, _ = strconv.ParseInt(errcodeList[k], 10, 64)
					}
				}
			}
		}
	}

	return errorx.New(ErrorCode, "")
}

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

// CheckEmail 校验邮箱
func CheckEmail(email string) bool {
	// 邮箱正则表达式
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// 编译正则表达式
	reg := regexp.MustCompile(pattern)
	// 执行匹配
	return reg.MatchString(email)
}
