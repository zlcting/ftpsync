// main.go
package main

import (
	"fmt"
	"ftpsync/myfsnotify"
	"ftpsync/utils"
)

func main() {

	watch := myfsnotify.NewNotifyFile()
	for _, v := range utils.GlobalObject.Sync {
		watch.WatchDir(v.Sourcepath, v.Targetpath) //添加监控目录
	}

	go func(*myfsnotify.NotifyFile) {
		for {
			select {
			case path := <-watch.Path:
				{
					fmt.Println("返回路径 : ", path)
				}

			}
		}
	}(watch)

	select {}
	return
}
