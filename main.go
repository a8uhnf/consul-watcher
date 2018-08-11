package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
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
	watchConsulConfig()
}

type HelloStruct struct {
	Good string `json:"good"`
}

func watchConsulConfig() {
	log.Println("------ Watch Config -----")
	w := &watch.Plan{
		Datacenter: "localhost:8500",
	}
	fmt.Println(w)

	var runtime_viper = viper.New()


	runtime_viper.AddRemoteProvider("consul", "localhost:8500", "hello")
	runtime_viper.SetConfigType("yaml")

	err := runtime_viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	s := &HelloStruct{}
	runtime_viper.Unmarshal(&s)


	stop := make(chan int)
	go func() {
		for {
			time.Sleep(time.Second * 5) // delay after each request

			// currently, only tested with etcd support
			err := runtime_viper.WatchRemoteConfig()
			if err != nil {
				log.Printf("unable to read remote config: %v", err)
				continue
			}
			fmt.Println(runtime_viper.Get("good"))
			log.Println("Changed the remote config")
		}
	}()

	<-stop
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
