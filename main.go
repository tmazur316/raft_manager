package main

import (
	"fmt"
	"log"
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

func main() {

	//CreateNewCluster
	c := NewCluster(0)

	cmd, err := BootstrapCluster(&c, "node1", "", "")
	if err != nil {
		log.Fatal(err)
	}

	cmd2, err := AddServer(&c, "node2", "127.0.0.1:6000", "127.0.0.1:6500", "127.0.0.1:5500")
	if err != nil {
		log.Fatal(err)
	}

	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	if err = cmd2.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cluster: %s\n", c)
}

func BootstrapCluster(c *Cluster, Id, nodeAddr, httpAddr string) (*exec.Cmd, error) {
	if nodeAddr == "" {
		nodeAddr = "127.0.0.1:5000"
	}

	if httpAddr == "" {
		httpAddr = "127.0.0.1:5500"
	}

	cmd := exec.Command("/home/tomek/Pulpit/start_scripts/bootstrap.sh", Id, nodeAddr, httpAddr)
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	c.AddToConfig(Server{
		ServerId: Id,
		ApiAddr:  httpAddr,
	})

	return cmd, err
}

func AddServer(c *Cluster, Id, nodeAddr, httpAddr, joinAddr string) (*exec.Cmd, error) {
	cmd := exec.Command("/home/tomek/Pulpit/start_scripts/join.sh", Id, nodeAddr, httpAddr, joinAddr)
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	c.AddToConfig(Server{
		ServerId: Id,
		ApiAddr:  httpAddr,
	})

	return cmd, err
}
