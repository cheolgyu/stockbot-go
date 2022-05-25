package model

type Config struct {
	Id         int
	Upper_code string
	Upper_name string
	Code       string
	Name       string
}

type Code struct {
	Id        int
	Code      string
	Code_type int
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
