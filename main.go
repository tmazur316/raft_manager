package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	configFile := flag.String("config", "", "config file of a raft cluster")
	flag.Parse()
	log.SetOutput(os.Stderr)

	var config Cluster

	if *configFile != "" {
		c, err := LoadConfig(*configFile)
		if err != nil {
			log.Fatal("Failed to open config file")
		}
		config = *c

	} else {
		config = NewCluster(0)
	}

	fmt.Println("Raft manager started. Waiting for commands")
	fmt.Println("Available commands: cluster, data, snapshot, config")

	for {
		fmt.Print("$ ")
		r := bufio.NewReader(os.Stdin)
		line, err := r.ReadString('\n')
		if err != nil {
			log.Println("Wrong command")
		}

		cmd := strings.Fields(line)

		if len(cmd) == 0 {
			continue
		}

		switch cmd[0] {
		case "cluster":
			if err := ManageCluster(&config, cmd); err != nil {
				fmt.Println(err)
			}
		case "data":
			if err := ManageData(&config, cmd); err != nil {
				fmt.Println(err)
			}
		case "snapshot":
			if err := ManageSnapshot(&config, cmd); err != nil {
				fmt.Println(err)
			}
		case "config":
			if err := ManageConfig(&config, cmd, configFile); err != nil {
				fmt.Println(err)
			}
		case "exit":
			break
		default:
			fmt.Println("Unrecognized command")
			fmt.Println("Available commands: cluster, data, snapshot, config")
		}

		if cmd[0] == "exit" {
			break
		}
	}

	if err := ShutdownCluster(&config); err != nil {
		fmt.Println("Cluster shutdown failed")
	}

	for {
		if *configFile != "" {
			return
		}

		fmt.Println("No config file loaded. Do you want to load? [y/n] [default = y]")

		s, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println("Wrong command")
		}

		s = strings.TrimSuffix(s, "\n")

		if s == "n" {
			fmt.Println("Exit without saving config...")
			return
		}

		fmt.Println("Enter config file path")

		path, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println("Wrong file path")
		}

		path = strings.TrimSuffix(path, "\n")

		cmd := []string{"config", "file", "new", path}

		if err := ManageConfig(&config, cmd, configFile); err != nil {
			fmt.Println("Config file load failed")
		} else {
			break
		}
	}

	for {
		fmt.Println("saving current config...")

		if err := SaveConfig(&config, *configFile); err != nil {
			fmt.Println("attempt to save config failed. Retry [y/n]? [default = y]")

			s, err := bufio.NewReader(os.Stdin).ReadString('\n')

			s = strings.TrimSuffix(s, "\n")

			if err != nil {
				log.Println("Wrong command")
			}

			if s == "n" {
				fmt.Println("Exit without saving config...")
				break
			}
		} else {
			break
		}
	}

	fmt.Println("Config saved")
}
