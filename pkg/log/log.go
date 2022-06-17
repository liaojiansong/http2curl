package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

type Field = zapcore.Field
type config struct {
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level             string   `json:"level"              mapstructure:"level"`
	Format            string   `json:"format"             mapstructure:"format"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`
	Development       bool     `json:"development"        mapstructure:"development"`
	Name              string   `json:"name"               mapstructure:"name"`
}

type Options func(o *config)

var std *zap.Logger

func SetLevel(level string) Options {
	return func(o *config) {
		o.Level = level
	}
}

func Init(opts ...Options) {
	o := &config{
		Level:             zapcore.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            "console",
		EnableColor:       false,
		Development:       false,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	for _, opt := range opts {
		opt(o)
	}
	//encoderConfig := zapcore.EncoderConfig{
	//	MessageKey:    "message",
	//	LevelKey:      "level",
	//	TimeKey:       "timestamp",
	//	NameKey:       "logger",
	//	CallerKey:     "caller",
	//	StacktraceKey: "stacktrace",
	//	LineEnding:    zapcore.DefaultLineEnding,
	//	EncodeCaller:  zapcore.ShortCallerEncoder,
	//
	//}

	//loggerConfig := &zap.Config{
	//	Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
	//	Development:       o.Development,
	//	DisableCaller:     o.DisableCaller,
	//	DisableStacktrace: o.DisableStacktrace,
	//	Sampling: &zap.SamplingConfig{
	//		Initial:    100,
	//		Thereafter: 100,
	//	},
	//	Encoding: o.Format,
	//	//EncoderConfig:    encoderConfig,
	//	OutputPaths:      o.OutputPaths,
	//	ErrorOutputPaths: o.ErrorOutputPaths,
	//}

	var err error
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("init zap logger fail; err:%v", err.Error())
	}
	std = logger
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...Field) {
	std.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...Field) {
	std.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...Field) {
	std.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...Field) {
	std.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanic(msg string, fields ...Field) {
	std.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...Field) {
	std.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...Field) {
	std.Fatal(msg, fields...)
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() error {
	return std.Sync()
}
