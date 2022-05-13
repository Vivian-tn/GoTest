package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "./test.csv"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("file does not exist")
	} else {
		//打开文件(只读模式)，创建io.read接口实例
		csvFile, err := os.Open(filename)
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}
		defer csvFile.Close()

		// 解析文件 创建csv读取接口实例
		p := csv.NewReader(csvFile)

		/*
		  说明：
		   1、读取csv文件返回的内容为切片类型，可以通过遍历的方式使用或Slicer[0]方式获取具体的值。
		   2、同一个函数或线程内，两次调用Read()方法时，第二次调用时得到的值为每二行数据，依此类推。
		   3、大文件时使用逐行读取，小文件直接读取所有然后遍历，两者应用场景不一样，需要注意。
		*/
		//获取一行内容，一般为第一行内容
		read, _ := p.Read() //返回切片类型：[chen  hai wei]
		fmt.Println(read)

		//读取所有内容
		ReadAll, err := p.ReadAll() //返回切片类型：[[s s ds] [a a a]]
		fmt.Println(ReadAll)
	}

}
