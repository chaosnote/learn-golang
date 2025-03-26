package utils

import (
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 專案(log)輸出

type Log interface {
	Logger() *zap.Logger
}

type trace struct {
	mu     sync.Mutex
	output *zap.Logger
}

func (t *trace) Logger() *zap.Logger {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.output
}

//-------------------------------------------------------------------------------------------------[Func]

// time_encode 排除(+8)時區
func time_encode(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("01-02 15:04:05.000"))
}

// level_encoder 自定義等級編碼函數
func level_encoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + strings.ToUpper(l.String()[:1]) + "]") // 只取首字母，並加上中括號
}

// GenDebug console mode
func GenDebug(skip int) Log {
	c := zap.NewDevelopmentConfig()
	c.EncoderConfig.LevelKey = zapcore.OmitKey
	c.EncoderConfig.EncodeTime = time_encode
	c.EncoderConfig.CallerKey = zapcore.OmitKey
	// c.Encoding = "json"

	//console, e := c.Build() // default
	_console, e := c.Build(zap.AddCallerSkip(skip))
	if e != nil {
		panic(e)
	}
	return &trace{
		output: _console,
	}
}

// GenFile file mode
func GenFile(dir string, skip int) Log {
	h := lumberjack.Logger{
		Filename:   path.Join(dir, ".log"), // 文件輸出位置
		MaxSize:    10,                     // 文件大小 MB
		LocalTime:  true,                   // 是否使用本地時間
		Compress:   false,                  // 是否壓縮檔案 ( 大小差滿多的，但不確定效能會不會影響很大 )
		MaxAge:     30,                     // 預設值是不刪除舊檔(單位天), 修改為 30 天
		MaxBackups: 50,                     // 保留多少個備份檔( 受限 MaxAge )，預設全保留 ( 500MB )
	}
	ws := zapcore.AddSync(&h)

	ec := zap.NewProductionEncoderConfig()
	ec.LevelKey = zapcore.OmitKey
	// ec.MessageKey = zapcore.OmitKey // 調整為保留 message
	ec.EncodeTime = time_encode
	ec.EncodeCaller = zapcore.FullCallerEncoder // 完整路徑

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(ec), // zapcore.NewConsoleEncoder(ec),
		ws,
		zapcore.InfoLevel,
	)

	return &trace{
		output: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(skip)), // File = zap.New(core, zap.AddCaller()) // default
	}
}

/*
GenComplex (console + file) mode

file rotation

https://github.com/uber-go/zap/blob/master/FAQ.md#does-zap-support-log-rotation

logger clone

https://github.com/uber-go/zap/blob/ed5598a9e26d75aae01f4846ce31b3fe785c57c7/logger.go#L284

zap multi core

https://cloud.tencent.com/developer/article/1645126

about mu

https://go.dev/src/log/log.go#L48
*/
func GenComplex(dir string, skip int, debug bool) Log {
	h := lumberjack.Logger{
		Filename:   path.Join(dir, ".log"), // 文件輸出位置
		MaxSize:    10,                     // 文件大小 MB
		LocalTime:  true,                   // 是否使用本地時間
		Compress:   false,                  // 是否壓縮檔案 ( 大小差滿多的，但不確定效能會不會影響很大 )
		MaxAge:     30,                     // 預設值是不刪除舊檔(單位天), 修改為 30 天
		MaxBackups: 50,                     // 保留多少個備份檔( 受限 MaxAge )，預設全保留 ( 500MB )
	}
	ws := zapcore.AddSync(&h)

	ec := zap.NewProductionEncoderConfig()
	// ec.LevelKey = zapcore.OmitKey
	ec.EncodeTime = time_encode
	ec.EncodeCaller = zapcore.ShortCallerEncoder // ( 長路徑 zapcore.FullCallerEncoder ) ( 短路徑 zapcore.ShortCallerEncoder)
	ec.EncodeLevel = level_encoder

	list := []zapcore.Core{
		zapcore.NewCore(zapcore.NewJSONEncoder(ec), ws, zapcore.InfoLevel),
	}

	if debug {
		list = append(list, zapcore.NewCore(zapcore.NewConsoleEncoder(ec), zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	}

	tee := zapcore.NewTee(list...)

	return &trace{
		output: zap.New(tee, zap.AddCaller(), zap.AddCallerSkip(skip)), // File = zap.New(core, zap.AddCaller()) // default
	}
}
