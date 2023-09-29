package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

var foundPort = make(chan int, 1)

func FetchPort(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("http://10.49.122.144:%d/ping", port)

	resp, err := http.Get(url)
	if err != nil {
		// fmt.Printf("Error fetching from %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Printf("Error reading response from %s: %v\n", url, err)
		return
	}
	foundPort <- port
	fmt.Printf("Response from %s: %s\n", url, string(body))
}

func main() {
	//Théo secret : 42078f34d2388a81df131352543155f748b436f78d947f109c8b28b00f4afe90
	var wg sync.WaitGroup

	for port := 0; port < 10000; port++ {
		wg.Add(1)
		go FetchPort(port, &wg)
	}
	go func() {
		wg.Wait()
		close(foundPort)
	}()
	foundPort := <-foundPort
	fmt.Printf("Found port: %d\n", foundPort)

	jsonBody := []byte(`{"user":"Théo"}`)

	// Adding the user
	resp1 := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/signup", foundPort), jsonBody)
	fmt.Println(resp1)

	// Check the user
	resp2 := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/check", foundPort), jsonBody)
	fmt.Println(resp2)

	// GetUserSecret
	jsonBody = []byte(`{"user":"Théo"}`)
	secret := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/getUserSecret", foundPort), jsonBody)

	payload := map[string]string{
		"user":   "Théo",
		"secret": secret,
	}

	// Convert the map to JSON
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	fmt.Println("secret : ", string(jsonBody))

	// GetUserLevel
	jsonBody = []byte(`{"user":"Théo", "secret": "42078f34d2388a81df131352543155f748b436f78d947f109c8b28b00f4afe90"}`)
	resp4 := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/getUserLevel", foundPort), jsonBody)
	fmt.Println(resp4)

	// GetUserPoints
	resp5 := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/getUserPoints", foundPort), jsonBody)
	fmt.Println(resp5)

	// Dabatase App : 72 44 90
	// Tiny Path [ctf-school Today is a good day innit ? ]
	// This is clearly not a binary : Q DC3 ) 1 4
	// Copy Trash 5FPprcvF-T75f91DQ2C
}

func FetchUrl(url string, jsonBody []byte) string {
	qBody := bytes.NewBuffer(jsonBody)
	resp, err := http.Post(url, "application/json", qBody)

	if err != nil {
		// fmt.Printf("Error fetching from %s: %v\n", url, err)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Printf("Error reading response from %s: %v\n", url, err)
		log.Fatal(err)
	}
	// fmt.Printf("Response from %s: %s\n", url, string(body))
	return string(body)
}
