package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"

	mesh "github.com/keti-openfx/openfx/executor/go/mesh"
	"github.com/keti-openfx/openfx/executor/go/mongoDB/model"
	mongo "github.com/keti-openfx/openfx/executor/go/mongoDB/src"
	sdk "github.com/keti-openfx/openfx/executor/go/pb"
)

const (
	// Collection name for workingset
	col_ws = "dangerInfo"

	// Code name for workingset
	ws_code = "DA001"
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
	//defer conn.Disconnect(context.TODO())

	/////////////////////////////////////////////////////////
	// Unmarshaling user's json input
	// then, insert data to database
	/////////////////////////////////////////////////////////
	data := make(map[string]interface{})
	json.Unmarshal(req.Input, &data)

	// Insert data to dangerInfo collection
	err = mongo.InsertWorkingSetInfo(conn, db_info.Authentications.MongoDB.Database, col_ws, data)
	if err != nil {
		return err.Error()
	}

	/////////////////////////////////////////////////////////
	// Setting for service mesh
	// set trigger name & create []byte input with reg num and code for calling trigger
	/////////////////////////////////////////////////////////
	triggerName := "meta-updater"
	input := make(map[string]interface{})
	input["reg_num"] = data["reg_num"].(string)
	input["code"] = ws_code

	marshaled, _ := json.Marshal(input)

	// call trigger
	res := mesh.MeshCall(triggerName, marshaled)
	log.Printf(res)

	return "success\n"
}
