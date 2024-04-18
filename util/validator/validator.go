package validator

import (
	"github.com/go-playground/validator/v10"
	"unicode"
)

// 自定义验证器函数，检查字段是否包含中文字符
func ChineseValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, r := range value {
		if unicode.In(r, unicode.Scripts["Han"]) {
			// 如果包含中文字符，则验证通过
			return true
		}
	}
	// 如果没有找到中文字符，则验证失败
	return false
}
