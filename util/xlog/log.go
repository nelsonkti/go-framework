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
func NewLogger(logPath string, opts ...LogOptionFunc) (*Log, error) {
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

	return logger, nil
}

func With(l *Log) *Log {
	// 创建一个新的Log实例
	log := &Log{Logger: l.Logger}
	//// 直接将新Log实例的Logger字段设置为指向原始Log实例的Logger字段的指针
	//
	//fmt.Println("Person 1:", &l)          // 输出 Person 1 的值
	//fmt.Println("Person 2:", &log)        // 输出 Person 2 的值
	//fmt.Println("Person 11:", l.Logger)   // 输出 Person 1 的值
	//fmt.Println("Person 22:", log.Logger) // 输出 Person 2 的值
	return log
}

//
//func NewLogger2(logPath string, opts ...CoreOption) *zap.Logger {
//	w := zapcore.AddSync(&lumberjack.Logger{
//		Filename:   logPath,
//		MaxSize:    500, // megabytes
//		MaxBackups: 3,
//		MaxAge:     28, // days
//	})
//
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	core := zapcore.NewCore(
//		zapcore.NewJSONEncoder(encoderConfig),
//		w,
//		zap.InfoLevel,
//	)
//
//	zap.New(core)
//
//	return zap.New(core)
//}

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
