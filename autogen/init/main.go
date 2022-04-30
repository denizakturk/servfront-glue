package main

import (
	"encoding/json"
	"github.com/denizakturk/servfront-glue/autogen/model"
	"log"
	"os"
)

func main() {
	var services map[string]model.Service
	services = make(map[string]model.Service)
	services["ExampleService"] = model.Service{
		PackageName:    "PackageNameHere",
		ControllerName: "ControllerNameHere",
		Model: []model.Model{
			{
				Name: "ControllerRequestModel",
				Field: []model.ModelField{
					{
						Name: "ExampleModelFieldName",
						Type: "string",
					},
				},
			},
			{
				Name: "ControllerResponseModel",
				Field: []model.ModelField{
					{
						Name: "ExampleModelFieldName",
						Type: "string",
					},
				},
			},
		},
		Method: []model.ControllerMetod{
			{
				Name:              "ExampleController",
				RequestModelName:  "ControllerRequestModel",
				ResponseModelName: "ControllerResponseModel",
			},
		},
	}
	autogenInit := services
	jsonMarshal, jmErr := json.MarshalIndent(autogenInit, " ", "\t")
	if nil != jmErr {
		log.Println(jmErr)
	}
	os.WriteFile("./schema.json", jsonMarshal, os.ModePerm)
}
