#!/bin/bash

#$1 serverId
#$2 raftAddr
#$3 httpAddr

x-terminal-emulator -e /home/tomek/Pulpit/Go_Projects/bin/raft_tests -id=$1 -bootstrap=true -rAddr=$2 -httpAddr=$3
