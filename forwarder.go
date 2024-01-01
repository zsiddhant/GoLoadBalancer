package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ForwardRequest(req *http.Request, w http.ResponseWriter, server *Server, response chan http.Response) {
	if server.Alive == false {
		Print("Backend Server is not healthy")
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	url := fmt.Sprintf("http://%s:%d%s", server.Host, server.Port, req.RequestURI)
	proxyRes, err := RequestExecute(req.Method, url, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer proxyRes.Body.Close()
	response <- *proxyRes
}

func RequestExecute(method string, url string, body []byte) (*http.Response, error) {
	proxyReq, err := http.NewRequest(method, url, bytes.NewReader(body))
	proxyReq.Header = make(http.Header)
	client := &http.Client{}
	proxyRes, err := client.Do(proxyReq)
	return proxyRes, err
}
