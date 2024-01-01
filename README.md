# GoLoadBalancer

This is a simple load balancer created in Go, working on a layer-seven - application load balancer. It uses round-robin algorithm to route HTTP requests from clients to a pool of HTTP servers.



### Register Server
```shell
curl --location --request GET 'localhost:8080/add-backend' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name" : "B1",
    "host" : "localhost",
    "port" : 8081
}'
```

### De-Register Server

```shell
curl --location --request GET 'localhost:8080/remove-backend' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name" : "B1"
}'
```

### Health Check

For each server registered will be pinged every 30 sec, if failed to respond with 200 OK HTTP status code will be removed routing.

### Access Logs

For each request received and response being forwarded back to client, it will be logged in /tmp/go-load-balancer-access-logs.txt


## Build

```sh
go build .
```

## Run

```shell
./GoLoadBalancer -PORT={port}
```