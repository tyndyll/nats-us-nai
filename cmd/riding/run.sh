#/bin/bash

if [[ -z "$(docker ps -q -f name=nats-main)" ]]; then
  docker run -d --name nats-main -p 4222:4222 -p 6222:6222 -p 8222:8222 nats
fi

go run platform/main.go &
sleep 2

for i in {1..5}; do
  go run driver/main.go &
done

sleep 5

go run customer/main.go

echo "Customer completed"

for i in `ps aux | grep "go run" | awk '{print $2}'`; do kill -9 $i 2>/dev/null ; done
for i in `ps aux | grep "go-build" | awk '{print $2}'`; do kill -9 $i 2>/dev/null ; done

docker rm -f  nats-main