package ant

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	// 版本
	VERSION = "1.0"
	// 测试环境
	DEV = "dev"
	// 生产环境
	PROD = "prod"
)

var (
	BApp *App
)

func init() {
	BApp = NewBApp()
}

//创建App构造体
func NewBApp() *App {
	return &App{
		handle: &ControllerRegister{
			Router: make(map[string]*ControllerInfo),
		},
	}
}

//启动框架
func Run() {

	var (
		server_listen string = ""
		err           error
	)
	if BConf.Host != "" && BConf.Port != 0 {
		server_listen = fmt.Sprintf("%s:%d", BConf.Host, BConf.Port)
	}
	log.Println("server listn :", server_listen)
	err = http.ListenAndServe(server_listen, Controller{})
	//如果监听出现问题 输出错误终止运行
	if err != nil {
		log.Fatal(err.Error())
	}
}

//添加路由
func AddRoute(url, FuncName string, c ControllerInterface) {
	BApp.handle.Add(url, FuncName, c)
}

//添加静态文件路径
func AddStaticPath(url, path string) {
	BConf.StaticDir[strings.Trim(url, "/")] = path
}

//获取配置
func GoWebConf(key string) interface{} {
	return BConf.getConf(key)
}
