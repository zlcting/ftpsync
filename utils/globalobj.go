package utils

import (
	"encoding/json"
	"io/ioutil"
)

/*GlobalObj 定义全局常量 用户根据 conf 里的文件conf.json来配置
 */
type GlobalObj struct {
	Name    string
	Version string
	Host    string
	Sftp    GlobalSftpMap    //sftp当前服务器主机
	Sync    []GlobalSyncpMap //同步路径参数
}

//GlobalSftpMap ftp
type GlobalSftpMap struct {
	Host string
	Name string
	Pass string
	Port int
}

//GlobalSyncpMap 路径
type GlobalSyncpMap struct {
	Name       string
	Sourcepath string
	Targetpath string
}

//GlobalObject 全局配置
var GlobalObject *GlobalObj

//Reload 读取用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/conf.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供init方法,默认加载
*/
func init() {
	//初始化GlobalObject变量,设置一些默认值
	GlobalObject = &GlobalObj{
		Name:    "ZinxServerApp",
		Version: "V0.1",
		Host:    "0.0.0.0",
	}
	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
