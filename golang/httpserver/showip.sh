#!/bin/bash

echo "show container ip addr"
Pid=$(docker inspect `docker ps |awk '{if($2=="eff4858/httpserver:1.0" && $NF=="httpserver")print $1}'` --format '{{ .State.Pid }}')
sudo nsenter -t $Pid -n ip addr
