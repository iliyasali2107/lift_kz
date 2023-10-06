package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ResponseStruct struct {
	Nonce string `json:"nonce"`
}

func GetNonce() (string, error) {
	// Define the URL you want to send the POST request to
	url := os.Getenv("BASE_URL") + "/api/auth"

	// Create a map or struct containing the data you want to send as JSON
	data := map[string]interface{}{}

	// Convert the data to JSON format
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return "", err
	}

	// Create a new HTTP request with the POST method, URL, and request body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Set the request headers if needed (e.g., content-type)
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read and process the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// Print the response status code and body
	fmt.Println("Status Code:", resp.Status)
	fmt.Println("Response Body:", string(responseBody))
	var responseObj ResponseStruct
	if err := json.Unmarshal(responseBody, &responseObj); err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return "", err
	}
	nonceValue := responseObj.Nonce
	return nonceValue, nil
}
