package main

import (
	"net/http"

	logger "github.com/sirupsen/logrus"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
	r.Write(w)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	logger.Info("Hello, world!")
}