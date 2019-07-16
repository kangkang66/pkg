# pkg
go语言工具包

## xaes
go版本实现aes的加解密
```
aes := NewAes("4ea5a27becaa2e054b950b3b7c8c95034949e4vb")

val,err := aes.Encrypt([]byte("kangkang"))
fmt.Println(val,err)

dstr :="FtEDCQ8kRMeijBZ+CoJzidB5g902MdGdug58E1H+i7s="
dval,err := aes.Decrypt(dstr)
fmt.Println(string(dval),err)
```

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

## browser
基于http库实现的http客户端,常驻进程模仿浏览器行为，自动处理代理，cookie，头数据
- 可设置代理
- 可设置 cookie
- 自动保存并应用响应的 cookie
- 自动为重新向的请求添加 cookie


```$xslt
browser = curl.NewBrowser();

//设置header
header := map[string]string{
		"Referer":"http://fuli.asia/luyilu/2017/0127/2910.html",
	}
browser.AddHeader(header)

//设置代理
browser.SetProxyUrl("http://221.174.153.63:9999");

//GET请求
response,err := browser.Get(url);

//POST请求
loginUrl := "http://rtpush.com/login"
loginParams := map[string]string{
    "username":"15656073550",
    "password":"e10adc3949ba59abbe56e057f20f883e",
}
content,err := b.Post(loginUrl,loginParams)

//上传文件
_, filename, _, _ := runtime.Caller(0)
f:= path.Join(path.Dir(filename), "1.html")
loginParams := map[string]string{
    "username":"15656073550",
    "password":"e10adc3949ba59abbe56e057f20f883e",
}
upUrl := "http://rtpush.com/api/v1/upload"
content,err = b.UploadFile(upUrl, "file", f, loginParams)
log.Print(string(content),err)
```