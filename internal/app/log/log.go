package log

import (
	"os"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
)

var logger *Log

const (
	FileName   = "./logs/"
	MaxSize    = 100 //日志文件长度大小 m
	MaxBackups = 30  //保留最近多少个日志文件
	MaxAge     = 7   //日志保留的最大天数(只保留最近多少天的日志)
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
)

type Log struct {
	// 日志级别
	Level string

	//logger
	logger *zap.Logger
	//sugar logger
	sugaredLogger *zap.SugaredLogger

	//其他参数
	Opts logOptions
}

//配置参数 可默认
type logOptions struct {
	// 是否切割记录文件
	EnableFileWrite bool
	EnableDbWrite   bool //开启存储到数据库

	Caller     bool   //显示行号
	ServerName string //服务名称
	Version    string //版本号
	FileName   string //文件路径位置
	MaxSize    int    //每个文件的最大尺寸 M
	MaxBackups int    //最多保存多少个备份
	MaxAge     int    //文件最多保存多少天
	Compress   bool   //是否压缩
}

//
type LogOption interface {
	apply(*logOptions)
}

type tempFunc func(*logOptions)

var _ LogOption = (*funcLogOption)(nil)

type funcLogOption struct {
	f tempFunc
}

func (flo *funcLogOption) apply(logOptions *logOptions) {
	flo.f(logOptions)
}

func newFuncLogOption(tempF tempFunc) *funcLogOption {
	return &funcLogOption{f: tempF}
}

var DefaultLogOpts = logOptions{ //默认
	EnableFileWrite: false,
	Caller:          true,

	FileName:   FileName + makeFileName(),
	MaxSize:    MaxSize,
	MaxBackups: MaxBackups,
	MaxAge:     MaxAge,
	Compress:   true,
}

func Init(level string /*日志级别*/, opts ...LogOption) {
	log := &Log{
		Level: level,
		Opts:  DefaultLogOpts,
	}
	for _, opt := range opts {
		opt.apply(&log.Opts)
	}
	build(log)
	logger = log
}

func build(log *Log) {
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   log.Opts.FileName,
		MaxSize:    log.Opts.MaxSize,
		MaxBackups: log.Opts.MaxBackups,
		MaxAge:     log.Opts.MaxAge,
		Compress:   log.Opts.Compress,
	})

	var level zapcore.Level
	switch log.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	default:
		level = zap.InfoLevel
	}
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.TimeKey = "time"
	cnf := zapcore.NewJSONEncoder(conf)
	//日志输出方式
	var wss []zapcore.WriteSyncer
	//输出到控制台
	wss = append(wss, zapcore.AddSync(os.Stdout))
	if log.Opts.EnableFileWrite {
		wss = append(wss, ws) //日志切割文件
	}
	if log.Opts.EnableDbWrite { //数据库
		wss = append(wss, zapcore.AddSync(hk))
	}
	writeSs := zapcore.NewMultiWriteSyncer(wss...)
	core := zapcore.NewCore(
		cnf,
		writeSs,
		level,
	)

	//zap option 参数
	var zapOpts []zap.Option
	////增加日志生成时间戳
	//zapOpts = append(zapOpts, zap.Fields(zap.Int64("time_unix", time.Now().Unix())))
	if log.Opts.Version != "" {
		zapOpts = append(zapOpts, zap.Fields(zap.String("version", log.Opts.Version)))
	}
	if log.Opts.ServerName != "" {
		zapOpts = append(zapOpts, zap.Fields(zap.String("server_name", log.Opts.ServerName)))
	}

	log.logger = zap.New(
		core,
		zapOpts...,
	)

	if log.Opts.Caller {
		log.logger = log.logger.WithOptions(
			zap.Development(), //开启文件和行号
			zap.AddCaller(),   //开启开发模式，堆栈跟踪
			zap.AddCallerSkip(1),
		)
	}
	//sugaredLogger
	log.sugaredLogger = log.logger.Sugar()
}

//开启调试
func WithCaller(enable bool) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.Caller = enable
	})
}

//文件位置 "./logs/"
func WithFilePath(filePath string) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.FileName = filePath + makeFileName()
	})
}

//WithEnableWriteFile 是否开启文件写入
func WithEnableWriteFile(enable bool) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.EnableFileWrite = enable
	})
}

//WithEnableDbWrite 是否开启写入数据库
func WithEnableDbWrite(enable bool) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.EnableDbWrite = enable
	})
}

//日志文件大小 单位m
func WithMaxSize(maxSize int) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.MaxSize = maxSize
	})
}

//日志最多保存多少个备份 文件个数
func WithMaxBackups(maxBackups int) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.MaxBackups = maxBackups
	})
}

//文件最多保存多少天数
func WithMaxAge(maxAge int) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.MaxAge = maxAge
	})
}

//是否开启压缩
func WithCompress(compress bool) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.Compress = compress
	})
}

//版本号
func WithVersion(version string) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.Version = version
	})
}

//服务名称
func WithServerName(name string) LogOption {
	return newFuncLogOption(func(options *logOptions) {
		options.ServerName = name
	})
}

//flush buffered
func Sync() {
	if logger == nil {
		return
	}
	_ = logger.sugaredLogger.Sync()
}

//生成文件名称
func makeFileName() string {
	fileName := time.Now().Format("20060102") + ".log"
	return fileName
}
