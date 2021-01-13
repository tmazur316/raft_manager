package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ManageSnapshot(c *Cluster, command []string) error {
	switch command[1] {
	case "create":
		if err := CreateSnapshot(c, command); err != nil {
			return err
		}
	}
	return nil
}

func CreateSnapshot(c *Cluster, command []string) error {
	var servers []string

	for i := 2; i < len(command); i++ {
		servers = append(servers, command[i])
	}

	if len(servers) <= 0 {
		return errors.New("no server id specified")
	}

	var addr []string

	for _, v := range servers {
		a, err := GetAddr(c, v)
		if err != nil {
			return err
		}
		addr = append(addr, a)
	}

	cl := &http.Client{}

	for i, a := range addr {
		url := fmt.Sprintf("http://%s/snapshot", a)
		req, err := http.NewRequest("POST", url, nil)

		resp, err := cl.Do(req)
		if err != nil {
			return err
		}

		r, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("nodeId: %s\t%s\n", command[i+2], string(r))

		if err := resp.Body.Close(); err != nil {
			return err
		}
	}

	return nil
}

func GetAddr(c *Cluster, id string) (string, error) {
	for _, v := range c.Nodes {
		if v.ServerId == id {
			return v.ApiAddr, nil
		}
	}
	return "", errors.New("wrong server Id")
}
