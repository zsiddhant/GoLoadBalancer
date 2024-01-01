package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

func registerServerHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var backendRequest RegisterBackendRequest
	if err := json.Unmarshal(body, &backendRequest); err != nil {
		fmt.Println("Invalid request")
		return
	}
	RegisterServer(backendRequest)
}

func deRegisterServerHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var removeBackendRequest DeRegisterBackendRequest
	if err := json.Unmarshal(body, &removeBackendRequest); err != nil {
		fmt.Println("Invalid request")
		return
	}
	DeRegisterServer(removeBackendRequest)
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
	http.HandleFunc("/add-backend", registerServerHandler)
	http.HandleFunc("/remove-backend", deRegisterServerHandler)

	port := flag.Int("PORT", 8080, "Port to start service on")
	url := fmt.Sprintf("localhost:%d", *port)
	err := http.ListenAndServe(url, nil)
	if err != nil {
		fmt.Println("Error while starting the server %s", err)
	}
}
