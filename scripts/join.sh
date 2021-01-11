#!/bin/bash

#$1 Id 
#$2 nodeAddr 
#$3 httpAddr 
#$4 joinAddr 

x-terminal-emulator -e /home/tomek/Pulpit/Go_Projects/bin/raft_tests -id=$1 -rAddr=$2 -httpAddr=$3 -joinAddr=$4
