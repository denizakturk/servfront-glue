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

const SCHEMA_FILENAME = "schema.json"
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
	var services model.Services
	json.Unmarshal(jsonString, &services)
	for _, val := range services {
		parameters := val
		// Models
		os.Mkdir(parameters.ModelPath, os.ModePerm)

		for _, val := range parameters.Model {
			modelFileName := codegen.To_snake_case(val.Name)
			if _, errModelFile := os.Stat(parameters.ModelPath + "/" + modelFileName + ".go"); errors.Is(errModelFile, os.ErrNotExist) {
				modelParam := servfront.ModelTmplParams{Name: val.Name}
				var fields []servfront.ModelTmplParamFields
				for _, fieldVal := range val.Field {
					f := servfront.ModelTmplParamFields{Name: fieldVal.Name, Type: fieldVal.Type}
					fields = append(fields, f)
				}

				modelParam.Fields = fields

				codegen.GenerateFromStringTemplate(servfront.MODEL_TEMPLATE, parameters.ModelPath+"/"+modelFileName+".go", modelParam)
			}
		}

		// Controller Middleware
		os.Mkdir(parameters.MiddlewarePath, os.ModePerm)
		if _, err := os.Stat(parameters.MiddlewarePath + "/middleware.go"); errors.Is(err, os.ErrNotExist) {
			// Generate Middleware
			codegen.GenerateFromStringTemplate(servfront.MIDDLEWARE_TEMPLATE, parameters.MiddlewarePath+"/middleware.go", nil)
		}

		// Controller
		os.Mkdir(parameters.ControllerPath, os.ModePerm)
		controllerFileName := codegen.To_snake_case(parameters.ControllerName)
		if _, err := os.Stat(parameters.ControllerPath + "/" + controllerFileName + ".go"); errors.Is(err, os.ErrNotExist) {
			controllerParams := servfront.ControllerTmplParams{ControllerName: parameters.ControllerName}
			for _, val := range parameters.Method {
				method := servfront.ControllerMethod{Name: val.Name, RequestModelName: val.RequestModelName, ResponseModelName: val.ResponseModelName}
				controllerParams.Method = append(controllerParams.Method, method)
			}
			// Generate Controller
			codegen.GenerateFromStringTemplate(servfront.CONTROLLER_TEMPLATE, parameters.ControllerPath+"/"+controllerFileName+".go", controllerParams)
		}

		// Frontend Models
		os.Mkdir(parameters.FrontendPath, os.ModePerm)
		for _, methodVal := range parameters.Method {
			if methodVal.RequestModelName != "" {
				for _, val := range parameters.Model {
					if methodVal.RequestModelName == val.Name || methodVal.ResponseModelName == val.Name {
						frontendModelFile := codegen.To_snake_case(val.Name)
						if _, errModelFile := os.Stat(parameters.FrontendPath + "/" + frontendModelFile + ".ts"); errors.Is(errModelFile, os.ErrNotExist) {
							modelParam := servfront.ModelTmplParams{Name: val.Name}
							var fields []servfront.ModelTmplParamFields
							for _, fieldVal := range val.Field {
								f := servfront.ModelTmplParamFields{Name: fieldVal.Name, Type: fieldVal.Type}
								fields = append(fields, f)
							}
							modelParam.Fields = fields
							if methodVal.RequestModelName == val.Name {
								codegen.GenerateFromStringTemplate(node.REQUEST_CLASS_TEMPLATE, parameters.FrontendPath+"/"+frontendModelFile+".ts", modelParam)
							} else {
								codegen.GenerateFromStringTemplate(node.RESPONSE_INTERFACE_TEMPLATE, parameters.FrontendPath+"/"+frontendModelFile+".ts", modelParam)
							}
						}
					}
				}
			}
		}
	}
}
