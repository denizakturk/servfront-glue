package main

import (
	"fmt"
	"path/filepath"
)

//////go:generate go run servfront-glue/generator ./example/hello.go.tmpl ./helper.go ./example/param.json
////go:generate go run ./generator ./template/servfront/controller.go.tmpl ./index.go ./example/param.json
func main() {
	fmt.Println(filepath.Abs("./"))
}
