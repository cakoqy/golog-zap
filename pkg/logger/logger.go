package logger

import (
	"fmt"
	"io"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Format       string               `toml:"format"`
	Level        zapcore.LevelEnabler `toml:"level"`
	SuppressLogo bool                 `toml:"suppress-logo"`
}

// NewConfig returns a new instance of Config with defaults.
func NewConfig() Config {
	return Config{
		Format: "auto",
		Level:  zapcore.Level(-1),
	}
}

// New return a new instance of zap logger
func (c *Config) New(defaultOutput io.Writer) (*zap.Logger, error) {
	w := defaultOutput
	format := c.Format
	if format == "console" {
		// Disallow the console logger if the output is not a terminal.
		return nil, fmt.Errorf("unknown logging format: %s", format)
	}

	// If the format is empty or auto, then set the format depending
	// on whether or not a terminal is present.
	if format == "" || format == "auto" {
		format = "logfmt"
	}

	// config := zap.NewProductionEncoderConfig()
	config := zapcore.EncoderConfig{
		LevelKey:     "level",
		TimeKey:      "time",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.RFC3339NanoTimeEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(config)
	return zap.New(zapcore.NewCore(
		encoder,
		zapcore.Lock(zapcore.AddSync(w)),
		c.Level,
	)), nil
}

// New return a new instance of zap logger
func New(logFile string) *zap.Logger {
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("file=logFile err=%s", err.Error())
	}
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	config := NewConfig()
	l, _ := config.New(multiLogFile)
	return l
}
