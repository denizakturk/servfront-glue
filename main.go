package main

//go:generate ./generator template/hello.go.tmpl schema.json
func main() {
	HelloFunc()
}
