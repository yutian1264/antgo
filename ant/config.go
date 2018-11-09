package ant

import (
	"log"
	"path/filepath"
	"reflect"
	"antgo/antlib"

	"strings"
)

var (
	//基础配置文件
	BConf *Conf
	//项目访问路径
	AppPath string
	//运行模式 dev prod
	RunMode string
	//项目目录
	WorkPath string
	//支持 view 层解析格式
	TplExt = []string{"tpl", "html", "htm"}
)

//配置构造体
type Conf struct {
	Host      string              //运行域名
	Port      int64               //运行端口
	AppName   string              //项目名称
	RunMode   string              //运行模块
	TplPATH   string              //模板路径
	TplExt    string              //模板后缀
	StaticDir map[string]string   //静态文件目录

}


//配置初始化
func init() {
	BConf = newConf()
	confPath := filepath.Join("", "conf", "ant.conf")
	parseConfig(confPath)
	if TplExtCheck(BConf.TplExt) == false {
		log.Fatal("`tpl_ext` can only be html,htm,tpl")
	}
	if BConf.RunMode == DEV {
		log.Println(BConf)
	}
}

func newConf() *Conf {
	return &Conf{
		Port:      8080,
		AppName:   "ant",
		RunMode:   DEV,
		TplPATH:   "views",
		TplExt:    "tpl",
		StaticDir: map[string]string{"static": "static"},
	}
}

func parseConfig(confPath string){
	//文件读取
	antlib.AntInit(confPath)
	for _, i := range []interface{}{BConf} {
		assignSingleConfig(i)
	}

	if sd := antlib.AppConfig.String("StaticDir"); sd != "" {
		BConf.StaticDir = map[string]string{}
		sds := strings.Fields(sd)
		for _, v := range sds {
			if url2fsmap := strings.SplitN(v, ":", 2); len(url2fsmap) == 2 {
				BConf.StaticDir["/"+strings.Trim(url2fsmap[0], "/")] = url2fsmap[1]
			} else {
				BConf.StaticDir["/"+strings.Trim(url2fsmap[0], "/")] = url2fsmap[0]
			}
		}
	}
}

func assignSingleConfig(p interface{}){

	pt := reflect.TypeOf(p)
	if pt.Kind() != reflect.Ptr {
		return
	}
	pt = pt.Elem()
	if pt.Kind() != reflect.Struct {
		return
	}
	pv := reflect.ValueOf(p).Elem()

	for i := 0; i < pt.NumField(); i++ {
		pf := pv.Field(i)
		if !pf.CanSet() {
			continue
		}
		name := pt.Field(i).Name
		switch pf.Kind() {
		case reflect.String:
			pf.SetString(antlib.AppConfig.DefaultString(name, pf.String()))
		case reflect.Int, reflect.Int64:
			pf.SetInt(antlib.AppConfig.DefaultInt64(name, pf.Int()))
		case reflect.Bool:
			pf.SetBool(antlib.AppConfig.DefaultBool(name, pf.Bool()))
		case reflect.Struct:
		default:
			//do nothing here
		}
	}
}

//模板后缀检查
func TplExtCheck(ext string) bool {

	for _, v := range TplExt {

		if ext == v {
			return true
		}
	}
	return false
}
func (c *Conf) getConf(key string) interface{} {
	val := reflect.ValueOf(c)
	v := val.Elem().FieldByName(key)
	//如果存在对应的字段
	if v.IsValid() {
		return v.Interface()
	} else {
		return nil
	}

}

func SetViewsPath(path string) *Conf {
	BConf.TplPATH = path
	return BConf
}

func SetStaticPath(url string, path string) *Conf {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	if url != "/" {
		url = strings.TrimRight(url, "/")
	}
	BConf.StaticDir[url] = path
	return BConf
}


func DelStaticPath(url string) *Conf {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	if url != "/" {
		url = strings.TrimRight(url, "/")
	}
	delete(BConf.StaticDir, url)
	return BConf
}
