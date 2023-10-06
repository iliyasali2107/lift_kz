package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Define the structure for the response
type AddSignatureResponse struct {
	DocumentID                       string        `json:"documentId"`
	SignID                           int           `json:"signId"`
	Data                             string        `json:"data,omitempty"`
	AutomaticallyCreatedUserSettings *UserSettings `json:"automaticallyCreatedUserSettings,omitempty"`
	DataArchived                     bool          `json:"dataArchived"`
	CanBeArchived                    bool          `json:"canBeArchived"`

	// ErrorResponse fields
	Message   string `json:"message,omitempty"`
	RequestID int64  `json:"requestID,omitempty"`
}

// UserSettingsResponse represents the user settings response
type UserSettings struct {
	UserID                    string `json:"userId,omitempty"`
	EmailNotificationsEnabled bool   `json:"emailNotificationsEnabled,omitempty"`
	Email                     string `json:"email,omitempty"`
	ModifiedAt                int64  `json:"modifiedAt,omitempty"`
}

type Request struct {
	SignType  string `json:"signType,omitempty"`
	Signature string `json:"signature"`
}

// Function to perform the POST request and receive the response
func AddSignatureToDocument(id, signature string, baseURL string) (*AddSignatureResponse, error) {
	// Construct the URL for the specific document ID
	url := fmt.Sprintf("%s/api/%s", baseURL, id)

	// Create the request body
	requestBody := Request{
		SignType:  "cms",
		Signature: signature,
	}

	// Convert the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("MARSHAL ERROR")
		return nil, err
	}

	// Create a new request with the JSON data
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("POST REQUEST ERROR")
		return nil, err
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do ERROR")
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code for errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	// Parse the response body into the AddSignatureResponse struct
	var response AddSignatureResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Decode ERROR")
		return nil, err
	} else if err != nil || response.Message != "" {
		// If there's an error message, return the ErrorResponse
		return nil, errors.Join(fmt.Errorf("RequestID: %s%s", strconv.FormatInt(response.RequestID, 10), " Error msg: "), fmt.Errorf(response.Message))

	}

	// Access the response data
	fmt.Println("Document ID:", response.DocumentID)
	fmt.Println("Sign ID:", response.SignID)
	fmt.Println("Data:", response.Data)
	if response.AutomaticallyCreatedUserSettings != nil {
		fmt.Println("User ID:", response.AutomaticallyCreatedUserSettings.UserID)
		fmt.Println("Email Notifications Enabled:", response.AutomaticallyCreatedUserSettings.EmailNotificationsEnabled)
		fmt.Println("Email:", response.AutomaticallyCreatedUserSettings.Email)
		fmt.Println("Modified At:", response.AutomaticallyCreatedUserSettings.ModifiedAt)
	}
	fmt.Println("Data Archived:", response.DataArchived)
	fmt.Println("Can Be Archived:", response.CanBeArchived)

	return &response, nil
}
