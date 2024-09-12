package main

import (
	"fmt"
	"log"
	"net/http"
	urlshort "urlshort/handler"
)

func main() {
	mux := defaultMux()
	yamlFilePath := "yaml-data.yaml"

	yamlHandler, err := urlshort.YAMLHandler(yamlFilePath, mux)
	if err != nil {
		log.Fatalf("failed to process request:\n\t%v", err)
	}

	log.Fatal(http.ListenAndServe(":8080", yamlHandler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
