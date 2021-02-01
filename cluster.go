package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	ServerId string
	NodeAddr string
	ApiAddr  string
}

type Cluster struct {
	Nodes  []Server
	Leader string
}

func (c *Cluster) AddToConfig(s Server) {
	c.Nodes = append(c.Nodes, s)
}

func NewCluster(size int) Cluster {
	return Cluster{Nodes: make([]Server, size)}
}

func ManageCluster(c *Cluster, command []string) error {
	switch command[1] {
	case "bootstrap":
		if err := BootstrapCluster(c, command, false); err != nil {
			return err
		}
	case "add":
		if err := AddServer(c, command, false); err != nil {
			return err
		}
	case "start":
		if err := StartCluster(c); err != nil {
			return err
		}
	case "remove":
		if err := RemoveServer(c, command); err != nil {
			return err
		}
	}
	return nil
}

func BootstrapCluster(c *Cluster, command []string, configStart bool) error {
	if len(command) != 8 {
		return errors.New("wrong command")
	}

	var Id = new(string)
	var nodeAddr = new(string)
	var httpAddr = new(string)

	for i := 2; i < len(command); i += 2 {
		switch command[i] {
		case "-nodeId":
			*Id = command[i+1]
		case "-nodeAddr":
			*nodeAddr = command[i+1]
		case "-httpAddr":
			*httpAddr = command[i+1]
		default:
			return errors.New("wrong command")
		}
	}

	host := strings.Split(*nodeAddr, ":")[0]
	p := strings.Split(*httpAddr, ":")

	e := fmt.Sprintf("./raft_tests -id=%s -rAddr=%s -httpAddr=%s:%s -bootstrap=true", *Id, *nodeAddr, host, p[1])

	fmt.Printf("Enter a private key file path for host %s\n", host)

	path, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return err
	}

	path = strings.TrimSuffix(path, "\n")

	abs, err := filepath.Abs("./scripts/ssh.go")
	if err != nil {
		return err
	}

	cmd := exec.Command("./scripts/conn.sh", abs, p[0], path, e)
	if err := cmd.Start(); err != nil {
		return err
	}

	if !configStart {
		c.AddToConfig(Server{
			ServerId: *Id,
			NodeAddr: *nodeAddr,
			ApiAddr:  *httpAddr,
		})
		c.Leader = *httpAddr
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func AddServer(c *Cluster, command []string, configStart bool) error {
	if len(command) != 10 {
		return errors.New("wrong command")
	}

	var Id = new(string)
	var nodeAddr = new(string)
	var httpAddr = new(string)
	var joinAddr = new(string)

	for i := 2; i < len(command); i += 2 {
		switch command[i] {
		case "-nodeId":
			*Id = command[i+1]
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

	host := strings.Split(*nodeAddr, ":")[0]
	p := strings.Split(*httpAddr, ":")

	e := fmt.Sprintf("./raft_tests -id=%s -rAddr=%s -httpAddr=%s:%s -joinAddr=%s", *Id, *nodeAddr, host, p[1], *joinAddr)

	fmt.Printf("Enter a private key file path for host %s\n", host)

	path, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return err
	}

	path = strings.TrimSuffix(path, "\n")

	abs, err := filepath.Abs("./scripts/ssh.go")
	if err != nil {
		return err
	}

	cmd := exec.Command("./scripts/conn.sh", abs, p[0], path, e)
	if err := cmd.Start(); err != nil {
		return err
	}

	if !configStart {
		c.AddToConfig(Server{
			ServerId: *Id,
			NodeAddr: *nodeAddr,
			ApiAddr:  *httpAddr,
		})
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func RemoveServer(c *Cluster, command []string) error {
	if len(command) != 6 {
		return errors.New("wrong command")
	}

	var nodeId = new(string)
	var httpAddr = new(string)

	for i := 2; i < len(command); i += 2 {
		switch command[i] {
		case "-nodeId":
			*nodeId = command[i+1]
		case "-httpAddr":
			*httpAddr = command[i+1]
		default:
			return errors.New("wrong command")
		}
	}

	url := fmt.Sprintf("http://%s/remove/%s", c.Leader, *nodeId)
	cl := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusMethodNotAllowed {
		addr, _ := ioutil.ReadAll(resp.Body)
		UpdateLeader(c, string(addr))
		resp.Body.Close()

		if err := RemoveServer(c, command); err != nil {
			return err
		}

		return nil
	}

	if err := removeFromConfig(c, *nodeId); err != nil {
		return err
	}

	if err := resp.Body.Close(); err != nil {
		return err
	}

	return nil
}

func StartCluster(c *Cluster) error {
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

	if err := BootstrapCluster(c, cmd, true); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

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

		if err := AddServer(c, cmd, true); err != nil {
			return err
		}
	}

	return nil
}

func ShutdownCluster(c *Cluster) error {
	for _, node := range c.Nodes {
		url := fmt.Sprintf("http://%s/shutdown", node.ApiAddr)
		cl := &http.Client{}
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return err
		}

		resp, err := cl.Do(req)
		if err != nil {
			return err
		}

		if err := resp.Body.Close(); err != nil {
			return err
		}
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

	indent, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	if _, err := f.Write(indent); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func removeFromConfig(c *Cluster, nodeID string) error {
	var index int
	var found bool

	for i, v := range c.Nodes {
		if v.ServerId == nodeID {
			index = i
			found = true
			break
		}
	}

	if !found {
		return errors.New("unable to remove node from config")
	}

	c.Nodes = append(c.Nodes[:index], c.Nodes[index+1:]...)
	return nil
}

func UpdateLeader(c *Cluster, leaderAddr string) {
	for _, node := range c.Nodes {
		if node.NodeAddr == leaderAddr {
			c.Leader = node.ApiAddr
			break
		}
	}
}
