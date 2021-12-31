package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"servfront-glue/autogen/model"
	"servfront-glue/codegen"
	"servfront-glue/template/backend/go/servfront"
)

const SCHEMA_FILENAME = "schame.json"
const SCHEMA_FILE_PATH = "./" + SCHEMA_FILENAME

func main() {

	whereFromRunning, _ := filepath.Abs("./")
	fmt.Println("Running path from " + whereFromRunning)

	if _, schemaFileErr := os.Stat(SCHEMA_FILE_PATH); errors.Is(schemaFileErr, os.ErrNotExist) {
		panic(SCHEMA_FILENAME + " file not found, please before run github.com/servfront-glue/autogen/init")
	}

	schemaFile, err := os.Open(SCHEMA_FILE_PATH)
	if nil != err {
		panic(err)
	}
	defer schemaFile.Close()

	jsonString, _ := ioutil.ReadAll(schemaFile)
	var parameters model.AutogenInit

	json.Unmarshal(jsonString, &parameters)
	// Models
	os.Mkdir("./model", os.ModePerm)

	for _, val := range parameters.Model {
		modelFileName := codegen.To_snake_case(val.Name)
		if _, errModelFile := os.Stat("./model/" + modelFileName + ".go"); errors.Is(errModelFile, os.ErrNotExist) {
			modelParam := servfront.ModelTmplParams{Name: val.Name}
			var fields []servfront.ModelTmplParamFields
			for _, fieldVal := range val.Field {
				f := servfront.ModelTmplParamFields{Name: fieldVal.Name, Type: fieldVal.Type}
				fields = append(fields, f)
			}

			modelParam.Fields = fields
			codegen.GenerateFromTemplate("./../template/backend/go/servfront/model.go.tmpl", "./model/"+modelFileName+".go", modelParam)
		}
	}

	// Controller Middleware
	os.Mkdir("./middleware", os.ModePerm)
	if _, err := os.Stat("./middleware/middleware.go"); errors.Is(err, os.ErrNotExist) {
		// Generate Middleware
		codegen.GenerateFromTemplate("./../template/backend/go/servfront/middleware.go.tmpl", "./middleware/middleware.go", nil)
	}

	// Controller
	os.Mkdir("./controller", os.ModePerm)
	controllerFileName := codegen.To_snake_case(parameters.ControllerName)
	if _, err := os.Stat("./controller/" + controllerFileName + ".go"); errors.Is(err, os.ErrNotExist) {
		controllerParams := servfront.ControllerTmplParams{ControllerName: parameters.ControllerName}
		for _, val := range parameters.Method {
			method := servfront.ControllerMethod{Name: val.Name, RequestModelName: val.RequestModelName, ResponseModelName: val.ResponseModelName}
			controllerParams.Method = append(controllerParams.Method, method)
		}
		// Generate Controller
		codegen.GenerateFromTemplate("./../template/backend/go/servfront/controller.go.tmpl", "./controller/"+controllerFileName+".go", controllerParams)
	}

	// Frontend Models
	os.Mkdir("./frontend", os.ModePerm)
	for _, methodVal := range parameters.Method {
		if methodVal.RequestModelName != "" {
			for _, val := range parameters.Model {
				if methodVal.RequestModelName == val.Name || methodVal.ResponseModelName == val.Name {
					frontendModelFile := codegen.To_snake_case(val.Name)
					if _, errModelFile := os.Stat("./frontend/" + frontendModelFile + ".ts"); errors.Is(errModelFile, os.ErrNotExist) {
						modelParam := servfront.ModelTmplParams{Name: val.Name}
						var fields []servfront.ModelTmplParamFields
						for _, fieldVal := range val.Field {
							f := servfront.ModelTmplParamFields{Name: fieldVal.Name, Type: fieldVal.Type}
							fields = append(fields, f)
						}
						modelParam.Fields = fields
						if methodVal.RequestModelName == val.Name {
							codegen.GenerateFromTemplate("./../template/client/node/request_class.ts.tmpl", "./frontend/"+frontendModelFile+".ts", modelParam)
						} else {
							codegen.GenerateFromTemplate("./../template/client/node/response_interface.ts.tmpl", "./frontend/"+frontendModelFile+".ts", modelParam)
						}
					}
				}
			}
		}
	}

}
