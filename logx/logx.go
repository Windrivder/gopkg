package logx

import (
	"io"
	"os"

	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/windrivder/gopkg/errorx"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zerolog.Logger

type (
	Level   = zerolog.Level
	Logger  = zerolog.Logger
	Context = zerolog.Context
	Event   = zerolog.Event
)

type Options struct {
	// 0Debug 1Info 2Warn 3Error 4Fatal 5Panic 6NoLevel 7Disabled 8Trace
	Level      Level  `json:"Level"`
	File       string `json:"File"`
	MaxSize    int    `json:"MaxSize"` // MB
	MaxAge     int    `json:"MaxAge"`  // Day
	MaxBackups int    `json:"MaxBackups"`
	Compress   bool   `json:"Compress"`
}

func NewOptions(v *viper.Viper) (o Options, err error) {
	if err = v.UnmarshalKey("Log", &o); err != nil {
		return o, errorx.Wrap(err, "unmarshal log option error")
	}

	return o, err
}

func New(o Options) (*Logger, error) {
	logger = func() *Logger {
		writers := []io.Writer{zerolog.ConsoleWriter{Out: os.Stderr}}

		if o.File != "" {
			writers = append(writers, &lumberjack.Logger{
				Filename:   o.File,
				MaxSize:    o.MaxSize,
				MaxAge:     o.MaxAge,
				MaxBackups: o.MaxBackups,
				Compress:   o.Compress,
			})
		}

		log := zerolog.New(io.MultiWriter(writers...)).With().Caller().Timestamp().Logger()
		zerolog.SetGlobalLevel(o.Level)
		return &log
	}()

	return logger, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
