package main

import (
	"fmt"
	"sync"
	"time"
)

var availableBackends = make([]*Server, 0)

type AddBackendRequest struct {
	Name string
	Host string
	Port int
}

type RemoveBackendRequest struct {
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

func AddBackend(request AddBackendRequest) {
	name := request.Name
	host := request.Host
	port := request.Port

	backend := Server{Host: host, Name: name, Port: port, Alive: true}
	if len(availableBackends) > 0 {
		temp := availableBackends[len(availableBackends)-1]
		temp.Next = &backend
	}
	availableBackends = append(availableBackends, &backend)
	go HealthCheckBackend(&backend)
}

func RemoveBackend(request RemoveBackendRequest) {
	name := request.Name

	if len(availableBackends) == 0 {
		Print("No servers available to remove")
		return
	}

	var indexToRemove int
	for i, server := range availableBackends {
		if server.Name == name {
			indexToRemove = i
		}
	}

	if indexToRemove != 0 {
		if len(availableBackends) == indexToRemove+1 {
			availableBackends[indexToRemove-1].Next = nil
		} else {
			availableBackends[indexToRemove-1].Next = availableBackends[indexToRemove+1].Next
		}
	}
	availableBackends = append(availableBackends[:indexToRemove], availableBackends[indexToRemove+1:]...)
}

func HealthCheckBackend(backend *Server) {
	for {
		url := fmt.Sprintf("http://%s:%d/", backend.Host, backend.Port)
		response, err := RequestExecute("GET", url, make([]byte, 0))
		if err != nil || response.StatusCode != 200 {
			backend.mu.Lock()
			backend.Alive = false
			RemoveBackend(RemoveBackendRequest{Name: backend.Name})
			backend.mu.Unlock()
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func FetchBackend() []*Server {
	return availableBackends
}
