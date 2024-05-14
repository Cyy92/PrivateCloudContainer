package main

import (
	sdk "github.com/keti-openfx/openfx/executor/go/pb"
	"time"
)

//import mesh "github.com/keti-openfx/openfx/executor/go/mesh"

func Handler(req sdk.Request) string {
	// mesh call
	//
	// functionName := "<FUNCTIONNAME>"
	// input := string(req.Input)
	// result := mesh.MeshCall(functionName, []byte(input))
	// return result

	// single call
	time.Sleep(105 * time.Millisecond)
	return "[Go] " + string(req.Input)
}
