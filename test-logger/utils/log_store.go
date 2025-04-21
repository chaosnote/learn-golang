package utils

import (
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// encodeTimeExcludingTimezone 自定義時間編碼器，排除時區資訊
func encodeTimeExcludingTimezone(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("01-02 15:04:05.000"))
}

// 自定義日誌等級編碼器 (已註解)
// func levelEncoder(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
// 	encoder.AppendString("[" + strings.ToUpper(level.String()[:1]) + "]") // 只取首字母，並加上中括號
// }

//-----------------------------------------------

// LogFields 定義日誌欄位類型
type LogFields map[string]any

/**
 * LogStore 專案日誌儲存介面
 * ∟ ConsoleLogger
 * ∟ FileLogger
 * ∟ ElkLogger (假設未來可能支援 ELK)
 */
type LogStore interface {
	Debug(fields LogFields)
	Info(fields LogFields)
	Warn(fields LogFields)
	Error(fields LogFields)
	Panic(fields LogFields)
}

//-----------------------------------------------

type defaultLogStore struct {
	logger *zap.Logger
}

func (ls *defaultLogStore) Debug(fields LogFields) {
	ls.logger.Debug("", zap.Any("M", fields))
}
func (ls *defaultLogStore) Info(fields LogFields) {
	ls.logger.Info("", zap.Any("M", fields))
}
func (ls *defaultLogStore) Warn(fields LogFields) {
	ls.logger.Warn("", zap.Any("M", fields))
}
func (ls *defaultLogStore) Error(fields LogFields) {
	ls.logger.Error("", zap.Any("M", fields))
}
func (ls *defaultLogStore) Panic(fields LogFields) {
	ls.logger.Error("", zap.Any("M", fields))
}

//-----------------------------------------------

// NewConsoleLogger 建立一個輸出到控制台的 LogStore 實例
func NewConsoleLogger(callerSkip int) LogStore {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.LevelKey = zapcore.OmitKey
	config.EncoderConfig.MessageKey = zapcore.OmitKey
	config.EncoderConfig.EncodeTime = encodeTimeExcludingTimezone
	// config.EncoderConfig.TimeKey = zapcore.OmitKey
	// config.EncoderConfig.CallerKey = zapcore.OmitKey
	config.Encoding = "json"

	logger, err := config.Build(zap.AddCallerSkip(callerSkip))
	if err != nil {
		panic(err)
	}
	return &defaultLogStore{
		logger: logger,
	}
}

// NewFileLogger 建立一個輸出到檔案的 LogStore 實例
func NewFileLogger(directory string, callerSkip int) LogStore {
	hook := lumberjack.Logger{
		Filename:   path.Join(directory, ".log"), // 文件輸出路徑
		MaxSize:    10,                           // 文件最大大小 (MB)
		LocalTime:  true,                         // 使用本地時間
		Compress:   false,                        // 是否壓縮檔案
		MaxAge:     30,                           // 舊檔案保留天數
		MaxBackups: 50,                           // 最多備份檔案數量
	}
	writeSyncer := zapcore.AddSync(&hook)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = zapcore.OmitKey
	encoderConfig.MessageKey = zapcore.OmitKey
	encoderConfig.EncodeTime = encodeTimeExcludingTimezone
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 短檔案路徑

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // zapcore.NewConsoleEncoder(encoderConfig),
		writeSyncer,
		zapcore.InfoLevel,
	)

	return &defaultLogStore{
		logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(callerSkip)),
	}
}
