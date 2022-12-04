package core

import (
	"amazing/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
	"time"
)

var zapInitLock sync.Mutex

func InitZap() {
	defer zapInitLock.Unlock()
	zapInitLock.Lock()

	lumberJackLogger := &lumberjack.Logger{
		Filename:   global.Global.App.LogFileName,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	writeSyncer := zapcore.AddSync(lumberJackLogger)

	// json
	//jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// console
	//jsonEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	// 时间优化
	encoderConfig := zap.NewProductionEncoderConfig()

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(level.CapitalString())
	}
	encoderConfig.ConsoleSeparator = " "
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = customLevelEncoder

	jsonEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(jsonEncoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)
	// log := zap.New(core, zap.AddCaller())
	global.Log = zap.New(core, zap.AddCaller()).Sugar()
}
