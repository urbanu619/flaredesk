package core

import (
	"fmt"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	config2 "go_server/base/config"
	"io"
	"os"
	"strings"
	"time"
)

const (
	DefaultLog string = "app"
)

var (
	loggerMap = make(map[string]*zap.Logger)
	Log       *zap.SugaredLogger
)

func init() {
	tags := []string{DefaultLog}
	for _, lg := range tags {
		loggerMap[lg] = newLogger(lg)
	}
	Log = loggerMap[DefaultLog].Sugar()
}

func conf() *config2.Zap {
	if loggerConfig == nil {
		loggerInit()
	}
	return loggerConfig
}

var loggerConfig *config2.Zap

func loggerInit() {
	defaultConf := &config2.Zap{
		TagName:       "app",
		Level:         "debug",
		Prefix:        "AMS_",
		Format:        "_json",
		Director:      "logs",
		EncodeLevel:   "LowercaseLevelEncoder",
		StacktraceKey: "error",
		MaxAge:        3,
		RotationSize:  100,
		RotationCount: 3,
		ShowLine:      true,
		LogInConsole:  true,
		LogOutputFile: true,
	}
	zapConf := config2.EnvConf().Zap
	if zapConf == nil {
		loggerConfig = defaultConf
	} else {
		loggerConfig = zapConf
	}
}

type logConfTags struct {
	Tag         string `json:"tag" yaml:"tag"`
	FileName    string `json:"filename" yaml:"filename"`
	ErrFileName string `json:"errFileName" yaml:"errFileName"`
}

func getWriter(filename string) io.Writer {
	// 生成 rotate_logs 的Logger 实际生成的文件名 demo.log.YY mm dd HH
	hook, err := rotateLogs.New(
		strings.Replace(filename, ".log", "", -1)+".%Y%m%d%H.log",
		rotateLogs.WithLinkName(filename),
		// 根据文件大小切割日志
		rotateLogs.WithRotationSize(1024*1024*conf().LogSignSize()), // 每个日志文件大小设置:100M
		rotateLogs.WithRotationTime(time.Minute*1),                  // 日志轮询周期 默认1分钟：60秒
		rotateLogs.WithRotationCount(conf().LogSaveCount()),         // 日志保留份数
		// 根据时间分割日志
		// rotateLogs.WithMaxAge(time.Hour*24*2), // 保存2天日志
	)
	if err != nil {
		panic(err)
	}
	return hook
}

// 设置日志格式 非json可进行堆栈跟踪

func getEncoder() zapcore.Encoder {
	if conf().Format == "json" || conf().Format == "" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		EncodeLevel:    conf().ZapEncodeLevel(),
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  conf().StacktraceKey, // "error"
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006/01/02-15:04:05.000"))
}

func NewLogger(tag string, args ...interface{}) *zap.Logger {
	if loggerMap[tag] != nil {
		return loggerMap[tag]
	}
	return newLogger(tag, args...)
}

func newLogger(tag string, args ...interface{}) *zap.Logger {
	savePath := fmt.Sprintf("./%s/", conf().Director)
	if len(args) > 0 {
		savePath = fmt.Sprintf("./%s/", args[0])
	}
	tagConf := logConfTags{Tag: tag, FileName: fmt.Sprintf("%s.log", tag),
		ErrFileName: fmt.Sprintf("%s_err.log", tag)}

	encoder := getEncoder()
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= conf().TransportLevel()
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(savePath + tagConf.FileName)
	errorWriter := getWriter(savePath + tagConf.ErrFileName)
	cores := make([]zapcore.Core, 0)
	if conf().LogOutputFile {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel))
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel))
	}
	if conf().LogInConsole {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), infoLevel))
	}
	// 最后创建具体的Logger
	core := zapcore.NewTee(cores...)
	// 开启开发模式，堆栈跟踪
	//caller := zap.AddCaller()
	// 	config.DisableCaller = true // 禁用调用者信息（代码行）
	development := zap.Development()
	filed := zap.Fields(zap.String("Tag", string(tagConf.Tag)))
	stackTrace := zap.AddStacktrace(zap.ErrorLevel)      // 当错误等级error时 触发堆栈跟踪
	log := zap.New(core, stackTrace, development, filed) // caller,
	if conf().ShowLine {
		log = log.WithOptions(zap.AddCaller())
	}
	loggerMap[tagConf.Tag] = log
	return loggerMap[tagConf.Tag]
}
