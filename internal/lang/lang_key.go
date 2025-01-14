package lang

//LangKey 语言键值对
var LangKey = map[string]map[interface{}]interface{}{
	"zh-cn": {
		"success": "成功",
		"time":    "时间",

		//设备服务平台报错
		"无权限":   "无权限",
		"设备不存在": "设备不存在",
		"参数有误":  "参数有误",
	},
	"ja-jp": {
		"success": "成功",
		"time":    "時間",

		//设备服务平台报错
		"无权限":   "権限なし",
		"设备不存在": "デバイスが存在しません",
		"参数有误":  "パラメータが間違っています",
	},
	"en-us": {
		"success": "success",
		"time":    "Time",

		//设备服务平台报错
		"无权限":   "No permission",
		"设备不存在": "Device does not exist",
		"参数有误":  "Parameter error",
	},
}
