package antlib

import (
	"os"
)

//文件存在检测
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//是否为文件夹
func IsDir(path string) bool {
	if f, err := os.Stat(path); err == nil {
		if f.IsDir() {
			return true
		}
	}
	return false
}
