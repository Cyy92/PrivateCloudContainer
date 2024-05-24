package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/keti-openfx/openfx/executor/go/minio/model"
	minio "github.com/keti-openfx/openfx/executor/go/minio/src"
	sdk "github.com/keti-openfx/openfx/executor/go/pb"
)

const (
	bucketName = "test0810"
)

func Handler(req sdk.Request) string {
	// Setting user credential for accessing MinIO
	/////////////////////////////////////////////////////////
	var storage model.Credential

	marshal_cred, err := ioutil.ReadFile("handler/cred.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal([]byte(marshal_cred), &storage)
	if err != nil {
		log.Fatalln(err)
	}

	/////////////////////////////////////////////////////////
	// Unmarshaling user's json input
	// then, decode user input
	/////////////////////////////////////////////////////////
	data := make(map[string]interface{})
	json.Unmarshal(req.Input, &data)

	dec, dec_err := base64.StdEncoding.DecodeString(data["base64_file"].(string))
	if dec_err != nil {
		panic(dec_err)
	}

	if _, stat_err := os.Stat(data["path"].(string)); os.IsNotExist(stat_err) {
		trim_path := strings.Trim(data["path"].(string), data["obj_name"].(string))
		os.MkdirAll(trim_path, 0700)

		f, f_err := os.Create(data["path"].(string))
		if f_err != nil {
			panic(f_err)
		}
		defer f.Close()

		if _, w_err := f.Write(dec); w_err != nil {
			panic(w_err)
		}
		if s_err := f.Sync(); s_err != nil {
			panic(s_err)
		}
	}
	/////////////////////////////////////////////////////////
	// template setting
	// please change these info before send obj file
	/////////////////////////////////////////////////////////
	contentType := minio.GetContentType(data["obj_name"].(string))
	useSSL := false

	/////////////////////////////////////////////////////////
	// send obj file
	/////////////////////////////////////////////////////////
	put_err := minio.PutObject(storage.Authentications.MinIO.EndPoint, storage.Authentications.MinIO.AccessKey, storage.Authentications.MinIO.SecretKey, useSSL, data["obj_name"].(string), bucketName, data["path"].(string), contentType)
	if put_err != nil {
		return put_err.Error()
	}

	return "success\n"
}
