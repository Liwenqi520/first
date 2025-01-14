package logic

import (
	"context"
	"encoding/base64"
	"first/init/myredis"
	"first/init/mysql"
	"first/internal/common"
	"first/internal/defines"
	"first/internal/model"
	"first/internal/schema"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	// "first/gopkg/log"

	"github.com/Liwenqi520/errorx"
	"github.com/Liwenqi520/helper"
	"github.com/Liwenqi520/logger"
)

type User struct{}

const (
	BREWING = 1 // 酿造
	WINEJAR = 2 // 酒坛
	ClOUD   = 3 //云平台
)

// Login 登录
func (u User) Login(ctx context.Context, Params schema.LoginParams, Host, IP string) (res interface{}, err error) {
	spanCtx := logger.Start(ctx, "user login")
	defer logger.End(spanCtx)
	Type := Params.Type
	var LoginUser model.AuthUser
	if Type == 0 {
		Type = 1
	}
	switch Type {
	case 1:
		LoginUser, err = u.CheckUserPassword(ctx, Params)
		break
	case 2:
		LoginUser, err = u.CheckPhoneCode(ctx, Params)
		break
	case 3:
		LoginUser, err = u.CheckEmailUser(ctx, Params)
		break
	}
	if err != nil {
		return
	}
	if LoginUser.UserType != 1 {
		err = errorx.New(330056, "班长不允许登录平台")
		return
	}
	IsOnline := Common{}.GetIsOnline()
	Url := "http://" + Host
	if IsOnline == 1 {
		Url += "/api"
	}
	// 缓存用户登录信息
	storeInfo := schema.StoreInfo{
		ID:     LoginUser.ID,
		RoleID: LoginUser.RoleID,
		// RoleType:                 RoleType,
		CompanyID:      LoginUser.CompanyID,
		UserName:       LoginUser.UserName,
		NickName:       LoginUser.NickName,
		Image:          LoginUser.Image,
		Mobile:         LoginUser.Mobile,
		Email:          LoginUser.Email,
		Status:         LoginUser.Status,
		IsAdmin:        LoginUser.IsAdmin,
		PlatformType:   LoginUser.PlatformType,
		AreaPermission: LoginUser.AreaPermission,
		HavePlatforms:  LoginUser.HavePlatforms,
	}
	token, err := u.StoreInRedis(ctx, storeInfo)
	if err != nil {
		logger.Error(ctx, "login store error", logger.String("error", err.Error()))
		err = errorx.New(200001, "token缓存失败")
		return
	}
	var data schema.LoginRes
	// data.Token = token
	data.AccessToken = token
	data.RefreshToken = token
	data.UserName = LoginUser.UserName
	now := time.Now()
	new := now.Add(30 * 24 * time.Hour)
	data.Expires = new.Format("2006/01/02 15:04:05")
	data.Avatar = LoginUser.Image
	data.Roles = strings.Split(LoginUser.RoleID, ",")

	data.UserInfo.ID = LoginUser.ID
	data.UserInfo.RoleID = LoginUser.RoleID
	data.UserInfo.CompanyID = LoginUser.CompanyID
	data.UserInfo.UserName = LoginUser.UserName
	data.UserInfo.NickName = LoginUser.NickName
	data.UserInfo.Image = LoginUser.Image
	data.UserInfo.Mobile = LoginUser.Mobile
	data.UserInfo.Email = LoginUser.Email
	data.UserInfo.Status = LoginUser.Status
	data.UserInfo.IsAdmin = LoginUser.IsAdmin
	data.UserInfo.PlatformType = LoginUser.PlatformType
	data.UserInfo.AreaPermission = LoginUser.AreaPermission
	data.UserInfo.HavePlatforms = LoginUser.HavePlatforms
	res = data
	return
}

