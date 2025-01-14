package i18n

import (
	"reflect"
	"strings"
)

// LangPack 语言集合
var LangPack = make(map[string]map[interface{}]interface{})

// 默认语言
var defaultLang = "zh-cn"

// LoadLangPack 加载语言包
func LoadLangPack(langPacks ...map[string]map[interface{}]interface{}) {
	for _, langPack := range langPacks {
		for l, lv := range langPack {
			if LangPack[l] == nil {
				LangPack[l] = make(map[interface{}]interface{})
			}
			for k, v := range lv {
				LangPack[l][k] = v
			}
		}
	}
}

// T 翻译key对应的语言
// lang 语种
// key 翻译字段
// tpl 模板内容
func T(lang string, key interface{}, tpl ...map[string]string) (transRes string) {
	if lang == "" || len(strings.Split(lang, ";")) > 1 {
		lang = defaultLang
	}
	var langPack map[interface{}]interface{}
	if langPack = LangPack[lang]; langPack == nil {
		return ""
	}
	// 非字符串类型的key，直接获取结果
	if _, ok := key.(string); !ok {
		if s, ok := langPack[key].(string); ok {
			transRes = s
		} else {
			return ""
		}
		// 字段类型为string的处理
		// 支持多级嵌套
	} else if keyString, ok := key.(string); ok {
		keyChildren := strings.Split(keyString, ".")
		if len(keyChildren) == 1 {
			if s, ok := langPack[keyString].(string); ok {
				transRes = s
			} else {
				return ""
			}
		}
		// 嵌套的处理
		if len(keyChildren) > 1 {
			var newMap interface{}
			newMap = langPack
			for i := 0; i < len(keyChildren); i++ {
				if lt, ok := newMap.(map[interface{}]interface{}); ok {
					newMap = lt[keyChildren[i]]
					continue
				} else {
					newMap = reflect.ValueOf(newMap).MapIndex(reflect.ValueOf(keyChildren[i])).Interface()
				}
			}
			transRes, _ = newMap.(string)
		}
		// 其他不支持的类型字段，返回 ""
	} else {
		return ""
	}
	// 模板处理
	if len(tpl) > 0 && strings.ContainsAny(transRes, "{}") {
		tplContent := tpl[0]
		for k, v := range tplContent {
			transRes = strings.ReplaceAll(transRes, "{"+k+"}", v)
		}
	}
	return transRes
}
