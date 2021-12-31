package main

import (
	"encoding/json"
	"log"
	"os"
	"servfront-glue/autogen/model"
)

func main() {
	autogenInit := &model.AutogenInit{
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
	jsonMarshal, jmErr := json.MarshalIndent(autogenInit, " ", "\t")
	if nil != jmErr {
		log.Println(jmErr)
	}
	os.WriteFile("./schame.json", jsonMarshal, os.ModePerm)
}
