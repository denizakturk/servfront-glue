package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/denizakturk/servfront-glue/autogen/model"
	"github.com/denizakturk/servfront-glue/codegen"
	"github.com/denizakturk/servfront-glue/template/backend/go/servfront"
	"github.com/denizakturk/servfront-glue/template/client/node"
)

const SCHEMA_FILENAME = "schame.json"
const SCHEMA_FILE_PATH = "./" + SCHEMA_FILENAME

func main() {

	whereFromRunning, _ := filepath.Abs("./")
	fmt.Println("Running path from " + whereFromRunning)

	if _, schemaFileErr := os.Stat(SCHEMA_FILE_PATH); errors.Is(schemaFileErr, os.ErrNotExist) {
		panic(SCHEMA_FILENAME + " file not found, please before run github.com/denizakturk/servfront-glue/autogen/init")
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

			codegen.GenerateFromTemplate(servfront.MODEL_TEMPLATE, "./model/"+modelFileName+".go", modelParam)
		}
	}

	// Controller Middleware
	os.Mkdir("./middleware", os.ModePerm)
	if _, err := os.Stat("./middleware/middleware.go"); errors.Is(err, os.ErrNotExist) {
		// Generate Middleware
		codegen.GenerateFromTemplate(servfront.MIDDLEWARE_TEMPLATE, "./middleware/middleware.go", nil)
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
		codegen.GenerateFromTemplate(servfront.CONTROLLER_TEMPLATE, "./controller/"+controllerFileName+".go", controllerParams)
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
							codegen.GenerateFromTemplate(node.REQUEST_CLASS_TEMPLATE, "./frontend/"+frontendModelFile+".ts", modelParam)
						} else {
							codegen.GenerateFromTemplate(node.RESPONSE_INTERFACE_TEMPLATE, "./frontend/"+frontendModelFile+".ts", modelParam)
						}
					}
				}
			}
		}
	}

}
