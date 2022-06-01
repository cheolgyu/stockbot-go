package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Code struct {
	Id   primitive.ObjectID `bson:"_id"`
	Code string             `bson:"code"`
	Name string             `bson:"name"`
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
