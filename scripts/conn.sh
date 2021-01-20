#!/bin/bash

#$1 ssh.go absolute path
#$2 hostAddr
#$3 keyfile
#$4 command

x-terminal-emulator -e  go run "$1" -host=$2 -key=$3 -cmd="$4"