func (u User) CheckPhoneCode(ctx context.Context, Params schema.LoginParams) (LoginUser model.AuthUser, err error) {
	db := mysql.NewDB()
	Phone := strings.TrimSpace(Params.Phone)
	Code := strings.TrimSpace(Params.VerifyCode)
	if Phone == "" {
		err = errorx.New(330048, "手机号不能为空")
		return
	}
	if Code == "" {
		err = errorx.New(330049, "短信验证码不能为空")
		return
	}
	rdb := myredis.RedisCli()
	PhoneKey := defines.PhoneLogin + Phone
	StoreCode, err := rdb.Get(ctx, PhoneKey).Result()
	if StoreCode != Code {
		err = errorx.New(330050, "短信验证码错误")
		return
	}
	err = db.Model(&model.AuthUser{}).Where("mobile = ? and deleted_time = 0", Phone).First(&LoginUser).Error
	if err != nil {
		err = errorx.New(330017, "账号不存在")
		return
	}
	if LoginUser.DeletedTime != 0 {
		err = errorx.New(330017, "账号不存在")
		return
	}
	if LoginUser.Status != 1 {
		err = errorx.New(330016, "该账号已被禁用，请联系管理员")
		return
	}
	return
}

func (u User) CheckEmailUser(ctx context.Context, Params schema.LoginParams) (LoginUser model.AuthUser, err error) {
	db := mysql.NewDB()
	Email := strings.TrimSpace(Params.Email)
	Code := strings.TrimSpace(Params.VerifyCode)
	if Email == "" {
		err = errorx.New(330045, "邮箱不能为空")
		return
	}
	if Code == "" {
		err = errorx.New(330046, "邮箱验证码不能为空")
		return
	}
	rdb := myredis.RedisCli()
	EmailKey := defines.EmailLogin + Email
	StoreCode, err := rdb.Get(ctx, EmailKey).Result()
	if StoreCode != Code {
		err = errorx.New(330047, "邮箱验证码错误")
		return
	}
	err = db.Model(&model.AuthUser{}).Where("email = ? and deleted_time = 0", Email).First(&LoginUser).Error
	if err != nil {
		err = errorx.New(330017, "账号不存在")
		return
	}
	if LoginUser.DeletedTime != 0 {
		err = errorx.New(330017, "账号不存在")
		return
	}
	if LoginUser.Status != 1 {
		err = errorx.New(330016, "该账号已被禁用，请联系管理员")
		return
	}
	return
}

// CheckUserPassword 验证用户账号密码
func (u User) CheckUserPassword(ctx context.Context, Params schema.LoginParams) (LoginUser model.AuthUser, err error) {
	db := mysql.NewDB()
	Username := strings.TrimSpace(Params.Account)
	Password := strings.TrimSpace(Params.Password)
	if Username == "" {
		err = errorx.New(330035, "登录账号不能为空")
		return
	}
	if Password == "" {
		err = errorx.New(330036, "密码不能为空")
		return
	}
	err = db.Model(&model.AuthUser{}).Where("user_name = ?", Username).Order("created_time desc").First(&LoginUser).Error
	if err != nil {
		err = errorx.New(330017, "账号不存在")
		return
	}
	if LoginUser.DeletedTime != 0 {
		err = errorx.New(330017, "账号不存在")
		return
	}
	if LoginUser.Status != 1 {
		err = errorx.New(330016, "该账号已被禁用，请联系管理员")
		return
	}
	plainPwd, err := u.PwdCodeCheck(ctx, Password)
	if err != nil {
		return
	}
	ok := helper.ComparePasswords(LoginUser.Password, plainPwd)
	if !ok {
		logger.Error(ctx, "pwd", logger.String("解析后的pwd:", string(plainPwd)), logger.String("db.password:", LoginUser.Password))
		err = errorx.New(330015, "您输入的账号或密码错误，请确认后重新登录")
		return
	}
	return
}

func (u User) CheckUsernameExist(Name, ID string) bool {
	db := mysql.NewDB()
	var UserNameExist model.AuthUser
	QueryUser := db.Model(&model.AuthUser{}).Where("user_name = ? and deleted_time = 0", Name)
	if ID != "" {
		QueryUser = QueryUser.Where("id != ?", ID)
	}
	QueryUser.Take(&UserNameExist)
	return UserNameExist.ID != ""
}

