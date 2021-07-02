package logx

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zapLogger logger struct
type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

// newZapLogger new zap logger
func newZapLogger(o Options) (Logger, error) {
	if log != nil {
		return log, nil
	}

	encoder := getJsonEncoder()
	log = &zapLogger{
		sugaredLogger: zap.New(
			zapcore.NewTee(
				zapcore.NewCore(
					encoder,
					getLogWriter(o.File, o),
					getLogLevel(o.Level),
				),
				zapcore.NewCore(
					encoder,
					zapcore.AddSync(os.Stdout),
					getLogLevel(o.Level),
				),
			),
			zap.AddCaller(),
			zap.Development(),
			zap.AddCallerSkip(2),
		).Sugar(),
	}

	return log, nil
}

// Debug logger
func (l *zapLogger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

// Info logger
func (l *zapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

// Warn logger
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

// Error logger
func (l *zapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

// Fatal logger
func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

// Debugf logger
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

// Infof logger
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

// Warnf logger
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

// Errorf logger
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

// Fatalf logger
func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

// Panicf logger
func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

// WithFields logger
func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}

// getJsonEncoder
func getJsonEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.CallerKey = "file"
	encoderConfig.StacktraceKey = "trace"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.MillisDurationEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter
func getLogWriter(filename string, o Options) zapcore.WriteSyncer {
	writeSyncer := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    o.MaxSize,
		MaxAge:     o.MaxAge,
		MaxBackups: o.MaxBackups,
		Compress:   o.Compress,
	}

	return zapcore.AddSync(writeSyncer)
}

// getLogLevel
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}

	return zapcore.FatalLevel
}
