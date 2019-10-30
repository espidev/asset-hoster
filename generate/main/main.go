package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
)

func main() {
	assets := http.Dir("assets")
	err := vfsgen.Generate(assets, vfsgen.Options{})
	if err != nil {
		log.Fatalln(err)
	}
}