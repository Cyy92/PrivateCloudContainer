package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"

	"github.com/keti-openfx/openfx/executor/go/mongoDB/model"
	mongo "github.com/keti-openfx/openfx/executor/go/mongoDB/src"
	sdk "github.com/keti-openfx/openfx/executor/go/pb"
)

const (
	// Collection name for workingset info
	col_wsinfo = "workingsetInfo"
)

func Handler(req sdk.Request) string {
	// Setting user credential for accessing MongoDB
	/////////////////////////////////////////////////////////
	var db_info model.Credential

	marshal_cred, err := ioutil.ReadFile("handler/cred.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal([]byte(marshal_cred), &db_info)
	if err != nil {
		log.Fatalln(err)
	}

	/////////////////////////////////////////////////////////
	// Create client for MongoDB
	/////////////////////////////////////////////////////////
	conn := mongo.MongoConn(db_info.Authentications.MongoDB.URI, db_info.Authentications.MongoDB.UserName, db_info.Authentications.MongoDB.PassWord, db_info.Authentications.MongoDB.Database)
	defer conn.Disconnect(context.TODO())

	/////////////////////////////////////////////////////////
	// Get workingset info
	/////////////////////////////////////////////////////////
	fmt.Println(string(req.Input) + " : This input has no effect on the result")
	ws := mongo.GetWorkingsetInfo(conn, db_info.Authentications.MongoDB.Database, col_wsinfo)

	return ws + "\n"
}
