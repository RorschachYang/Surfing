package util

import (
	"fmt"
	"log"
	"os"
	"time"
)

func PrintLog(logText string) {

	// 获取配置文件中的日志路径
	logDir := GetConfigString("Log.logDir")

	// 获取当前日期并格式化为"2006-01-02"格式
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02")

	// 拼接目标文件路径
	filePath := logDir + dateString + ".log"

	// 判断文件是否存在，如果不存在则创建
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	// 设置日志输出文件
	logFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// 输出日志
	log.Println(logText)
	fmt.Println(logText)

}
