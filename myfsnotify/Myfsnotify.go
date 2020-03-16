package myfsnotify

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

//NotifyFile 包的指针结构
type NotifyFile struct {
	watch *fsnotify.Watcher
	Path  chan ActionPath
}

//ActionPath 文件操作
type ActionPath struct {
	Path       string
	ActionType fsnotify.Op
	desc       string
	SoucePath  string
	TargetPath string
}

//NewNotifyFile 返回fsnotify对象指针
func NewNotifyFile() *NotifyFile {
	w := new(NotifyFile)
	w.watch, _ = fsnotify.NewWatcher()
	w.Path = make(chan ActionPath, 10)
	return w
}

//WatchDir 监控目录
func (notifyfile *NotifyFile) WatchDir(dir string, target string) {
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//判断是否为目录，监控目录,目录下文件也在监控范围内，不需要加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = notifyfile.watch.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("监控 : ", path)
		}
		return nil
	})

	go notifyfile.WatchEvent(dir, target) //协程
}

//WatchEvent 监控目录
func (notifyfile *NotifyFile) WatchEvent(dir, target string) {
	for {
		select {
		case ev := <-notifyfile.watch.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("创建文件 : ", ev.Name)
					//获取新创建文件的信息，如果是目录，则加入监控中
					file, err := os.Stat(ev.Name)
					if err == nil && file.IsDir() {
						notifyfile.watch.Add(ev.Name)
						fmt.Println("添加监控 : ", ev.Name)
					}

					go notifyfile.PushEventChannel(ev.Name, fsnotify.Create, "添加监控", dir, target)
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("写入文件 : ", ev.Name)
					go notifyfile.PushEventChannel(ev.Name, fsnotify.Write, "写入文件", dir, target)

				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("删除文件 : ", ev.Name)
					//如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						notifyfile.watch.Remove(ev.Name)
						fmt.Println("删除监控 : ", ev.Name)
					}
					//go notifyfile.PushEventChannel(ev.Name, fsnotify.Remove, "删除文件")

				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					//如果重命名文件是目录，则移除监控 ,注意这里无法使用os.Stat来判断是否是目录了
					//因为重命名后，go已经无法找到原文件来获取信息了,所以简单粗爆直接remove
					fmt.Println("重命名文件 : ", ev.Name)
					notifyfile.watch.Remove(ev.Name)

					//go notifyfile.PushEventChannel(ev.Name, fsnotify.Rename, "重命名文件")

				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					fmt.Println("修改权限 : ", ev.Name)

					//go notifyfile.PushEventChannel(ev.Name, fsnotify.Chmod, "修改权限")

				}
			}
		case err := <-notifyfile.watch.Errors:
			{
				fmt.Println("error : ", err)
				return
			}
		}
	}
}

//PushEventChannel 将发生事件加入channel
func (notifyfile *NotifyFile) PushEventChannel(Path string, ActionType fsnotify.Op, desc string, source string, target string) {
	notifyfile.Path <- ActionPath{Path: Path, ActionType: ActionType, desc: desc, SoucePath: source, TargetPath: target}

}
