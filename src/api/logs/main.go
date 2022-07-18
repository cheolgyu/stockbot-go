package main

import (
	"context"
	"encoding/json"
	"fmt"
	syslog "log"
	"net/http"
	"time"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/mlog"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var logs_collection *mongo.Collection

func init() {
	client, _ = common.Connect()
	logs_collection = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_LOGS)
}
func main() {
	http.HandleFunc("/logging", logging)
	syslog.Println(http.ListenAndServe(":5000", nil))
}

type LogInfo struct {
	mlog.LOG `bson:"inline"`
	RecvInfo `bson:"inline"`
}
type RecvInfo struct {
	Time       string      `bson:"recv_time" `
	RemoteAddr interface{} `bson:"recv_remote_addr" `
}

func logging(w http.ResponseWriter, req *http.Request) {
	log := LogInfo{}
	log.RecvInfo.Time = time.Now().String()
	log.RecvInfo.RemoteAddr = req.RemoteAddr

	err := json.NewDecoder(req.Body).Decode(&log.LOG)
	if err != nil {
		fmt.Printf("json error : %+v \n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

	}

	fmt.Fprintf(w, "%+v \n", log)

	insert(log)
}

func insert(loginfo LogInfo) {
	_, err := logs_collection.InsertOne(context.TODO(), loginfo)
	if err != nil {
		syslog.Fatalln(err)

	}

}
