package main

import (
	"log"
	"net/http"
	"wxcloudrun-golang/service"
)

func main() {

	http.HandleFunc("/", service.HelloWorldHandler)
	http.HandleFunc("/tags", service.TagsHandler)
	http.HandleFunc("/send", service.SendHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}
