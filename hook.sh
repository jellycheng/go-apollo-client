#!/bin/bash

#bash ./hook.sh 环境代号 阿波罗项目名appid 集群名 命名空间名 env文件
dev=$1
appid=$2
cluster=$3
namespace=$4
envfile=$5

echo "接收参数： ${dev} | ${appid} | ${cluster} | ${namespace} | ${envfile}"

