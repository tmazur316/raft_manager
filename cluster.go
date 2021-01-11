package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"time"
)

type Server struct {
	ServerId string
	NodeAddr string
	ApiAddr  string
}

type Cluster struct {
	Nodes []Server
}

func (c *Cluster) AddToConfig(s Server) {
	c.Nodes = append(c.Nodes, s)
}

func NewCluster(size int) Cluster {
	return Cluster{Nodes: make([]Server, size)}
}

func ManageCluster(c *Cluster, command []string) {
	//TODO error handling
	switch command[1] {
	case "bootstrap":
		BootstrapCluster(c, command, false)
	case "add":
		AddServer(c, command, false)
	case "start":
		StartCluster(c)
	}
}

func BootstrapCluster(c *Cluster, command []string, configStart bool) error {
	if len(command) != 8 {
		return errors.New("wrong command")
	}

	var nodeId = new(string)
	var nodeAddr = new(string)
	var httpAddr = new(string)

	for i := 2; i < len(command); i += 2 {
		switch command[i] {
		case "-nodeId":
			*nodeId = command[i+1]
		case "-nodeAddr":
			*nodeAddr = command[i+1]
		case "-httpAddr":
			*httpAddr = command[i+1]
		default:
			return errors.New("wrong command")
		}
	}

	cmd := exec.Command("scripts/bootstrap.sh", *nodeId, *nodeAddr, *httpAddr)
	err := cmd.Start()

	if err != nil {
		return errors.New("bootstrap error")
	}

	if !configStart {
		c.AddToConfig(Server{
			ServerId: *nodeId,
			NodeAddr: *nodeAddr,
			ApiAddr:  *httpAddr,
		})
	}

	if err := cmd.Wait(); err != nil {
		return errors.New("bootstrap error")
	}

	return nil
}

func AddServer(c *Cluster, command []string, configStart bool) error {
	if len(command) != 10 {
		return errors.New("wrong command")
	}

	var nodeId = new(string)
	var nodeAddr = new(string)
	var httpAddr = new(string)
	var joinAddr = new(string)

	for i := 2; i < len(command); i += 2 {
		switch command[i] {
		case "-nodeId":
			*nodeId = command[i+1]
		case "-nodeAddr":
			*nodeAddr = command[i+1]
		case "-httpAddr":
			*httpAddr = command[i+1]
		case "-joinAddr":
			*joinAddr = command[i+1]
		default:
			return errors.New("wrong command")
		}
	}
	cmd := exec.Command("scripts/join.sh", *nodeId, *nodeAddr, *httpAddr, *joinAddr)
	err := cmd.Start()

	if err != nil {
		return errors.New("bootstrap error")
	}

	if !configStart {
		c.AddToConfig(Server{
			ServerId: *nodeId,
			NodeAddr: *nodeAddr,
			ApiAddr:  *httpAddr,
		})
	}

	if err := cmd.Wait(); err != nil {
		return errors.New("bootstrap error")
	}

	return nil
}

func LoadConfig(filename string) (*Cluster, error) {
	var config Cluster

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(f)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(c *Cluster, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	e := json.NewEncoder(f)
	if err := e.Encode(c); err != nil {
		return err
	}

	return nil
}

func StartCluster(c *Cluster) error {
	//TODO error handling
	if len(c.Nodes) < 1 {
		return errors.New("cluster must contain at least one server")
	}

	l := c.Nodes[0]
	cmd := []string{
		"cluster",
		"bootstrap",
		"-nodeId",
		l.ServerId,
		"-nodeAddr",
		l.NodeAddr,
		"-httpAddr",
		l.ApiAddr,
	}

	BootstrapCluster(c, cmd, true)
	//TODO see if I can eliminate this sleep
	time.Sleep(2 * time.Second)

	for i := 1; i < len(c.Nodes); i++ {
		f := c.Nodes[i]
		cmd = []string{
			"cluster",
			"add",
			"-nodeId",
			f.ServerId,
			"-nodeAddr",
			f.NodeAddr,
			"-httpAddr",
			f.ApiAddr,
			"-joinAddr",
			l.ApiAddr,
		}

		AddServer(c, cmd, true)
	}

	return nil
}
