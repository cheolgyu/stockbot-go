package kr

import (
	"log"
	"os"

	"github.com/cheolgyu/stockbot/src/fetch/kr/config"
)

func init() {
	log.Println("i am fetch file init")
	mkdir := []string{
		config.DOWNLOAD_DIR_COMPANY_DETAIL,
		config.DOWNLOAD_DIR_COMPANY_STATE,
		config.DOWNLOAD_DIR_PRICE,
		config.DOWNLOAD_DIR_MARKET,
	}

	for _, item := range mkdir {
		err := os.MkdirAll(item, 0755)
		if err != nil {
			panic(err)
		}
	}
}

type File struct{}

func (o *File) Open(fileName string) *os.File {
	file, err := os.Open(fileName)

	o.CheckError(err)
	return file
}

func (o *File) Write(f *os.File, text string) {
	_, err := f.WriteString(text + "\n")

	o.CheckError(err)
}

// 새로쓰기
func (o *File) CreateFile(fileName string) *os.File {
	file, err := os.OpenFile(fileName, config.FILE_FLAG_TRUNC, 0644)
	o.CheckError(err)
	return file
}

func (o *File) CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
