package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func ManageConfig(c *Cluster, command []string, configFile *string) error {
	switch command[1] {
	case "current":
		if err := CurrentConfig(c, command); err != nil {
			return err
		}
	case "file":
		if err := CurrentFile(command, configFile); err != nil {
			return err
		}
	case "save":
		if err := Save(c, configFile); err != nil {
			return err
		}
	}
	return nil
}

func CurrentConfig(c *Cluster, command []string) error {
	if len(command) > 2 {
		return errors.New("unrecognized arguments")
	}

	indent, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(indent))
	return nil
}

func CurrentFile(command []string, configFile *string) error {
	if len(command) > 2 {
		if command[2] == "new" {
			return ConfigFileChange(command, configFile)
		} else {
			return errors.New("unrecognized arguments")
		}
	}

	fmt.Println(*configFile)

	return nil
}

func ConfigFileChange(command []string, configFile *string) error {
	if len(command) != 4 {
		return errors.New("wrong number of arguments")
	}

	f, err := os.Stat(command[3])
	if errors.Is(err, os.ErrNotExist) || f.Mode().IsDir() {
		return err
	}

	*configFile = command[3]

	return nil
}

func Save(c *Cluster, configFile *string) error {
	if err := SaveConfig(c, *configFile); err != nil {
		return err
	}
	return nil
}
