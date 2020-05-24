package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"sync"
	"gopkg.in/yaml.v2"
)

// Request http
type Request struct {
	ArtCtn string `json:"ArtCtn"`
	ArtCteTm string `json:"ArtCteTm"`
	ChlCap string `json:"ChlCap"`
	ChlID string `json:"ChlId"`
}

//yml設定
var config = &conf{}
type conf []struct {
	URL     string `yaml:"url"`
	ChlID string `yaml:"ChlId"`
}

var wg sync.WaitGroup

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println("Unmarshal: %v", err)
	}
	return c
}

func main() {
	config.getConf()
	data, err := json.Marshal(config)

    if err != nil {
        fmt.Println("err:\t", err.Error())
        return
    }

    c := conf{}
	err = json.Unmarshal(data,&c)
	if err != nil {
		fmt.Println("err:\t", err.Error())
        return
 	}
	
	wg.Add(len(c))
	for i := 0; i < len(c); i++ {
        go run(c[i].URL,c[i].ChlID)
    }

	wg.Wait()

}

func run(url string, chlID string) error {
	method := "GET"
	client := &http.Client{}

	var createTime = ""
	var count = 0
	
	for {
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
		}

		res, err := client.Do(req)
		if err != nil {
			defer res.Body.Close()
		}

		body, err := ioutil.ReadAll(res.Body)
		if err == nil && len(body) > 0 {

			request := []Request{}
			_ = json.Unmarshal(body, &request)
			
			if request[0].ArtCteTm > createTime && request[0].ChlID == chlID {
				count ++
				createTime = request[0].ArtCteTm
				fmt.Println(request[0].ChlCap,fmt.Sprintf(" 第%v次PO文 ",count),request[0].ArtCteTm)
				fmt.Println(request[0].ArtCtn)
				fmt.Println("\n")
			}

		}
		time.Sleep(500 * time.Millisecond)
	}

	
	return nil
}
