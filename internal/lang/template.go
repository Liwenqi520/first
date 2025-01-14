package lang

//LangTemplate 语言模板
var LangTemplate = map[string]map[interface{}]interface{}{
	"zh-cn": {
		"has {times} to try": "账号或密码错误，剩余{times}次机会",
	},
	"ja-jp": {
		"has {times} to try": "アカウントまたはパスワードが正しくない、{times}の機会が残っている",
	},
	"en-us": {
		"has {times} to try": "Incorrect account or password, {times} opportunities remaining",
	},
}
