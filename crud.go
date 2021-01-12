package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func ManageData(c *Cluster, command []string) {
	switch command[1] {
	case "put":
		Create(c, command)
	case "get":
		Read(c, command)
	}
}

func Create(c *Cluster, command []string) error {
	data := map[string]string{}

	for i := 2; i < len(command); i++ {
		kv := strings.Split(command[i], "=")
		data[kv[0]] = kv[1]
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s", c.Leader)
	resp, _ := http.Post(url, "application/json", bytes.NewReader(b))

	r, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(r))

	//todo error handling of bad keys
	resp.Body.Close()

	return nil
}

func Read(c *Cluster, command []string) error {
	var keys []string

	for i := 2; i < len(command); i++ {
		keys = append(keys, command[i])
	}

	cl := &http.Client{}

	for _, v := range keys {
		url := fmt.Sprintf("http://%s/%s", c.Leader, v)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		resp, err := cl.Do(req)
		if err != nil {
			return err
		}

		r, _ := ioutil.ReadAll(resp.Body)
		fmt.Print(string(r))

		resp.Body.Close()
	}

	return nil
}

func Update() {

}

func Delete() {

}
