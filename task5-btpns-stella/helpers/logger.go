package helpers

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// the type of callback function which recives parameter of error message
type ErrorFunc func(message string, fields ...zapcore.Field)

// level error callback function
var (
	ErrInfoCallback    ErrorFunc = info
	ErrWarningCallback ErrorFunc = warning
	ErrDebugCallback   ErrorFunc = debug
	ErrErrorsCallback  ErrorFunc = errors
	ErrFatalCallback   ErrorFunc = fatal
)

// setup logger if this package called
func init() {
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	config.EncoderConfig = encoderConfig

	var err error
	// log, err = zap.NewProduction(zap.AddCallerSkip(1))
	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		// panic(err)
		log.Error(err.Error())
	}
}

// function logger for info level
func info(message string, fields ...zapcore.Field) {
	log.Info(message, fields...)
	fmt.Print("\n")
}

// function logger for error level
func errors(message string, fields ...zapcore.Field) {
	log.Error(message, fields...)
	fmt.Print("\n")
}

// function logger for debug level
func debug(message string, fields ...zapcore.Field) {
	log.Debug(message, fields...)
	fmt.Print("\n")
}

// function logger for warning level
func warning(message string, fields ...zapcore.Field) {
	log.Warn(message, fields...)
	fmt.Print("\n")
}

// function logger for fatal level, be careful because this log call panic.
func fatal(message string, fields ...zapcore.Field) {
	// log.Panic(message)
	log.Error(message)
}
