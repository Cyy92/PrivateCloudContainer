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
	// Collection name for working set
	col_ws = "dangerInfo"
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
	// then, update data in database and send obj file to storage
	/////////////////////////////////////////////////////////
	data := make(map[string]interface{})
	json.Unmarshal(req.Input, &data)

	err = mongo.UpdateWorkingSetInfo(conn, db_info.Authentications.MongoDB.Database, col_ws, data["reg_num"].(string), data["working_reg_num"].(string), data)
	if err != nil {
		return err.Error()
	}

	for _, file := range data["obj_data"].([]interface{}) {
		if file.(map[string]interface{})["base64_file"] != nil && file.(map[string]interface{})["obj_name"] != nil {
			serviceName := "obj-storage"
			path := mongo.GetFilePath(conn, db_info.Authentications.MongoDB.Database, col_ws, data["reg_num"].(string), file.(map[string]interface{})["obj_name"].(string))

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
	// Setting for service mesh
	// set trigger name & create []byte input with reg num and working set reg num for calling trigger
	/////////////////////////////////////////////////////////
	triggerName := "meta-updater"
	input := make(map[string]interface{})
	input["reg_num"] = data["reg_num"].(string)
	input["working_reg_num"] = data["working_reg_num"].(string)
	input["code"] = "DA001"

	marshaled, _ := json.Marshal(input)

	// call trigger
	res := mesh.MeshCall(triggerName, marshaled)
	log.Printf(res)

	return "success\n"
}
