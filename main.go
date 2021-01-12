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
	fmt.Println("Available commands: cluster")

	for {
		fmt.Print("$ ")
		r := bufio.NewReader(os.Stdin)
		line, err := r.ReadString('\n')
		if err != nil {
			log.Println("Wrong command")
		}

		cmd := strings.Fields(line)

		switch cmd[0] {
		case "exit":
			break
		case "cluster":
			ManageCluster(&config, cmd)
		case "data":
			ManageData(&config, cmd)
		}

		if cmd[0] == "exit" {
			break
		}
	}

	for {
		fmt.Println("saving current config...")

		if err := SaveConfig(&config, "/home/tomek/Pulpit/config"); err != nil {
			fmt.Println("attempt to save config failed. Retry [y/n]? [default = y]")

			s, err := bufio.NewReader(os.Stdin).ReadString('\n')

			if err != nil {
				log.Println("Wrong command")
			}

			if s == "n" {
				fmt.Println("Exit without saving config...")
				break
			}
		}

		break
	}

	fmt.Println("Config saved")
}
