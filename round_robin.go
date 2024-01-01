package main

import (
	"fmt"
)

var lastServedBackend *Server

func RoundRobinExecution() *Server {
	backends := FetchServers()
	if len(backends) == 0 {
		fmt.Println("No available backend")
	}

	var backendToServe *Server
	if lastServedBackend == nil {
		backendToServe = backends[0]
	} else {
		backendToServe = lastServedBackend.Next
		if backendToServe == nil {
			backendToServe = backends[0]
		}
	}
	lastServedBackend = backendToServe
	return backendToServe
}
