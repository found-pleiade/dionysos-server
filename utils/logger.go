package utils

import (
	"log"
	"os"

	c "github.com/Brawdunoir/dionysos-server/variables"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

// InitLogger initializes the Logger.
// The Logger is then available in the utils package.
func InitLogger() error {
	config, defaultLogLevel := createConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, err := os.OpenFile("dionysos.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
	return nil
}

func createConfig() (zapcore.EncoderConfig, zapcore.Level) {
	switch c.Environment {
	case c.ENVIRONMENT_TESTING:
		return zap.NewDevelopmentEncoderConfig(), zapcore.InfoLevel
	case c.ENVIRONMENT_DEVELOPMENT:
		return zap.NewDevelopmentEncoderConfig(), zapcore.DebugLevel
	case c.ENVIRONMENT_PRODUCTION:
		return zap.NewProductionEncoderConfig(), zapcore.WarnLevel
	default:
		log.Fatal("Unknown environment: " + c.Environment)
		return zapcore.EncoderConfig{}, 0
	}
}
