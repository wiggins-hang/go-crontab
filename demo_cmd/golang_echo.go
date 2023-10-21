package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {

	filePath := "test.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	writeTime := time.Now().Format("2006-01-02 - 15:04:05")
	write.WriteString(fmt.Sprintf("golang_test %s \n", writeTime))
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}
