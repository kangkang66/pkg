package xzap

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"sync"
	"time"
)

type ZLog struct {
	log *zap.Logger
	outputPaths []string
	closeOutputPathsFunc func()
	errorOutPaths []string
	closeErrorOutPathsFunc func()
	level zapcore.Level
}

var once = &sync.Once{}
var zl	*ZLog
var err error


func NewZLog(outputPaths []string, level zapcore.Level) (*ZLog,error){
	if zl != nil {
		return zl,err
	}

	instNo := os.Getenv("INST_NO")
	if instNo=="" {
		instNo = "1"
	}
	errorOutPaths := []string{}
	for k,v := range outputPaths {
		if v == "stderr" || v == "stdin" || v == "stdout" {
			continue
		}
		ps := strings.Split(v,".")
		ps[0] = ps[0] + "-" + instNo
		v = strings.Join(ps,".")

		outputPaths[k] = v
		errorOutPaths = append(errorOutPaths, v + ".error")
	}

	zl = &ZLog{
		outputPaths:outputPaths,
		errorOutPaths:errorOutPaths,
		level:level,
	}
	err = zl.init()
	if err != nil {
		return zl,err
	}
	fmt.Println("new zlog",outputPaths, errorOutPaths)
	once.Do(zl.splitByTime)
	return zl,err
}

func (this *ZLog) init() (err error) {
	var sink zapcore.WriteSyncer
	//初始化日志输入文件
	sink, this.closeOutputPathsFunc, err = zap.Open(this.outputPaths...)
	if err != nil {
		this.closeOutputPathsFunc()
		return err
	}
	allWriter := zapcore.AddSync(sink)
	//初始化错误日志输入文件
	sink, this.closeErrorOutPathsFunc, err = zap.Open(this.errorOutPaths...)
	if err != nil {
		this.closeErrorOutPathsFunc()
		return err
	}
	errorWriter := zapcore.AddSync(sink)
	//初始化日志格式配置
	config := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime: this.timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(config)

	//一次写行为到多个输出端
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, allWriter, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= this.level
		})),
		zapcore.NewCore(encoder, errorWriter, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	)

	this.log = zap.New(core,zap.AddCaller(),zap.AddStacktrace(zap.ErrorLevel))
	this.log = this.log.WithOptions(zap.AddCallerSkip(1))

	return
}

func (this *ZLog) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
}

func (this *ZLog) Debug(msg string, fields ...zap.Field) {
	this.log.Debug(msg, fields ...)
}

func (this *ZLog) Info(msg string, fields ...zap.Field) {
	this.log.Info(msg, fields ...)
}

func (this *ZLog) Error(msg string, fields ...zap.Field) {
	this.log.Error(msg, fields ...)
}

func (this *ZLog) splitByTime() {
	fmt.Println("do splitByTime")
	go func() {
		var lastSplitHour = -1
		for {
			time.Sleep(200 * time.Millisecond)

			//每分钟写入一个测试日志
			/*if time.Now().Second() == 0 {
				this.Debug("zlog")
			}*/

			//整点切换文件
			if time.Now().Minute() == 59 {
				currHour := time.Now().Hour()
				if currHour == lastSplitHour {
					continue
				}
				lastSplitHour = currHour

				for _,file := range this.outputPaths{
					_,err := os.Stat(file)
					if err == nil {
						newFile := file + "." + time.Now().Format("2006-01-02_15")
						err = os.Rename(file, newFile)
						if err != nil {
							fmt.Println(err)
						}else{
							fmt.Println("RenameFile", newFile)
						}
					}
				}
				if currHour == 23 {
					for _,file := range this.errorOutPaths{
						_,err := os.Stat(file)
						if err == nil {
							newFile := file + "." + time.Now().Format("2006-01-02_15")
							err = os.Rename(file, newFile)
							if err != nil {
								fmt.Println(err)
							}else{
								fmt.Println("RenameFile", newFile)
							}
						}
					}
				}

				this.log.Sync()
				this.closeOutputPathsFunc()
				this.closeErrorOutPathsFunc()
				this.init()
			}
		}
	}()
}