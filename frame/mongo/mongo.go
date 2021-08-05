package mongo

import (
	"context"
	"net/url"

	"github.com/google/wire"
	"github.com/qiniu/qmgo"
	"github.com/spf13/viper"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/logx"
)

type Options struct {
	Host       string
	Username   string
	Password   string
	Database   string
	Collection string
}

// NewOptions
func NewOptions(v *viper.Viper) (o Options, err error) {
	if err = v.UnmarshalKey("MongoDB", &o); err != nil {
		return o, errorx.Wrap(err, "unmarshal mongodb server option error")
	}

	return o, nil
}

func New(o Options) (client *qmgo.Client, cleanFunc func(), err error) {
	dbUri := &url.URL{
		Scheme: "mongodb",
		Host:   o.Host,
		User:   url.UserPassword(o.Username, o.Password),
	}

	config := &qmgo.Config{
		Uri:      dbUri.String(),
		Database: o.Database,
		Coll:     o.Collection,
		Auth:     &qmgo.Credential{Username: o.Username, Password: o.Password},
	}

	client, err = qmgo.NewClient(context.Background(), config)
	if err != nil {
		return nil, nil, err
	}

	// ping
	if err := client.Ping(10); err != nil {
		return nil, nil, err
	}

	cleanFunc = func() {
		if err := client.Close(context.Background()); err != nil {
			logx.Error().Msgf("close mongodb client: %+v", err)
		}
	}

	return client, cleanFunc, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
