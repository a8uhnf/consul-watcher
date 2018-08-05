package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
)

func main() {
	fmt.Println("Hello World!!!")
	url := os.Getenv("CONSUL_URL")
	path := os.Getenv("CONSUL_PATH")

	log.Printf("Reading config file %s from %s", path, url)


	rq, err := http.NewRequest("GET", fmt.Sprintf("http://%s/v1/kv/%s", url, path), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(rq)
}
