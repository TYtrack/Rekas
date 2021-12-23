###
 # @Author: TYtrack
 # @Date: 2021-12-23 15:00:10
 # @LastEditors: TYtrack
 # @LastEditTime: 2021-12-23 15:03:42
 # @FilePath: /Rekas/data_structure/test_group.sh
### 
#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 &
./server -port=8002 &
./server -port=8003 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait