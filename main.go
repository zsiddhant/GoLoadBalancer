package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

func addBackendHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var backendRequest AddBackendRequest
	if err := json.Unmarshal(body, &backendRequest); err != nil {
		fmt.Println("Invalid request")
		return
	}
	AddBackend(backendRequest)
}

func removeBackendHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var removeBackendRequest RemoveBackendRequest
	if err := json.Unmarshal(body, &removeBackendRequest); err != nil {
		fmt.Println("Invalid request")
		return
	}
	RemoveBackend(removeBackendRequest)
}

func forwardingHandler(w http.ResponseWriter, r *http.Request) {
	response := make(chan http.Response)
	backendToServe := RoundRobinExecution()
	go ForwardRequest(r, w, backendToServe, response)
	timeout := false
	for timeout == false {
		select {
		case res := <-response:
			Print("Response received ", res)
			respBody, _ := io.ReadAll(res.Body)
			w.WriteHeader(res.StatusCode)
			_, err := w.Write(respBody)
			if err != nil {
				Print("Error received from BE ", err)
			}
		case <-time.After(1 * time.Second):
			timeout = true
		}
	}

}

func main() {
	http.HandleFunc("/", forwardingHandler)
	http.HandleFunc("/add-backend", addBackendHandler)
	http.HandleFunc("/remove-backend", removeBackendHandler)

	port := flag.Int("PORT", 8080, "Port to start service on")
	url := fmt.Sprintf("localhost:%d", *port)
	err := http.ListenAndServe(url, nil)
	if err != nil {
		fmt.Println("Error while starting the server %s", err)
	}
}
