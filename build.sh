#!/bin/bash

work=/opt/apps/mcs
name=mcs
running=$(pgrep $name | wc -l)
if [ "$running" -gt 0 ]; then
  pgrep ${name} | awk '{print $1}' | xargs kill -9
fi
cd $work || (echo ${work}" does not exist" && exit 1)
go mod tidy
go build -o $name
nohup ./${name} > /dev/null 2>&1 &
