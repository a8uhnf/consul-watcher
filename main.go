package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mitchellh/consulstructure"
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
	// watchExample()
	testWatchEx()
	<-stop
}

func testWatchEx() {
	// Create a configuration struct that'll be filled by Consul.
	type Config struct {
		Addr     string
		DataPath string `consul:"data_path"`
	}

	// Create our decoder
	updateCh := make(chan interface{})
	errCh := make(chan error)
	decoder := &consulstructure.Decoder{
		Target:   &Config{},
		Prefix:   "hello",
		UpdateCh: updateCh,
		ErrCh:    errCh,
	}

	// Run the decoder and wait for changes
	go decoder.Run()
	for {
		select {
		case v := <-updateCh:
			fmt.Printf("Updated config: %#v\n", v.(*Config))
		case err := <-errCh:
			fmt.Printf("Error: %s\n", err)
		}
	}
}
