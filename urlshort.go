package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	urlshort "urlshort/handler"
)

func main() {
	yamlFilePath := flag.String("yaml", "yaml-data.yaml", "set custom yaml file")
	flag.Parse()

	mux := defaultMux()

	yamlHandler, err := urlshort.YAMLHandler(*yamlFilePath, mux)
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
