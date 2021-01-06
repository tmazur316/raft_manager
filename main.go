package main

import (
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("/home/tomek/Pulpit/start_scripts/bootstrap.sh")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	cmd2 := exec.Command("/home/tomek/Pulpit/start_scripts/join.sh")
	err2 := cmd2.Start()
	if err2 != nil {
		log.Fatal(err2)
	}

	cmd.Wait()
	cmd2.Wait()
}
