package ant

import (
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"antgo/antlib"
)

//控制器信息
type ControllerInfo struct {
	controllerType reflect.Type
	controllerName string
	funcName       string //方法名称
}

//控制器注册
type ControllerRegister struct {
	Router map[string]*ControllerInfo
}

//控制器注册添加  路由器添加
func (p *ControllerRegister) Add(url, FuncName string, c ControllerInterface,mappingMethods ...string) {
	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()
	//去掉左边 ／
	if strings.Index(url, "/") != -1 {
		url = strings.TrimLeft(url, "/")
	}
	//检测是否存在对应的方法
	if reflectVal.MethodByName(FuncName).IsValid() == false {
		log.Fatal("'" + FuncName + "' method doesn't exist in the controller " + t.Name())
	}
	//初始化
	route := &ControllerInfo{}
	route.controllerType = t
	route.funcName = FuncName
	route.controllerName = t.Name()

	p.Router[url] = route

}

//重写http Handle interface
func (this Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var url = req.URL.Path

	httpc := &ctx{writer: w, request: req}
	//如果运行模式为 develop
	if BConf.RunMode == DEV {
		log.Println(req.URL.String())
	}
	//去掉左边 /
	if strings.Index(url, "/") != -1 {
		url = strings.TrimLeft(url, "/")
	}
	//查看是否配置对应路由 如果没有 查找静态配置
	if v, ok := BApp.handle.Router[url]; ok != false {
		var param []reflect.Value
		vc := reflect.New(v.controllerType)
		//使用断言方式 调用对应的init方法进行初始化
		execController, ok := vc.Interface().(ControllerInterface)
		if !ok {
			log.Fatal("controller is not ControllerInterface")
		}
		//调用初始化方法
		execController.Init(httpc, v.controllerName, v.funcName)
		//反射调用运行方法
		method := vc.MethodByName(v.funcName)
		method.Call(param)
	} else {
		staticFilePro(w, req, url)
	}
}

func staticFilePro(w http.ResponseWriter, req *http.Request, url_path string) {
	var static_key, static_path string
	path_split := strings.Split(url_path, "/")
	static_key = path_split[0]
	//如果不存在对应静态路由 抛出404
	if route_path, ok := BConf.StaticDir[static_key]; ok == false {
		http.Error(w, "not found page", 404)
	} else {
		//去掉开头的 /
		url_path = strings.TrimLeft(url_path, "/")
		//去掉第一个/ 前面的内容
		file_path := strings.TrimLeft(url_path, static_key)
		//生成文件路径
		static_path = filepath.Join(WorkPath, route_path, file_path)
		//如果为dir 输出404
		if antlib.IsDir(static_path) {
			http.Error(w, "not found page", 404)
			return
		}
		//打印路径
		if BConf.RunMode == DEV {
			log.Println(static_path)
		}
		//运行文件服务 如果不存在为 404 如果存在读取文件 并显示
		http.ServeFile(w, req, static_path)
	}
}
