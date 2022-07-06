#!/bin/sh

pkill cmcc
pkill cmcc

# -1=debug, 0=info, 1=warn..., default to info
export GNET_LOGGING_LEVEL=0
export GNET_LOGGING_FILE="/Users/huangzhonghui/smsvg/logs/sms.log"
export CMCC_CONF_PATH="/Users/huangzhonghui/.cmcc.yaml"

mkdir -p /Users/huangzhonghui/smsvg/logs

# optional args --port 1234 --multicore=false
# default  args --port 9000 --multicore=true
nohup ./cmcc --port 9000 --multicore=true >panic.log 2>&1 &

sleep 3
tail -10 /Users/huangzhonghui/smsvg/logs/sms.log
sleep 7
top -pid "$(cat cmcc.pid)"
