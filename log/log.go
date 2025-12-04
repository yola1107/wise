package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const timeFmt = "2006/01/02 15:04:05.000"

type Mode int32

const (
	Dev Mode = iota
	Prod
)

type Config struct {
	Mode   Mode
	Level  string
	App    string
	RoomID int64
	Dir    string
}

var (
	mu  sync.Mutex
	log = newLogger(&Config{Mode: Dev, Level: "debug"})
)

func New(cfg *Config) {
	mu.Lock()
	defer mu.Unlock()
	_ = log.Sync()
	log = newLogger(cfg)
}

func newLogger(cfg *Config) *zap.Logger {
	if cfg == nil {
		_, _ = fmt.Fprintln(os.Stderr, "logger: using default development logger with nil config")
		cfg = &Config{Mode: Dev, Level: "debug"}
	}
	if cfg.App == "" {
		cfg.App = "app"
	}
	lv := zap.NewAtomicLevel()
	if err := lv.UnmarshalText([]byte(cfg.Level)); err != nil {
		_ = lv.UnmarshalText([]byte("debug"))
		_, _ = fmt.Fprintf(os.Stderr, "logger: invalid log level %q, defaulting to DEBUG\n", cfg.Level)
	}
	cores := []zapcore.Core{
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encCfg(false)),
			zapcore.Lock(os.Stderr),
			lv,
		),
	}
	if cfg.Mode == Prod {
		name := filepath.Join(cfg.Dir, fmt.Sprintf("%s_%d", cfg.App, cfg.RoomID))
		cores = append(cores, fileCore(name+".log", lv))
		cores = append(cores, fileCore(name+"_error.log", zap.ErrorLevel))
	}
	return zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
}

func fileCore(file string, lv zapcore.LevelEnabler) zapcore.Core {
	w := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg(true)),
		zapcore.AddSync(w),
		lv,
	)
}

func encCfg(file bool) zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format(timeFmt) + "]")
	}
	cfg.EncodeCaller = func(c zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + c.FullPath() + "]")
	}
	cfg.ConsoleSeparator = " "
	if !file {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncodeCaller = zapcore.FullCallerEncoder
	} else {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	return cfg
}

func Debugf(format string, args ...any) {
	log.Debug(fmt.Sprintf(format, args...))
}

func Infof(format string, args ...any) {
	log.Info(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...any) {
	log.Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...any) {
	log.Error(fmt.Sprintf(format, args...))
}

func Sync() { _ = log.Sync() }
