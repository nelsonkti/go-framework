package xlog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CoreOption 核心选项
type LogOptionFunc func(*LogOption)

type LogOption struct {
	Filter       Filter
	GlobalFields []zap.Field
}

// NewLogger 创建并返回一个配置好的zap.Logger实例
func NewLogger(logPath string, opts ...LogOptionFunc) *Log {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	core := &Core{
		Core: zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			w,
			zap.InfoLevel,
		),
		GlobalFields: []zap.Field{},
	}

	logOption := &LogOption{}
	// 应用配置选项
	for _, opt := range opts {
		opt(logOption)
	}

	core.GlobalFields = logOption.GlobalFields

	logger := NewLog(zap.New(core), logOption.Filter)

	return logger
}

func With(l *Log) *Log {
	// 创建一个新的Log实例
	log := &Log{Logger: l.Logger}
	return log
}

// WithGlobalFields 添加全局字段到Core。
func WithGlobalFields(fields ...zap.Field) LogOptionFunc {
	return func(c *LogOption) {
		c.GlobalFields = append(c.GlobalFields, fields...)
	}
}

// WithFilters 添加字段过滤器到Core。
func WithFilters(opts ...FilterOption) LogOptionFunc {
	return func(c *LogOption) {
		options := Filter{
			Key:   make(map[interface{}]struct{}),
			Value: make(map[interface{}]struct{}),
		}
		for _, o := range opts {
			o(&options)
		}
		c.Filter = options
	}
}
