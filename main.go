package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)

func main() {
	fmt.Println("Hello World!!!")
	url := os.Getenv("CONSUL_URL")
	path := os.Getenv("CONSUL_PATH")

	log.Printf("Reading config file %s from %s", path, url)

	rq, err := http.Get(fmt.Sprintf("http://%s/v1/kv/%s", url, path))
	if err != nil {
		panic(err)
	}
	fmt.Println(rq)
	defer rq.Body.Close()
	b, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	// consul_test_api()
}

func watchConsulConfig() {
	log.Println("------ Watch Config -----")
	w := &watch.Plan{
		Datacenter: "localhost:8500",
	}
	fmt.Println(w)
}

func consul_test_api() {
	b, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	// PUT a new KV pair
	p := &api.KVPair{Key: "REDIS_MAXCLIENTS", Value: b}
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
