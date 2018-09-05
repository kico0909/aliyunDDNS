package mylog

import (
	"os"
	"log"
	"time"
		"io"
)

var f *os.File
var maxSize int64

var path string
var filename string



func SetLogMaxSize(size int64){
	maxSize = size
}

func SetLogPath(setPath, setFilename string){
	path = setPath
	filename = setFilename
}

func LogStart()bool{
	var err error
	fullPath := path+filename

	// 判断文件是否存在,不存在则创建,存在则判断是否超体积
	RunLogfileSize, err := os.Stat( fullPath + ".txt" )

	if err != nil {
		f, err = os.Create(fullPath+".txt")
		if err != nil {
			log.Println(err)
			return false
		}
		return true

	}else{

		if RunLogfileSize.Size() > maxSize {
			// 转存
			tmp_f, _ := os.Create(path + time.Now().Format("2006-01-02 15:04:05") + "_" + filename + ".txt" )
			tmp_f_o, _ := os.OpenFile(fullPath+".txt", os.O_RDONLY, 0664)
			defer tmp_f.Close()
			defer tmp_f_o.Close()
			io.Copy(tmp_f,tmp_f_o)

			// 重新打开
			f, err = os.OpenFile(fullPath+".txt", os.O_TRUNC|os.O_RDWR, 0777)
		} else {
			f, err = os.OpenFile(fullPath+".txt", os.O_RDWR|os.O_APPEND, 0777)

			if err != nil {
				log.Println(err)
				return false
			}
			return true
		}
	}

	return true
}

func Record( log string){
	f.WriteString(log)
}

func LogStop(){
	f.Close()
}
