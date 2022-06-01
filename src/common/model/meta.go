package model

type Code struct {
	Code string `bson:"code"`
	Name string `bson:"name"`
}

type Config struct {
	Id         int
	Upper_code string
	Upper_name string
	Code       string
	Name       string
}

type DownloadInfo struct {
	Code    Code
	StartDt string
	EndDt   string
}

type Opening struct {
	YY      int
	MM      int
	DD      int
	Week    int
	Quarter int
}
