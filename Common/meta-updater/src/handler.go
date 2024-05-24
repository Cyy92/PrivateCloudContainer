package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"

	"github.com/keti-openfx/openfx/executor/go/mongoDB/model"
	mongo "github.com/keti-openfx/openfx/executor/go/mongoDB/src"
	sdk "github.com/keti-openfx/openfx/executor/go/pb"
)

const (
	// Collection name for metadata manager
	col_mdmanger = "mdManager"
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
	var wsReg string

	/////////////////////////////////////////////////////////
	// Update Documentaion with key from trigger input
	/////////////////////////////////////////////////////////
	if data["working_reg_num"] == nil {
		wsReg = "none"
		err = mongo.UpdateMetaByTrigger(conn, db_info.Authentications.MongoDB.Database, col_mdmanger, data["reg_num"].(string), wsReg, data["code"].(string))
		if err != nil {
			return err.Error()
		}

		return "Successfully update database\n"
	} else {
		wsReg = data["working_reg_num"].(string)
		err = mongo.UpdateMetaByWorkingSet(conn, db_info.Authentications.MongoDB.Database, col_mdmanger, data["reg_num"].(string), wsReg, data["code"].(string))
		if err != nil {
			return err.Error()
		}

		return "Successfully update database\n"
	}
}
