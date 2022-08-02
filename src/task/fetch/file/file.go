package file

import (
	"log"
	"os"
)

func Mkdir(list []string) {
	for _, item := range list {
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
// func (o *File) CreateFile(fileName string) *os.File {
// 	file, err := os.OpenFile(fileName, FILE_FLAG_TRUNC, 0644)
// 	o.CheckError(err)
// 	return file
// }

func (o *File) CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// 기존 파일을 새로 만들어버림.
func CreateFile(fileName string) *os.File {
	file, err := os.Create(fileName)
	checkError(err)
	return file
}
