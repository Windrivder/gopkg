package wxopen

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/rpc/rest"
)

type Options struct {
	Debug     bool   `json:"Debug"`
	AppId     string `json:"AppId"`
	AppSecret string `json:"AppSecret"`
	Token     string `json:"Token"`
	AesKey    string `json:"AesKey"`
}

func NewOptions(v *viper.Viper) (o Options, err error) {
	if err = v.UnmarshalKey("Wxopen", &o); err != nil {
		return o, errorx.Wrap(err, "unmarshal wxopen option error")
	}

	return o, err
}

type Platform struct {
	Options Options
	Client  *rest.Client
}

var ProviderSet = wire.NewSet(NewOptions,
	wire.NewSet(wire.Struct(new(Platform), "*")),
)
