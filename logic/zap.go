package logic

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

var logger *zap.Logger
var logFile *os.File

func InitLogger() {
	var err error
	logFile, err = os.OpenFile("ZapLog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("无法创建或打开日志文件：%v", err)
	}
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006/01/02 15:04:05"))
	}
	err = LoadConfig()
	if err != nil {
		return
	}
	logLevel := getLogLevel(Data.Loglevel)

	logCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		//zapcore.AddSync(logFile), //写到日志文件
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(logFile), zapcore.AddSync(os.Stdout)), //同时写到控制台和日志文件
		logLevel,
	)
	logger = zap.New(logCore, zap.AddCaller())
}

func getLogLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "DEBUG":
		return zap.DebugLevel
	case "INFO":
		return zap.InfoLevel
	case "WARN":
		return zap.WarnLevel
	case "ERROR":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
func CloseLogger() {
	if err := logger.Sync(); err != nil {
		log.Printf("错误：无法同步日志记录器：%v", err)
	}
	if err := logFile.Close(); err != nil {
		log.Printf("错误：无法关闭日志文件：%v", err)
	}
}
