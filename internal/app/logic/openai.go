package logic

import (
	"first/init/mysql"
	"first/internal/common"
	"first/internal/model"
	"first/internal/schema"
	"log"
	"math"
	"time"

	"github.com/Liwenqi520/errorx"
	"github.com/Liwenqi520/helper"

	openai "github.com/all-in-aigc/openai-go"
)

type OpenAi struct{}

var APIKEY = "sk-nZLjmKc5QSJdVobt27607d42A1Dd4c7a9aAd74D11b761e8f"
var BaseUri = "https://api.mixrai.com/v1"
var Model = "gpt-4o"

// 使用微信支付完成下单接口
var WxPayUrl = "https://api.mixrai.com/v1/wxpay/unifiedorder"

// 使用微信支付完成支付接口
var WxPayNotifyUrl = "https://api.mixrai.com/v1/wxpay/notify"

// 使用微信支付完成退款接口
var WxPayRefundUrl = "https://api.mixrai.com/v1/wxpay/refund"

// 使用微信支付完成查询订单接口
var WxPayQueryUrl = "https://api.mixrai.com/v1/wxpay/query"

// 使用微信支付完成关闭订单接口
var WxPayCloseUrl = "https://api.mixrai.com/v1/wxpay/close"

// 使用微信支付完成下载账单接口
var WxPayDownloadBillUrl = "https://api.mixrai.com/v1/wxpay/downloadbill"

// 使用微信支付完成查询退款接口
var WxPayQueryRefundUrl = "https://api.mixrai.com/v1/wxpay/queryrefund"

func (ai OpenAi) GetClinet() (client *openai.Client) {
	apiKey := APIKEY
	client, _ = openai.NewClient(&openai.Options{
		ApiKey:  apiKey,
		Timeout: 30 * time.Second,
		Debug:   true,
		BaseUri: BaseUri,
	})
	return
}

func (ai OpenAi) ChatCompletions(Params schema.OpenAiParams, Admin common.UserInfo) (ResponseMsg interface{}, err error) {
	StartUnix := time.Now().Unix()
	cli := ai.GetClinet()
	uri := "/chat/completions"
	params := map[string]interface{}{
		"model": Model,
		"messages": []map[string]interface{}{
			{"role": "user", "content": Params.Message},
		},
		// "functions": getFuncs(),
	}
	// query user money
	AvailableAmount, err := ai.CheckUserMoney(Admin.ID)
	if err != nil {
		err = errorx.New(660001, "用户余额查询失败，请稍后再试")
		return
	}
	if AvailableAmount <= 0 {
		err = errorx.New(880002, "用户余额不足，请先充值")
		return
	}
	res, err := cli.Post(uri, params)
	if err != nil {
		log.Fatalf("request api failed: %v", err)
		return
	}
	EndUnix := time.Now().Unix()
	UsedToken := res.Get("usage.total_tokens").Int()
	// 根据用户使用的Token数量，扣除用户余额
	ai.DeductUserMoney(Admin.ID, UsedToken)
	// keep in log
	ai.SaveOpenaiLog(Admin.ID, Params.Message, res.Get("choices.0.message.content").String(), EndUnix-StartUnix, UsedToken)
	//
	ResponseMsg = res.Get("choices.0.message.content").String()
	return
}

// UserRecharge 用户充值
func (ai OpenAi) UserRecharge(ReqParams schema.OpenAiRecharge, Admin common.UserInfo) (ResponseMsg interface{}, err error) {
	// 开启事务的db
	db := mysql.NewDB().Begin()
	now := time.Now().Unix()
	// 第一步：查询用户余额
	var UserCash model.UserCash
	err = db.Model(&model.UserCash{}).Where("user_id = ? and deleted_time = 0", Admin.ID).First(&UserCash).Error
	if err != nil {
		db.Rollback()
		return
	}
	// 第二步：充值金额到用户余额中
	AfterTotalAmount := UserCash.TotalAmount + ReqParams.Amount/100         // 充值后用户总金额【元】
	AfterAvailableAmount := UserCash.AvailableAmount + ReqParams.Amount/100 // 充值后用户可用金额【元】
	UpdateParams := make(map[string]interface{})
	UpdateParams["total_amount"] = AfterTotalAmount
	UpdateParams["available_amount"] = AfterAvailableAmount
	UpdateParams["updated_time"] = now
	err = db.Model(&model.UserCash{}).Where("user_id = ? and deleted_time = 0", Admin.ID).UpdateColumns(&UpdateParams).Error
	if err != nil {
		db.Rollback()
		return
	}
	// 第三步：记录用户充值流水
	var UserTranscation model.Transcation
	UserTranscation.ID, _ = helper.UUID{}.GetUniqueKey()
	UserTranscation.UserID = Admin.ID
	UserTranscation.TotalAmount = AfterTotalAmount
	UserTranscation.ChangeAmount = ReqParams.Amount / 100
	UserTranscation.CreatedTime = now
	err = db.Create(&UserTranscation).Error
	if err != nil {
		db.Rollback()
		return
	}
	db.Commit()
	return
}

