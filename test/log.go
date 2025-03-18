package test

import (
	"fmt"
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log() {
	// logger:= configLog(configEncoder())
	logger:= configLog(configEncoder(), writeSync())
	defer logger.Sync()


	logger.Error("This is an error log")
}
func configLog(encoder zapcore.Encoder, writeSync zapcore.WriteSyncer) *zap.Logger{
	core:= zapcore.NewCore(encoder, writeSync, zapcore.InfoLevel)
	logger:= zap.New(core, zap.AddCaller()) // AddCaller() -> show file.go:line
	return logger
}

func configEncoder() zapcore.Encoder{
	encoder:= zap.NewProductionEncoderConfig()

	encoder.EncodeLevel = zapcore.CapitalLevelEncoder // info -> INFO
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder // time -> 2021-09-09T12:12:12.123456Z
	encoder.TimeKey = "time" // ts -> time
	encoder.EncodeCaller = zapcore.ShortCallerEncoder // /path/to/file.go:line -> file.go:line

	return zapcore.NewJSONEncoder(encoder)
} 
 
func writeSync() zapcore.WriteSyncer{
	logFile:= &lumberjack.Logger{
		Filename: "logger/log.txt",
		MaxSize: 10, // 10MB
		MaxBackups: 3, // số file backup
		MaxAge: 10, // 10 days
		Compress: true, // Nén file cũ
	}
	// fileWrite, err := os.OpenFile("logger/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil{
	// 	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	// }
	async:= zapcore.AddSync(logFile) // async log ra file
	consolog:= zapcore.AddSync(os.Stdout) // console log ra terminal

	return zapcore.NewMultiWriteSyncer(async, consolog)
}

func getWriteSync(filePath string) zapcore.WriteSyncer{
	fileWrite, err:= os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil{
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWrite), zapcore.AddSync(os.Stdout))
}

func handMakeWriteSync() zapcore.WriteSyncer{
	logFilePaht:= "log/log.txt"
	maxSize:= int64(10*1024*1024) // 10MB
	maxBackups:= 3

	fileInfo, err:= os.Stat(logFilePaht)
	if err != nil{
		if  os.IsNotExist(err){
			return getWriteSync(logFilePaht)
		}
	}
	if fileInfo.Size() < maxSize {
		return getWriteSync(logFilePaht)
	}
	if fileInfo.Size() >= maxSize{
		for i:= 1; i<= maxBackups; i++{
			newPath:= fmt.Sprintf("_%d.txt", i)
			logFilePaht:= strings.Replace(logFilePaht, ".txt", newPath, 1)
			if fileInfo, err:= os.Stat(logFilePaht); err != nil{
				if os.IsNotExist(err){
					return getWriteSync(logFilePaht)
				}
				if fileInfo.Size() < maxSize{
					return getWriteSync(logFilePaht)
				}
			}
		}
	}

	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))

}




func configZap() *zap.Logger{
	config:= zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder	
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder


	logger, _:= config.Build( zap.AddCaller() )
	return logger
}