// PwdCodeCheck 密码解密
func (u User) PwdCodeCheck(ctx context.Context, Password string) (pwd []byte, err error) {
	// 验证
	// 密码解码
	cipherText, err := base64.StdEncoding.DecodeString(Password)
	if err != nil {
		logger.Error(ctx, "base64", logger.String("error", "密码base64解码失败"))
		return nil, errorx.NewWithMsg(200001, "密码base64解码失败")
	}
	// 获取私钥
	rsaData, err := u.GetRSAPrivateKey()
	if err != nil {
		logger.Error(ctx, "rsa-login", logger.String("error", "获取私钥失败"))
		return nil, errorx.New(200001)
	}
	getPwd, err := helper.RSA{}.RSADecrypt(cipherText, rsaData)
	if err != nil {
		logger.Error(ctx, "rsa-login", logger.String("error", "rsa解密失败"))
		return nil, errorx.New(200001)
	}

	return getPwd, nil

}

// GetRsaPrivateKey 获取RSA私钥
func (u User) GetRSAPrivateKey() (privateKey []byte, err error) {
	file, err := ioutil.ReadFile("./etc/rsa/private.pem")
	if err != nil {
		logger.Info(context.Background(), "get rsa private key error")
		err = errorx.New(0, "获取秘钥错误")
	}
	return file, err
}

// StoreInRedis 信息写入redis
func (u User) StoreInRedis(ctx context.Context, userInfo schema.StoreInfo) (token string, err error) {
	// 生成token
	token = helper.Md5(userInfo.UserName + time.Now().String() + userInfo.ID)
	tokenKey := defines.SESSION_TOKEN + token
	userInfoKey := defines.SESSION_USER_INFO + userInfo.ID
	userTokensKey := defines.SESSION_USER_TOKENS + userInfo.ID

	rdb := myredis.RedisCli()
	// token字符串存储
	if _, err := rdb.Set(context.Background(), tokenKey, userInfo.ID, 0).Result(); err != nil {
		return "", err
	}
	// 保存个人信息
	if _, err := rdb.HMSet(context.Background(), userInfoKey,
		"id", userInfo.ID,
		"role_id", userInfo.RoleID,
		"role_type", userInfo.RoleType,
		"company_id", userInfo.CompanyID,
		"nick_name", userInfo.NickName,
		"user_name", userInfo.UserName,
		"image", userInfo.Image,
		"mobile", userInfo.Mobile,
		"email", userInfo.Email,
		"status", userInfo.Status,
		"is_admin", strconv.Itoa(int(userInfo.IsAdmin)),
		"platform_type", strconv.Itoa(int(userInfo.PlatformType)),
		"area_permission", userInfo.AreaPermission,
		"have_platforms", userInfo.HavePlatforms,
	).Result(); err != nil {
		return "", err
	}

	// 保存有效tokens
	if _, err := rdb.SAdd(context.Background(), userTokensKey, tokenKey).Result(); err != nil {
		return "", err
	}
	// 清理无效tokens
	strs, err := rdb.SMembers(context.Background(), userTokensKey).Result()
	if err != nil {
		return "", err
	}
	for _, v := range strs {
		count, _ := rdb.Exists(context.Background(), v).Result()
		if count != 1 {
			// 从集合里删除
			if _, err := rdb.SRem(context.Background(), userTokensKey, v).Result(); err != nil {
				return "", err
			}
		}
	}
	return
}

// DestroyToken销毁用户token
func (u User) DestroyToken(ctx context.Context, token string, userInfo common.UserInfo, IP string) (err error) {
	rdb := myredis.RedisCli()
	tokenKey := defines.SESSION_TOKEN + token
	userTokensKey := defines.SESSION_USER_TOKENS + userInfo.ID

	// 删除tokenKey
	_, err = rdb.Del(context.Background(), tokenKey).Result()
	if err != nil {
		logger.Error(ctx, "DestroyToken key", logger.String("error", err.Error()))
		return nil
	}
	// 从集合里删除
	if _, err := rdb.SRem(context.Background(), userTokensKey, tokenKey).Result(); err != nil {
		logger.Error(ctx, "DestroyToken", logger.String("error", err.Error()))
		return nil
	}
	return nil
}
