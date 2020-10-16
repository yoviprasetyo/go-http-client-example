package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// checkURL method.
func checkURL(api API) {
	switch api.Method {
	case "GET":
		dataGet(api)
	case "POST":
		dataPost(api)
	}
}

func dataGet(api API) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", api.URL, nil)

	if err != nil {
		fmt.Println(err)
	}

	if len(api.Headers) > 0 {
		for _, value := range api.Headers {
			req.Header.Add(value.Key, value.Value)
		}
	}

	resp, err := client.Do(req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("URL:", api.URL)
	if len(api.Headers) > 0 {
		for _, value := range api.Headers {
			fmt.Println("Header:", value.Key, value.Value)
		}
	}

	fmt.Println(string(body))
	fmt.Println("")
}

func dataPost(api API) {

	var (
		body   io.Reader
		client = &http.Client{}
	)

	switch {
	case api.JSON != nil:
		api.Headers = append(api.Headers, Header{Key: "Content-Type", Value: "application/json"})
	case api.FormData != nil:
		requestBody, err := json.Marshal(api.FormData)
		if err != nil {
			fmt.Println(err)
		}

		api.Headers = append(api.Headers, Header{Key: "Content-Type", Value: "application/json"})
		body = bytes.NewBuffer(requestBody)
	}

	req, err := http.NewRequest("POST", api.URL, body)
	if err != nil {
		fmt.Println(err)
	}

	if len(api.Headers) > 0 {
		for _, value := range api.Headers {
			req.Header.Add(value.Key, value.Value)
		}
	}

	resp, err := client.Do(req)

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("URL:", api.URL)
	if len(api.Headers) > 0 {
		for _, value := range api.Headers {
			fmt.Println("Header:", value.Key, value.Value)
		}
	}

	fmt.Println(string(responseBody))
	fmt.Println("")
}

// API struct.
type API struct {
	URL         string
	Method      string
	FormData    H
	JSON        H
	QueryString H
	Headers     []Header
}

// H struct.
type H map[string]interface{}

// Header struct.
type Header struct {
	Key   string
	Value string
}

// APIKey of ruang api.
var APIKey string

func main() {
	godotenv.Load(".env")
	APIKey = os.Getenv("API_KEY")
	apis := []API{
		API{
			URL:    "https://official-joke-api.appspot.com/random_joke",
			Method: "GET",
		},
		API{
			URL:    "https://ruangapi.com/api/v1/provinces",
			Method: "GET",
			Headers: []Header{
				Header{
					Key:   "Authorization",
					Value: APIKey,
				},
			},
		},
		API{
			URL:    "https://ruangapi.com/api/v1/search-engine",
			Method: "POST",
			Headers: []Header{
				Header{
					Key:   "Authorization",
					Value: APIKey,
				},
			},
			FormData: H{
				"url": "https://soizee.com",
				"se":  "google",
			},
		},
		API{
			URL:    "https://ruangapi.com/api/v1/shopee",
			Method: "POST",
			Headers: []Header{
				Header{
					Key:   "Authorization",
					Value: APIKey,
				},
			},
			FormData: H{
				"username": "muhammadhendra",
				"take":     "50",
			},
		},
		API{
			URL:    "https://ruangapi.com/api/v1/currency",
			Method: "POST",
			Headers: []Header{
				Header{
					Key:   "Authorization",
					Value: APIKey,
				},
			},
			FormData: H{
				"code": "gold-antam",
			},
		},
	}

	var wg sync.WaitGroup

	for _, api := range apis {
		wg.Add(1)
		go func(api API) {
			defer wg.Done()
			checkURL(api)
		}(api)
	}

	wg.Wait()
}
