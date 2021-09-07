package i18n

import (
	"github.com/spf13/viper"
	"github.com/windrivder/gopkg/errorx"
)

type Options struct {
	Locale string `json:"Locale"`
}

func NewOptions(v *viper.Viper) (o Options, err error) {
	if err = v.UnmarshalKey("i18n", &o); err != nil {
		return o, errorx.Wrap(err, "unmarshal i18n option error")
	}

	return o, err
}
