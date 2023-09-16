package main

import (
	"net/http"
)

func route() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", report)
}

func report(w http.ResponseWriter, r *http.Request) {

}
