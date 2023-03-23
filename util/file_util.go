package util

import (
	"os"
	"path/filepath"
)

// CreateFile函数在指定的目录下创建一个指定名称的文件
func CreateFile(path string, filename string) {

	// 创建目录及其父目录
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		PrintLog(err.Error())
		return
	}

	// 创建文件并关闭
	file, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		PrintLog(err.Error())
		return
	}
	defer file.Close()
}

// ReadStringFromFile函数从指定的文件中读取字符串并返回该字符串
func ReadStringFromFile(path string, filename string) string {

	// 拼接文件路径
	filePath := filepath.Join(path, filename)

	// 打开文件并读取内容
	file, err := os.Open(filePath)
	if err != nil {
		PrintLog(err.Error())
		return ""
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		PrintLog(err.Error())
		return ""
	}

	size := stat.Size()
	data := make([]byte, size)

	_, err = file.Read(data)
	if err != nil {
		PrintLog(err.Error())
		return ""
	}

	// 返回读取的字符串
	return string(data)
}

// WriteStringToFile函数将指定的字符串写入指定目录下的指定文件中
func WriteStringToFile(directory string, filename string, content string) {

	// 创建目录及其父目录
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		PrintLog(err.Error())
		return
	}

	// 拼接文件路径，将字符串写入文件
	filePath := filepath.Join(directory, filename)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		PrintLog(err.Error())
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		PrintLog(err.Error())
		return
	}

}
