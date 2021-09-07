package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/logx"
	"github.com/windrivder/gopkg/rpc/rest/middleware"
)

type Options struct {
	Name            string        `json:"Name"`
	Mode            string        `json:"Mode"`
	Host            string        `json:"Host"`
	Port            int           `json:"Port"`
	CertFile        string        `json:"CertFile"`
	KeyFile         string        `json:"KeyFile"`
	ShutdownTimeout time.Duration `json:"ShutdownTimeout"`
	ClientTimeout   time.Duration `json:"ClientTimeout"`
	Secret          string        `json:"Secret"`
	Expired         time.Duration `json:"Expired"`
	Locale          string        `json:"Locale"`
}

func NewOptions(v *viper.Viper) (o Options, err error) {
	if err = v.UnmarshalKey("Rest", &o); err != nil {
		return o, errorx.Wrap(err, "unmarshal rest option error")
	}

	return o, err
}

type Server struct {
	*echo.Echo
	o   Options
	log *logx.Logger
}

func New(o Options, log *logx.Logger, fn HandlerRoutersFunc) (IServer, func(), error) {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	binder, err := NewBinder(o.Locale)
	if err != nil {
		return nil, nil, err
	}
	e.Binder = binder

	s := &Server{o: o, log: log, Echo: e}

	fn(s)

	return s, func() { s.Stop() }, nil
}

func (s *Server) Start() (err error) {
	addr := fmt.Sprintf("%s:%d", s.o.Host, s.o.Port)
	s.log.Info().Str("addr", addr).Msg("http server starting...")

	go func() {
		s.Echo.Server.Addr = addr

		if s.o.CertFile == "" && s.o.KeyFile == "" {
			err = s.Echo.Server.ListenAndServe()
		} else {
			err = s.Echo.Server.ListenAndServeTLS(s.o.CertFile, s.o.KeyFile)
		}

		if err != nil && err != http.ErrServerClosed {
			s.log.Fatal().Msgf("start http server err: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	s.log.Info().Msg("http server stopping...")

	timeout := time.Second * s.o.ShutdownTimeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.Echo.Shutdown(ctx)
}

func (s *Server) ErrorHandler(eh HTTPErrorHandler) {
	s.Echo.HTTPErrorHandler = eh
}

var ProviderSet = wire.NewSet(New, NewOptions)
