# pkg
go语言工具包

## xzap
基于zap日志包改造的日志工具
- 支持按小时自动切割日志文件
- 自动把error级别日志单独输出到 outputPath.log.error 文件

```$xslt
out := []string{"stderr","/ssd/outputPath.log"}
zlog,err := xzap.NewZLog(out,zapcore.DebugLevel)
if err != nil {
    panic(err)
}

for {
    zlog.Debug("hello debug")
    zlog.Info("hello info")
    zlog.Error("hello error")
    time.Sleep(10 * time.Second)
}

//输出
/ssd/outputPath.log
/ssd/outputPath.log.error
```