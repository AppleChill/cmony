package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	yamlv2 "gopkg.in/yaml.v2"
)

// Request http
type Request struct {
	ArtCtn string `json:"ArtCtn"`
}

//yml設定
var config = &conf{}
type conf struct {
	URL     string `yaml:"url"`
}

var data1 = ""
var count = 0

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println("yamlFile.Get err   #%v ", err)
	}
	err = yamlv2.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println("Unmarshal: %v", err)
	}
	return c
}

func main() {
	config.getConf()
	for {
		run()
		// err := run()
		// if err != nil {
		// 	break
		// }
		time.Sleep(500 * time.Millisecond)
	}
}

func run() error {
	url := config.URL
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err == nil {

		request := []Request{}
		_ = json.Unmarshal(body, &request)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return err
		// }
		
		if request[0].ArtCtn != data1 {
			count ++
			data1 = request[0].ArtCtn
			fmt.Println(fmt.Sprintf("第%v次PO文",count))
			fmt.Println(request[0].ArtCtn)
			fmt.Println("\n")
		}

	}
	return nil
}
