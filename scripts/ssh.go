package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	key := flag.String("key", "", "Private key file")
	host := flag.String("host", "", "host to connect to")
	cmd := flag.String("cmd", "", "command to run on remote host")
	flag.Parse()

	k, err := ioutil.ReadFile(*key)
	if err != nil {
		log.Fatal(err)
	}

	sign, err := ssh.ParsePrivateKey(k)
	if err != nil {
		log.Fatal(err)
	}

	p := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	hostKey, err := knownhosts.New(p)
	if err != nil {
		log.Fatal(err)
	}

	c := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(sign),
		},
		HostKeyCallback: hostKey,
	}

	cl, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", *host), c)
	if err != nil {
		log.Fatal(err)
	}
	defer cl.Close()

	s, err := cl.NewSession()
	if err != nil {
		log.Fatal()
	}

	m := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := s.RequestPty("xterm", 80, 24, m); err != nil {
		s.Close()
		log.Fatal(err)
	}

	out, err := s.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stdout, out)

	stderr, err := s.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	go io.Copy(os.Stderr, stderr)

	err = s.Run(*cmd)
}
