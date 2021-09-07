package validators

import (
	"github.com/windrivder/gopkg/container/typex"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/util/valid"
)

// 用户名的正则匹配, 合法的字符有 0-9, A-Z, a-z, _
// 第一个字母不能为 _, 0-9
// 最后一个字母不能为 _, 且 _ 不能连续
type username struct{}

func Username() valid.Validator {
	return &username{}
}

func (n *username) Name() string {
	return "username"
}

func (n *username) Trans() typex.DictStrs {
	return typex.DictStrs{
		"en": "{0} length 4 to 32, contain digits, letters, special characters",
		"zh": "{0} 长度在 4 到 32 位，包含数字、字母、特殊字符",
	}
}

func (n *username) Validate(fl interface{}) error {
	v, ok := fl.(string)
	if ok {
		if len(v) < 4 || len(v) > 32 {
			return errorx.New("length 4 to 32")
		}

		if RegUsername.MatchString(v) {
			return nil
		}
	}

	return errorx.New("only digits, letters, special characters")
}
