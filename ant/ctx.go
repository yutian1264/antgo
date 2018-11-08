package ant

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ctx struct {
	writer  http.ResponseWriter //http 回应写入
	request *http.Request       //http 请求
}

//设置header
func (this *ctx) SetHeader(key, val string) {
	this.writer.Header().Set(key, val)
}

//获取header
func (this *ctx) GetHeader(key string) {
	this.request.Header.Get(key)
}

//输出
func (this *ctx) Echo(result string) error {
	_, err := this.writer.Write([]byte(result))
	return err
}

//json处理 是否缩进格式化输出
func (this *ctx) JSON(data interface{}, hasIndent bool) {
	this.SetHeader("Content-Type", "application/json; charset=utf-8")
	var content []byte
	var err error
	if hasIndent {
		content, err = json.MarshalIndent(data, "", "  ")
	} else {
		content, err = json.Marshal(data)
	}
	if err != nil {
		this.EchoError(err)
		return
	}

	err = this.Echo(string(content))

	if err != nil {
		this.EchoError(err)
	}
}

//页面跳转
func (this *ctx) Redirect(url string) {
	http.Redirect(this.writer, this.request, url, http.StatusFound)
}

//错误输出
func (this *ctx) EchoError(err error) {
	http.Error(this.writer, err.Error(), http.StatusInternalServerError)
}

//运行错误输出 错误处理方式有待考虑
func (this *ctx) RunError(err error) {
	log.Println(err)
	this.EchoError(errors.New("Page wrong"))
}
