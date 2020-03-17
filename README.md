### 前言
1.本人是个重度某度云用户，但是云盘的速度不开会员的情况下真是太慢了，而且也不能实时备份，所以想过自己搭建一个NAS。写个程序能可以实时把我的工作目录备份到NAS上。

2.工作中，我们开发组想搭建一个公共的开发机，但是samba协议挂在到本地以后无法用svn 或git 与版本库同步，最近突然想到其实也可以用实时同步的程序，把本地的脚本代码实时同步到开发机上，这个就可以解决开发机的问题了


### 开发环境
linux deepin 系统
golang 版本  1.14


### 实现的功能
1.可配置多个实时同步目录。
只需要再conf.json中添加好源目录和目标目录即可
```json
"Sync": [
        {
            "Name":"pc项目",
            "Sourcepath": "/home/zlc/goProject/src/ftpsync/test",
            "Targetpath": "/home/zlc/project/sftpdata/upload2"
        },
        {
            "Name":"触屏项目",
            "Sourcepath": "/home/zlc/goProject/src/ftpsync/test2",
            "Targetpath": "/home/zlc/project/sftpdata/upload"
        }
        ]
```
2.实现了监控文件目录的实时增加和删除

3.通过sftp协议同步文件到服务器

4.所有的更新操作都记录到日志文件

### 软件规划
1.以上基础上添加一键初始化所有目录到服务器
2.添加gui可视化界面
3.完善单元测试
4.交叉编译，支持主流平台使用Windows linux mac  等等

### 最后
再完成这个小软件的过程中，我发现早已经有同行有过相同的想法，但是都开源代码都不是很完善，今后我会再业余时间了尽量完善这个软件，也是对我自学golang以来学习成果的一个验证