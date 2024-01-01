package main

import (
	"fmt"
	"sync"
	"time"
)

var availableServers = make([]*Server, 0)

type RegisterBackendRequest struct {
	Name string
	Host string
	Port int
}

type DeRegisterBackendRequest struct {
	Name string
}

type Server struct {
	Name  string
	Host  string
	Port  int
	Alive bool
	mu    sync.Mutex
	Next  *Server
}

func RegisterServer(request RegisterBackendRequest) {
	name := request.Name
	host := request.Host
	port := request.Port

	backend := Server{Host: host, Name: name, Port: port, Alive: true}
	if len(availableServers) > 0 {
		temp := availableServers[len(availableServers)-1]
		temp.Next = &backend
	}
	availableServers = append(availableServers, &backend)
	go HealthCheckBackend(&backend)
}

func DeRegisterServer(request DeRegisterBackendRequest) {
	name := request.Name

	if len(availableServers) == 0 {
		Print("No servers available to remove")
		return
	}

	var indexToRemove int
	for i, server := range availableServers {
		if server.Name == name {
			indexToRemove = i
		}
	}

	if indexToRemove != 0 {
		if len(availableServers) == indexToRemove+1 {
			availableServers[indexToRemove-1].Next = nil
		} else {
			availableServers[indexToRemove-1].Next = availableServers[indexToRemove+1].Next
		}
	}
	availableServers = append(availableServers[:indexToRemove], availableServers[indexToRemove+1:]...)
}

func HealthCheckBackend(backend *Server) {
	for {
		url := fmt.Sprintf("http://%s:%d/", backend.Host, backend.Port)
		response, err := RequestExecute("GET", url, make([]byte, 0))
		if err != nil || response.StatusCode != 200 {
			backend.mu.Lock()
			backend.Alive = false
			DeRegisterServer(DeRegisterBackendRequest{Name: backend.Name})
			backend.mu.Unlock()
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func FetchServers() []*Server {
	return availableServers
}
