package mylog

import (
	"ftpsync/utils"
	"log"
	"os"
)

var Logger *log.Logger

func init() {

	file, err := os.Create(utils.GlobalObject.Logpath)

	if err != nil {
		log.Fatalln(err)
		log.Fatalln("fail to create test.log file!")
	}

	Logger = log.New(file, "", log.Llongfile)
	Logger.SetFlags(log.LstdFlags)
}
