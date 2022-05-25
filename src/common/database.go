package common

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, context.Context) {
	ctx := context.Background()
	credential := options.Credential{
		Username: MyEnv["MONGO_DB_USERNAME"],
		Password: MyEnv["MONGO_DB_PASSWORD"],
	}
	// 몽고DB 연결
	clientOptions := options.Client().ApplyURI(MyEnv["MONGO_DB_URL"]).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("몽고 DB에 연결했습니다!")

	// 내용을 적을 부분

	// 몽고DB 연결 끊기
	uesrsCollection := client.Database("test").Collection("users")
	fmt.Println(uesrsCollection)

	//defer client.Disconnect(ctx)

	// err := defer  client.Disconnect(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("몽고DB 연결을 종료했습니다!")

	return client, ctx
}
