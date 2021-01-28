package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func ManageData(c *Cluster, command []string) error {
	switch command[1] {
	case "put":
		if err := Create(c, command); err != nil {
			return err
		}
	case "get":
		if err := Read(c, command); err != nil {
			return err
		}
	case "update":
		if err := Update(c, command); err != nil {
			return err
		}
	case "delete":
		if err := Delete(c, command); err != nil {
			return err
		}
	}
	return nil
}

func Create(c *Cluster, command []string) error {
	data := map[string]string{}

	for i := 2; i < len(command); i++ {
		kv := strings.Split(command[i], "=")

		if len(kv) != 2 || len(kv[0]) == 0 || len(kv[1]) == 0 {
			return errors.New("bad key-value pair")
		}

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

	if err := resp.Body.Close(); err != nil {
		return err
	}

	return nil
}

func Read(c *Cluster, command []string) error {
	var keys []string

	for i := 2; i < len(command); i++ {
		keys = append(keys, command[i])
	}

	cl := &http.Client{}

	for _, v := range keys {
		url := fmt.Sprintf("http://%s/key/%s", c.Leader, v)

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

		if err := resp.Body.Close(); err != nil {
			return err
		}
	}

	return nil
}

func Update(c *Cluster, command []string) error {
	data := map[string]string{}

	for i := 2; i < len(command); i++ {
		kv := strings.Split(command[i], "=")

		if len(kv) != 2 || len(kv[0]) == 0 || len(kv[1]) == 0 {
			return errors.New("bad key-value pair")
		}

		data[kv[0]] = kv[1]
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	cl := &http.Client{}
	url := fmt.Sprintf("http://%s", c.Leader)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(b))

	resp, err := cl.Do(req)
	if err != nil {
		return err
	}

	r, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(r))

	if err := resp.Body.Close(); err != nil {
		return err
	}

	return nil
}

func Delete(c *Cluster, command []string) error {
	var keys []string

	for i := 2; i < len(command); i++ {
		keys = append(keys, command[i])
	}

	cl := &http.Client{}

	for _, v := range keys {
		url := fmt.Sprintf("http://%s/key/%s", c.Leader, v)
		req, err := http.NewRequest("DELETE", url, nil)

		resp, err := cl.Do(req)
		if err != nil {
			return err
		}

		r, _ := ioutil.ReadAll(resp.Body)
		fmt.Print(string(r))

		if err := resp.Body.Close(); err != nil {
			return err
		}
	}

	return nil
}
