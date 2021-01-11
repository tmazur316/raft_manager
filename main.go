package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetOutput(os.Stderr)
	fmt.Println("Raft manager started. Waiting for commands")
	fmt.Println("Available commands: cluster")
	c := NewCluster(0)

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
			go ManageCluster(&c, cmd)

		}

		if cmd[0] == "exit" {
			break
		}
	}
}
