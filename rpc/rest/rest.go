package rest

import (
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/logx"
)

type Options struct {
	Name            string        `json:"name"`
	Mode            string        `json:"mode"`
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	CertFile        string        `json:"cert_file"`
	KeyFile         string        `json:"key_file"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
	ClientTimeout   time.Duration `json:"client_timeout"`
	Secret          string        `json:"secret"`
	Expired         time.Duration `json:"expired"`
}

func NewOptions(v *viper.Viper) (o Options, err error) {
	if err = v.UnmarshalKey("rest", &o); err != nil {
		return o, errorx.Wrap(err, "unmarshal rest option error")
	}

	return o, err
}

type Server struct {
	o      Options
	svr    *http.Server
	logger logx.Logger
}
