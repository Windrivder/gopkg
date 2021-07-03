package qiniu

import (
	"github.com/google/wire"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
	"github.com/windrivder/gopkg/errorx"
)

type Options struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Scope     string `json:"scope"`
	Expires   uint64 `json:"expires"`
}

func NewOptions(v *viper.Viper) (o Options, err error) {
	if err = v.UnmarshalKey("qiniu", &o); err != nil {
		return o, errorx.Wrap(err, "unmarshal qiniu option error")
	}

	return o, err
}

type Qiniu struct {
	Options   Options
	putPolicy storage.PutPolicy
}

func New(o Options) *Qiniu {
	putPolicy := storage.PutPolicy{
		Scope:   o.Scope,
		Expires: o.Expires,
	}
	return &Qiniu{Options: o, putPolicy: putPolicy}
}

func (q *Qiniu) Token() string {
	mac := qbox.NewMac(q.Options.AccessKey, q.Options.SecretKey)
	upToken := q.putPolicy.UploadToken(mac)

	return upToken
}

var ProviderSet = wire.NewSet(New, NewOptions)
