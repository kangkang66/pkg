package browser

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
	"bytes"
	"mime/multipart"
	"os"
	"io"
)

type Browser struct {
	header map[string]string
	cookies []*http.Cookie
	client *http.Client
}

//初始化
func NewBrowser() *Browser {
	hc := &Browser{}
	hc.header = make(map[string]string)
	hc.client = &http.Client{}
	//为所有重定向的请求增加cookie
	hc.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 {
			for _,v := range hc.GetCookie() {
				req.AddCookie(v)
			}
		}
		return nil
	}
	return hc
}

//设置代理地址
func (self *Browser) SetProxyUrl(proxyUrl string)  {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyUrl)
	}
	transport := &http.Transport{Proxy:proxy}
	self.client.Transport = transport
}

func (self *Browser) AddHeader(header map[string]string)  {
	for k,v := range header{
		self.header[k] = v
	}
}

//设置请求cookie
func (self *Browser) AddCookie(cookies []*http.Cookie)  {
	self.cookies = append(self.cookies, cookies...)
}

//获取当前所有的cookie
func (self *Browser) GetCookie() ([]*http.Cookie) {
	return self.cookies
}

//发送Get请求
func (self *Browser) Get(requestUrl string) ([]byte, error) {
	return self.makeRequest("GET", requestUrl,nil)
}

//发送Post请求
func (self *Browser) Post(requestUrl string, params map[string]string) ([]byte, error) {
	header := map[string]string{"Content-Type":"application/x-www-form-urlencoded"}
	self.AddHeader(header)
	postData := self.encodeParams(params)

	return self.makeRequest("POST", requestUrl, strings.NewReader(postData))
}

//发送Post json请求
func (self *Browser) PostRaw(requestUrl string, json []byte) ([]byte, error) {
	header := map[string]string{"Content-Type":"application/json"}
	self.AddHeader(header)

	return self.makeRequest("POST", requestUrl, bytes.NewBuffer(json))
}

//上传文件
func (self *Browser) UploadFile(requestUrl, fieldName, filename string, params map[string]string)  ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil,err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		return nil,err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil,err
	}

	for key, val := range params {
		_ = bodyWriter.WriteField(key, val)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	header := map[string]string{"Content-Type":contentType}
	self.AddHeader(header)
	return self.makeRequest("POST",requestUrl,bodyBuf)
}


func (self *Browser) makeRequest(method , requestUrl string, body io.Reader) ([]byte, error) {
	request,_ := http.NewRequest(method, requestUrl, body)
	self.setHeader(request)
	self.setRequestCookie(request)
	response,err := self.client.Do(request)
	if err!=nil{
		return nil,err
	}
	defer response.Body.Close()

	respCks := response.Cookies()
	self.AddCookie(respCks)

	data, _ := ioutil.ReadAll(response.Body)
	return data, nil
}

//为请求设置header
func (self *Browser) setHeader(request *http.Request)  {
	for k,v := range self.header{
		if len(request.Header.Get(k))==0{
			request.Header.Set(k,v)
		}else {
			request.Header.Del(k)
			request.Header.Set(k,v)
		}
	}
}

//为请求设置 cookie
func (self *Browser) setRequestCookie(request *http.Request)  {
	for _,v := range self.cookies{
		val,_ := request.Cookie(v.Name)
		if(val == nil){
			request.AddCookie(v)
		}
	}
}

//参数 encode
func (self *Browser) encodeParams(params map[string]string) string {
	paramsData := url.Values{}
	for k,v := range params {
		paramsData.Set(k,v)
	}
	return paramsData.Encode()
}
