// main.go
package main

import (
	"fmt"
	"ftpsync/myfsnotify"
	"ftpsync/sftphandler"
	"ftpsync/utils"
)

func main() {

	watch := myfsnotify.NewNotifyFile()
	for _, v := range utils.GlobalObject.Sync {
		watch.WatchDir(v.Sourcepath, v.Targetpath) //添加监控目录
	}
	sftpClient := sftphandler.NewSftpHandler()
	go func(*myfsnotify.NotifyFile) {
		for {
			select {
			case path := <-watch.Path:
				{
					sftpClient.Upload(path.SoucePath, path.TargetPath)
					fmt.Println("返回路径 : ", path)
				}

			}
		}
	}(watch)

	select {}
	return
}
