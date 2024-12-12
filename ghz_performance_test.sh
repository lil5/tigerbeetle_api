#!/bin/bash

combinations=(
  "1 100000"
  "200 3000"
  "200 4000"
  "200 5000"
  "500 2000"
  "500 3000"
  "500 4000"
  "500 5000"
  "1000 10000"
)
for pair in "${combinations[@]}"; do
  conn=$(echo $pair | cut -d' ' -f1)
  conc=$(echo $pair | cut -d' ' -f2)
  echo "Testing with connections: $conn, concurrency: $conc"
  ghz --insecure \
    --call proto.TigerBeetle.CreateTransfers \
    --total 50000 \
    --concurrency $conc \
    --connections $conn \
    --data-file transfers.json \
    127.0.0.1:50051 | grep "Requests/sec:"
  sleep 5
done
