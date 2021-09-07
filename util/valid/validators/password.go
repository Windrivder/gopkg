package validators

import (
	"github.com/windrivder/gopkg/container/typex"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/util/valid"
)

// 用户密码的正则匹配, 合法的字符有 0-9, A-Z, a-z, 特殊字符
type password struct{}

func Password() valid.Validator {
	return &password{}
}

func (p *password) Name() string {
	return "password"
}

func (p *password) Trans() typex.DictStrs {
	return typex.DictStrs{
		"en": "{0} length 4 to 32, contain digits, letters, special characters",
		"zh": "{0} 长度在 4 到 32 位，包含数字、字母、特殊字符",
	}
}

func (p *password) Validate(i interface{}) error {
	v, ok := i.(string)
	if ok {
		if len(v) < 4 || len(v) > 32 {
			return errorx.New("length 4 to 32")
		}

		if RegPassword.MatchString(v) {
			return nil
		}
	}

	return errorx.New("only digits, letters, special characters")
}

func PasswordValidate(password string) error {
	return Password().Validate(password)
}
