package sftphandler

import (
	"fmt"
	"ftpsync/utils"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//SftpHandler 全局结构体
type SftpHandler struct {
	SftpClient *sftp.Client
}

//NewSftpHandler 初始化
func NewSftpHandler() *SftpHandler {
	w := new(SftpHandler)
	w.SftpClient, _ = connect(utils.GlobalObject.Sftp.Name, utils.GlobalObject.Sftp.Pass, utils.GlobalObject.Sftp.Host, utils.GlobalObject.Sftp.Port)
	return w
}

// 	if _, err := sftpClient.Stat(remoteDir); err != nil {
// 		panic("Remote dir dose not exist: " + remoteDir)
// 	}

//connect 生成链接 对象
func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

//首先上传文件的方法
func (sftpHandler *SftpHandler) uploadFile(localFilePath string, remotePath string) {
	//打开本地文件流
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		log.Fatal(err)
	}
	//关闭文件流
	defer srcFile.Close()
	//上传到远端服务器的文件名,与本地路径末尾相同
	var remoteFileName = path.Base(localFilePath)

	//判断当前目录是否存在
	if _, err := sftpHandler.SftpClient.Stat(remotePath); err != nil {
		sftpHandler.SftpClient.Mkdir(remotePath)
	}

	//打开远程文件,如果不存在就创建一个
	dstFile, err := sftpHandler.SftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)
	}
	//关闭远程文件
	defer dstFile.Close()
	//读取本地文件,写入到远程文件中(这里没有分快穿,自己写的话可以改一下,防止内存溢出)
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		log.Fatal(err)
	}
	dstFile.Write(ff)
	fmt.Println(localFilePath + "  copy file to remote server finished!")
}

//uploadDirectory 遍历上传远程文件夹
func (sftpHandler *SftpHandler) uploadDirectory(localPath string, remotePath string) {
	//打开本地文件夹流
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal("路径错误 ", err)
	}
	//先创建最外层文件夹
	sftpHandler.SftpClient.Mkdir(remotePath)
	//遍历文件夹内容
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		//判断是否是文件,是文件直接上传.是文件夹,先远程创建文件夹,再递归复制内部文件
		if backupDir.IsDir() {
			sftpHandler.SftpClient.Mkdir(remoteFilePath)
			sftpHandler.uploadDirectory(localFilePath, remoteFilePath)
		} else {
			sftpHandler.uploadFile(path.Join(localPath, backupDir.Name()), remotePath)
		}
	}

	fmt.Println(localPath + "  copy directory to remote server finished!" + remotePath + "目录L：")
}

//Upload 判断是否是路径属性
func (sftpHandler *SftpHandler) Upload(localPath string, remotePath string) {
	//获取路径的属性
	s, err := os.Stat(localPath)
	if err != nil {
		fmt.Println("文件路径不存在")
		return
	}

	//判断是否是文件夹
	if s.IsDir() {
		sftpHandler.uploadDirectory(localPath, remotePath)
	} else {
		sftpHandler.uploadFile(localPath, remotePath)
	}
}
