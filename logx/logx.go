package logx

import (
	"sync"

	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/windrivder/gopkg/errorx"
)

// instance of Logger
var (
	once sync.Once
	log  Logger
)

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// Logger is our contract for the logger
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	WithFields(keyValues Fields) Logger
}

// Options is log configuration struct
type Options struct {
	Level      string `json:"Level"`
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

func New(o Options) (Logger, error) {
	if log == nil {
		once.Do(func() {
			newZapLogger(o)
		})
	}

	return log, nil
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

var ProviderSet = wire.NewSet(New, NewOptions)
