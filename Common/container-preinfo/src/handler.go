package main

import (
	"context"
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
	// Collection name for container pre-info
	col_coninfo = "containerPreInfo"

	// Collection name for trigger & code
	col_trigger = "workingsetInfo"
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
	// Unmarshaling user's json input
	// then, insert data to database
	/////////////////////////////////////////////////////////
	var reg_num string
	data := make(map[string]interface{})
	json.Unmarshal(req.Input, &data)

	reg_num, err = mongo.InsertContainerInfo(conn, db_info.Authentications.MongoDB.Database, col_coninfo, data)
	if err != nil {
		return err.Error()
	}

	/////////////////////////////////////////////////////////
	// After insertion, send obj file to storage
	/////////////////////////////////////////////////////////
	for _, file := range data["obj_data"].([]interface{}) {
		if file.(map[string]interface{})["base64_file"] != nil && file.(map[string]interface{})["obj_name"] != nil {
			serviceName := "obj-storage"
			path := mongo.GetFilePath(conn, db_info.Authentications.MongoDB.Database, col_coninfo, reg_num, file.(map[string]interface{})["obj_name"].(string))

			input := make(map[string]interface{})
			input["base64_file"] = file.(map[string]interface{})["base64_file"].(string)
			input["obj_name"] = file.(map[string]interface{})["obj_name"].(string)
			input["path"] = path

			marshaled, _ := json.Marshal(input)

			// call trigger
			res := mesh.MeshCall(serviceName, marshaled)
			log.Printf(res)
		}
	}

	/////////////////////////////////////////////////////////
	// then, verify code for calling trigger
	/////////////////////////////////////////////////////////
	for _, value := range data["related_works"].([]interface{}) {
		if value.(map[string]interface{})["code"] != nil && value.(map[string]interface{})["code"].(string) != "00000" {
			// Firstly, get trigger name from table by code
			triggerName := mongo.GetTriggerName(conn, db_info.Authentications.MongoDB.Database, col_trigger, value.(map[string]interface{})["code"].(string))

			// Secondly, create []byte input with reg num and container num for calling trigger
			input := make(map[string]interface{})
			input["reg_num"] = reg_num

			meta := data["meta_data"].(map[string]interface{})
			container_num := meta["container_num"].(string)
			input["container_num"] = container_num

			marshaled, _ := json.Marshal(input)

			// Finally, call trigger
			res := mesh.MeshCall(triggerName, marshaled)
			log.Printf(res)
		}
	}

	return "success\n"
}
