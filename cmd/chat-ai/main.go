package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Getenv("CEREBRAS_API_KEY")
	url := "https://api.cerebras.ai/v1/chat/completions"

	// payload chat
	payload := map[string]interface{}{
		"model": "llama-3.3-70b",
		"messages": []map[string]string{
			{"role": "user", "content": "Halo, apa kabar?"},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
