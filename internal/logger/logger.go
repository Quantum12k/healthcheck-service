package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	Config struct {
		LogFilePath string `yaml:"log_file_path"`
		Debug       bool   `yaml:"debug"`
		MaxSizeMB   int    `yaml:"max_size_mb"`
		MaxBackups  int    `yaml:"max_backups"`
		MaxAgeDays  int    `yaml:"max_age_days"`
		LocalTime   bool   `yaml:"local_time"`
		Compress    bool   `yaml:"compress"`
	}
)

func New(cfg *Config) *zap.SugaredLogger {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	writer := zapcore.AddSync(io.Writer(&lumberjack.Logger{
		Filename:   cfg.LogFilePath,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		LocalTime:  cfg.LocalTime,
		Compress:   cfg.Compress,
	}))

	logLevel := zapcore.InfoLevel

	if cfg.Debug {
		logLevel = zapcore.DebugLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, logLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel),
	)

	return zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
}
