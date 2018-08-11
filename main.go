package main

import (
	"bytes"
	"encoding/json"
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

type HelloStruct struct {
	Good string `json:"good"`
}

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
	stop := make(chan int)
	// consul_test_api()
	// watchConsulConfig()
	go runServer()
	watchExample()
	<-stop
}

func runServer() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello World!!!")
	})

	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}

func watchExample() {
	m := make(map[string]interface{})
	m["datacenter"] = "dc1"
	m["type"] = "key"
	m["key"] = "/hello"
	m["handler_type"] = "http"
	m["http_handler_config"] = makeParams(`{"path":"http://localhost/hello","method":"GET","header":{},"timeout":"10s","tls_skip_verify":true}`)

	plan, err := watch.Parse(m)
	if err != nil {
		panic(err)
	}
	plan.Watcher = func(pp *watch.Plan) (watch.BlockingParamVal, interface{}, error) {
		fmt.Println("Hello World!!!", pp)
		t := &HelloBlockingParamVal{}
		st := &HelloStruct{}
		return t, st, nil
	}
	err = plan.Run(fmt.Sprintf("http://%s%s", "localhost", ":8500"))
	if err != nil {
		panic(err)
	}
}

type HelloBlockingParamVal struct{}

func (hbpv *HelloBlockingParamVal) Equal(other watch.BlockingParamVal) bool {
	fmt.Println("Hello World equal!!!", other)
	panic("hello")
	return true
}
func (hbpv *HelloBlockingParamVal) Next(previous watch.BlockingParamVal) watch.BlockingParamVal {
	fmt.Println("Hello World next!!!", previous)
	// time.Sleep(10 * time.Second)
	return previous
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

func makeParams(s string) map[string]interface{} {
	var out map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	if err := dec.Decode(&out); err != nil {
		panic(err)
	}
	return out
}
