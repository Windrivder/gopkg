package validators

import (
	"github.com/windrivder/gopkg/container/typex"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/util/valid"
)

type mobile struct{}

func Mobile() valid.Validator {
	return &mobile{}
}

func (m *mobile) Name() string {
	return "mobile"
}

func (m *mobile) Trans() typex.DictStrs {
	return typex.DictStrs{
		"en": "{0} must be a mobile format",
		"zh": "{0} 必须是一个手机格式",
	}
}

func (m *mobile) Validate(i interface{}) error {
	mobile, ok := i.(string)
	if ok {
		if RegMobile.MatchString(mobile) {
			return nil
		}
	}

	return errorx.New("must be a mobile format")
}

func MobileValidate(mobile string) error {
	return Mobile().Validate(mobile)
}
