package config

import (
	"go.uber.org/zap/zapcore"
	"strings"
)

type Zap struct {
	TagName       string `mapstructure:"tag_name" json:"tag_name" yaml:"tag_name"`                      // 日志名
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                               // 级别
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                            // 日志前缀
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                            // 输出类型 json
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                     // 日志文件夹
	EncodeLevel   string `mapstructure:"encode_level" json:"encode_level" yaml:"encode_level"`          // 编码级
	StacktraceKey string `mapstructure:"stacktrace_key" json:"stacktrace_key" yaml:"stacktrace_key"`    // 堆栈跟踪关键字
	MaxAge        int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`                         // 日志最大保留(天)
	RotationSize  int64  `mapstructure:"rotation_size" json:"rotation_size" yaml:"rotation_size"`       // 单个日志文件大小：M
	RotationCount uint   `mapstructure:"rotation_count" json:"rotation_count" yaml:"rotation_count"`    // 保存日志份数
	ShowLine      bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`                   // 显示行
	LogInConsole  bool   `mapstructure:"log_in_console" json:"log_in_console" yaml:"log_in_console"`    // 输出控制台
	LogOutputFile bool   `mapstructure:"log_output_file" json:"log_output_file" yaml:"log_output_file"` // 输出到文件
}

func (z *Zap) LogSaveCount() uint {
	if z.RotationCount < 3 {
		return uint(3)
	}
	return z.RotationCount
}

func (z *Zap) LogSignSize() int64 {
	if z.RotationSize < 100 {
		return int64(100)
	}
	return z.RotationSize
}

func (z *Zap) ZapEncodeLevel() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func (z *Zap) TransportLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
