package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger() *LoggerZap {
	zapLog:= zapcore.NewCore(configEncoderLogger(), writeSync(), getLogLevel())

	return &LoggerZap{zap.New(zapLog, zap.AddCaller())} // zap.AddCaller() -> option để show caller

}

func configEncoderLogger() zapcore.Encoder {
	enCodeConfig:= zap.NewProductionEncoderConfig()

	enCodeConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // info -> INFO
	enCodeConfig.TimeKey = "time" // ts -> time
	enCodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder // time -> 2021-09-09T12:12:12.123456Z
	enCodeConfig.EncodeCaller = zapcore.ShortCallerEncoder // /path/to/file.go:line -> file.go:line

	return zapcore.NewJSONEncoder(enCodeConfig)

}

func writeSync() zapcore.WriteSyncer {
	fileWrite:= lumberjack.Logger{
		Filename:  "log/log.xxx.log",
		MaxSize:   10, // 10MB
		MaxBackups: 10, // 10 file backup
		MaxAge:    10, // 10 days
		Compress:  true, // Nén file cũ
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(&fileWrite),
		zapcore.AddSync(os.Stdout),
	)
}

// debug -> info -> warn -> error -> panic -> fatal
func getLogLevel() zapcore.Level {
	logLevel:= "debug"
	var level zapcore.Level 
	
	switch logLevel {
		case "debug":
			level = zapcore.DebugLevel // debug
		case "info":
			level = zapcore.InfoLevel // thông tin chung
		case "warn":
			level = zapcore.WarnLevel // cảnh báo có thể xảy ra lỗi
		case "error":
			level = zapcore.ErrorLevel // lỗi nhưng ct vẫn chạy
		case "panic":
			level = zapcore.PanicLevel // lỗi và ct dừng lại
		case "fatal":
			level = zapcore.FatalLevel // ct dừng lại
	}

	return level
}