package main

import (
	"errors"
	"os"
	"servfront-glue/codegen"
)

func main() {
	os.Mkdir("./middleware", os.ModePerm)
	if _, err := os.Stat("./middleware/middleware.go"); errors.Is(err, os.ErrNotExist) {

		// Generate Middleware
		codegen.GenerateFromTemplate("./../template/backend/go/servfront/middleware.go.tmpl", "./middleware/middleware.go", "./../example/param.json")
	}
	os.Mkdir("./../controller", os.ModePerm)
	if _, err := os.Stat("./controller/index.go"); errors.Is(err, os.ErrNotExist) {

		// Generate Controller
		codegen.GenerateFromTemplate("./../template/backend/go/servfront/controller.go.tmpl", "./controller/index.go", "./../example/param.json")
	}
}