// DeductUserMoney 扣除用户余额[一次请求大概使用1000Token，1000Token=0.1元]
func (ai OpenAi) DeductUserMoney(UserID string, UsedToken int64) (err error) {
	// 开启事务的db
	db := mysql.NewDB().Begin()
	now := time.Now().Unix()
	// 第一步：查询用户余额
	var UserCash model.UserCash
	err = db.Model(&model.UserCash{}).Where("user_id = ? and deleted_time = 0", UserID).First(&UserCash).Error
	// 第二步：扣除用户余额
	// 扣除的金额为UsedToken除以一千，向上取整
	DeductAmount := math.Ceil(float64(UsedToken) / 1000)
	AfterTotalAmount := UserCash.TotalAmount - int64(DeductAmount)         // 扣后用户剩余总金额
	AfterAvailableAmount := UserCash.AvailableAmount - int64(DeductAmount) // 扣后用户剩余可用金额
	AfterUsedAmount := UserCash.UsedAmount + int64(DeductAmount)           // 扣后用户已使用金额
	UpdateParams := make(map[string]interface{})
	UpdateParams["total_amount"] = AfterTotalAmount
	UpdateParams["available_amount"] = AfterAvailableAmount
	UpdateParams["used_amount"] = AfterUsedAmount
	UpdateParams["last_used_time"] = now
	UpdateParams["updated_time"] = now
	err = db.Model(&model.UserCash{}).Where("user_id = ? and deleted_time = 0", UserID).UpdateColumns(&UpdateParams).Error
	if err != nil {
		db.Rollback()
		return
	}
	// 第三步：记录用户使用流水
	var UserTranscation model.Transcation
	UserTranscation.ID, _ = helper.UUID{}.GetUniqueKey()
	UserTranscation.UserID = UserID
	UserTranscation.TotalAmount = AfterTotalAmount
	UserTranscation.ChangeAmount = -int64(DeductAmount)
	UserTranscation.CreatedTime = now
	err = db.Create(&UserTranscation).Error
	if err != nil {
		db.Rollback()
		return
	}
	db.Commit()
	return
}

// CheckUserMoney 检查用户余额
func (ai OpenAi) CheckUserMoney(UserID string) (AvailableAmount int64, err error) {
	db := mysql.NewDB()
	var UserCash model.UserCash
	db.Model(&model.UserCash{}).Where("user_id = ? and deleted_time = 0", UserID).First(&UserCash)
	if UserCash.ID == "" { // 初始化用户余额
		var AddCash model.UserCash
		AddCash.ID, _ = helper.UUID{}.GetUniqueKey()
		AddCash.UserID = UserID
		AddCash.AvailableAmount = 0
		AddCash.CreatedTime = time.Now().Unix()
		db.Create(&AddCash)
		return
	}
	db.Where("user_id = ?", UserID).First(&UserCash)
	AvailableAmount = UserCash.AvailableAmount
	return
}

func (ai OpenAi) SaveOpenaiLog(AdminID, ReqParams, Response string, WasteTime, UsedToken int64) {
	db := mysql.NewDB()
	var OpenaiLog model.OpenaiLog
	OpenaiLog.UserID = AdminID
	OpenaiLog.ReqParams = ReqParams
	OpenaiLog.Response = Response
	OpenaiLog.WasteTime = WasteTime
	OpenaiLog.UsedToken = UsedToken
	OpenaiLog.CreatedTime = time.Now().Unix()
	db.Create(&OpenaiLog)
}
