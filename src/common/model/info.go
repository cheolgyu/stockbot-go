package model

type Info struct {
	Name    string `bson:"_id"`
	Updated string `bson:"updated"`
}
