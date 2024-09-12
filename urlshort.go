package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	urlshort "urlshort/handler"
)

func main() {
	dataFilePath := flag.String("f", "data.yaml", "set custom yaml or json file")
	flag.Parse()

	mux := defaultMux()

	handler, err := urlshort.MakeHandler(*dataFilePath, mux)
	if err != nil {
		log.Fatalf("failed to process request:\n\t%v", err)
	}

	// run server
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
