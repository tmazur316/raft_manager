package main

import (
	"errors"
	"os/exec"
)

type Server struct {
	ServerId string
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
	switch command[1] {
	case "bootstrap":
		BootstrapCluster(c, command)
	case "add":
		AddServer(c, command)
	}

}

func BootstrapCluster(c *Cluster, command []string) error {
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

	c.AddToConfig(Server{
		ServerId: *nodeId,
		ApiAddr:  *httpAddr,
	})

	if err := cmd.Wait(); err != nil {
		return errors.New("bootstrap error")
	}

	return nil
}

func AddServer(c *Cluster, command []string) error {
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

	c.AddToConfig(Server{
		ServerId: *nodeId,
		ApiAddr:  *httpAddr,
	})

	if err := cmd.Wait(); err != nil {
		return errors.New("bootstrap error")
	}

	return nil
}
