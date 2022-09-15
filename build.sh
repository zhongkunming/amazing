#!/bin/bash

work=/opt
name=service-hub
running=$(pgrep $name | wc -l)
if [ "$running" -gt 0 ];
  then pgrep ${name} | awk '{print $1}' | xargs kill -9
fi
cd $work/$name || (echo "文件不存在" && exit)
go mod tidy
go build
nohup ./${name} > /dev/null 2>&1 &
