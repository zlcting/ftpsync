package sftphandler

import (
	"ftpsync/utils"
	"log"

	"github.com/pkg/sftp"
)

var (
	err        error
	sftpClient *sftp.Client
)

func fileSftpHandler() {
	// Connect
	sftpClient, err = connect(utils.GlobalObject.Sftp.Name, utils.GlobalObject.Sftp.Pass, utils.GlobalObject.Sftp.Host, utils.GlobalObject.Sftp.Port)

	if err != nil {
		log.Fatal("SSH connect error: ", err)
	}

	defer sftpClient.Close()

	// Check if the remote dir exist
	if _, err := sftpClient.Stat(remoteDir); err != nil {
		panic("Remote dir dose not exist: " + remoteDir)
	}

}
