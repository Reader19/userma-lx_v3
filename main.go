package main

import (
	"log"
	"net/http"
	"userma-lx/router"
)

func main() {
	server := http.Server{
		Addr: "localhost:8080",
	}
	router.Router()
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}
