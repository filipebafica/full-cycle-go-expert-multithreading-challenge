package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const APITimeout = time.Second * 1

func makeRequest(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func main() {
	service1Ch := make(chan string)
	service2Ch := make(chan string)

	service1Url := "https://brasilapi.com.br/api/cep/v1/01310000"
	service2Url := "http://viacep.com.br/ws/01310000/json/"

	go func() {
		resp, err := makeRequest(service1Url)
		if err != nil {
			slog.Error(err.Error())
		}
		defer resp.Body.Close()

		var decoded map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&decoded)
		bytes, err := json.MarshalIndent(decoded, "", "  ")
		if err != nil {
			slog.Error(err.Error())
			return
		}
		service1Ch <- string(bytes)
	}()

	go func() {
		resp, err := makeRequest(service2Url)
		if err != nil {
			slog.Error(err.Error())
		}
		defer resp.Body.Close()

		var decoded map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&decoded)
		bytes, err := json.MarshalIndent(decoded, "", "  ")
		if err != nil {
			slog.Error(err.Error())
			return
		}
		service2Ch <- string(bytes)
	}()

	select {
	case result := <-service1Ch:
		fmt.Println("Received response from Brasil API: ", result)
		break
	case result := <-service2Ch:
		fmt.Println("Received response from ViaCep API: ", result)
		break
	case <-time.After(APITimeout):
		fmt.Println("Timed out waiting for response from APIs")
		break
	}
}
