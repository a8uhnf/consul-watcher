package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/hashicorp/consul/api"
"fmt"
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



func consul_test_api() {
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	// PUT a new KV pair
	p := &api.KVPair{Key: "REDIS_MAXCLIENTS", Value: []byte("1000")}
	_, err = kv.Put(p, nil)
	if err != nil {
		panic(err)
	}

	// Lookup the pair
	pair, _, err := kv.Get("REDIS_MAXCLIENTS", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)
}
