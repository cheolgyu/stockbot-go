package common

import (
	"log"

	"github.com/joho/godotenv"
)

var MyEnv map[string]string

func init() {
	var err error
	MyEnv, err = godotenv.Read(".env.local")
	if err != nil {
		log.Panic("env 파일읽기 실패:", err)
	}
	log.Println("env init 실행", MyEnv)
}
