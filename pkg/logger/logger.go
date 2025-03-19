package logger

import (
	"os"

	"github.com/dinhdev-nu/realtime_auth_go/setting"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(logConfg setting.Logger) *LoggerZap {
	zapLogCore:= zapcore.NewCore(configEncoderLogger(), writeSync(logConfg), getLogLevel(logConfg.Level))

	return  &LoggerZap{zap.New(zapLogCore, zap.AddCaller())}// zap.AddCaller() -> option để show caller
}

func configEncoderLogger() zapcore.Encoder {
	enCodeConfig:= zap.NewProductionEncoderConfig()

	enCodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder // info -> INFO
	enCodeConfig.TimeKey = "time" // ts -> time
	enCodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder // time -> 2021-09-09T12:12:12.123456Z
	enCodeConfig.EncodeCaller = zapcore.ShortCallerEncoder // /path/to/file.go:line -> file.go:line

	return zapcore.NewJSONEncoder(enCodeConfig)

}

func writeSync(logConfg setting.Logger) zapcore.WriteSyncer {
	fileWrite:= lumberjack.Logger{
		Filename:  logConfg.File,
		MaxSize:   logConfg.MaxSize, // 10MB
		MaxBackups: logConfg.MaxBackups, // 10 file backup
		MaxAge:    logConfg.MaxAge, // 30 days
		Compress:  logConfg.Compress, // compress
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(&fileWrite),
		zapcore.AddSync(os.Stdout),
	)
}

// debug -> info -> warn -> error -> panic -> fatal
func getLogLevel(levelConfig string) zapcore.Level {
	logLevel:= levelConfig
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
		default:
			level = zapcore.DebugLevel
	}

	return level
}