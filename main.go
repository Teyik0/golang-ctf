package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

	// 1. Signup
	user := "Théo"
	payload := map[string]string{
		"user": user,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Adding the user
	resp1 := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/signup", foundPort), jsonBody)
	fmt.Print(resp1)

	// Check the user
	resp2 := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/check", foundPort), jsonBody)
	fmt.Print(resp2)

	// GetUserSecret
	// secret := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/getUserSecret", foundPort), jsonBody)
	h := sha256.New()
	h.Write([]byte(user))
	secret := fmt.Sprintf("%x", h.Sum(nil))

	payload = map[string]string{
		"user":   user,
		"secret": secret,
	}

	// Convert the map to JSON
	jsonBody, err = json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// GetUserLevel
	for i := 0; i < 1; i++ {
		resp4 := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/getUserLevel", foundPort), jsonBody)
		fmt.Println(resp4)
	}

	// GetUserPoints
	for i := 0; i < 1; i++ {
		resp5 := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/getUserPoints", foundPort), jsonBody)
		fmt.Println(resp5)
	}

	// iNeedAHint
	var slicer []string
	for i := 0; i < 100; i++ {
		resp6 := FetchUrl(fmt.Sprintf("http://10.49.122.144:%d/iNeedAHint", foundPort), jsonBody)
		slicerAns := strings.Trim(resp6, "Coward over here asking for hints...\nHere you go, your random hint:\n")
		if !contains(slicer, slicerAns) {
			slicer = append(slicer, slicerAns)
		}
	}
	for _, element := range slicer {
		fmt.Println("| ", element)
	}
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
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
