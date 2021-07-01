package config

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// New 初始化 viper
func New(path string) (*viper.Viper, error) {
	var (
		err error
		v   = viper.New()
	)

	v.AddConfigPath(".")
	v.SetConfigFile(string(path))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, err
}

var ProviderSet = wire.NewSet(New)
