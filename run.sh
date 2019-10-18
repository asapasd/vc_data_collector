#!/bin/sh

mkdir -p log
today=$(date "+%Y%m%d")
/root/go/src/github.com/asapasd/vc_data_collector/main >> "./log/access_${today}.txt" 2>&1 &
echo "startup background"
ps aux | grep main
