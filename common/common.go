package common

import (
	"path/filepath"
	"strings"
)

// GetPathLastIndex 获取文件上一级目录
func GetPathLastIndex(files string) string {
	paths, _ := filepath.Split(files)
	// fmt.Println(paths, fileName)      //获取路径中的目录及文件名 E:\data\  test.txt
	// fmt.Println(filepath.Base(files)) //获取路径中的文件名test.txt
	// fmt.Println(path.Ext(files))      //获取路径中的文件的后缀 .txt
	return paths
}

// GetTargetPath 替换字符串
func GetTargetPath(s, old, new string) string {
	s = GetPathLastIndex(s)
	return strings.Replace(s, old, new, 1)
}
