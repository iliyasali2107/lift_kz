package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type GenerateElectronicDocumentCardRequest struct {
	ID                                      string `json:"id"`
	RequestBody                             []byte `json:"requestBody"`
	FileName                                string `json:"fileName,omitempty"`
	WithoutDocumentVisualization            bool   `json:"withoutDocumentVisualization,omitempty"`
	WithoutSignaturesVisualization          bool   `json:"withoutSignaturesVisualization,omitempty"`
	WithoutQRCodesInSignaturesVisualization bool   `json:"withoutQRCodesInSignaturesVisualization,omitempty"`
	Language                                string `json:"language,omitempty"`
}

func NewGenerateElectronicDocumentCardRequest(id string, fileName string, withoutDocVis, withoutSigVis, withoutQRVis bool, language string, requestBody []byte) *GenerateElectronicDocumentCardRequest {
	return &GenerateElectronicDocumentCardRequest{
		ID:                                      id,
		FileName:                                fileName,
		WithoutDocumentVisualization:            withoutDocVis,
		WithoutSignaturesVisualization:          withoutSigVis,
		WithoutQRCodesInSignaturesVisualization: withoutQRVis,
		Language:                                language,
		RequestBody:                             requestBody,
	}
}

type GenerateElectronicDocumentCardResponse struct {
	DocumentID   string `json:"documentId"`
	DDC          string `json:"ddc"`
	DataArchived bool   `json:"dataArchived"`

	// ErrorResponse fields
	Message   string `json:"message,omitempty"`
	RequestID int64  `json:"requestID,omitempty"`
}

// Function to perform the POST request and receive the response
func (docCard GenerateElectronicDocumentCardRequest) GenerateElectronicDocumentCard(baseURL string) (*GenerateElectronicDocumentCardResponse, error) {
	// Construct the URL for the specific document ID
	url := fmt.Sprintf("%s/api/%s/buildDDC?fileName=%s&withoutDocumentVisualization=%t&withoutSignaturesVisualization=%t&withoutQRCodesInSignaturesVisualization=%t&language=%s",
		baseURL, docCard.ID, docCard.FileName, docCard.WithoutDocumentVisualization, docCard.WithoutSignaturesVisualization, docCard.WithoutQRCodesInSignaturesVisualization, docCard.Language)

	// Create the request body (empty for this POST request)

	// Create  new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(docCard.RequestBody))
	if err != nil {
		fmt.Println("POST REQUEST ERROR")
		return nil, err
	}

	// Set the Content-Type header to application/octet-stream
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(docCard.RequestBody)))
	// Set the Content-Length header with the correct length of the request body (empty for this POST request)

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

	// Parse the response body into the GenerateElectronicDocumentCardResponse struct
	var response GenerateElectronicDocumentCardResponse
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
	// fmt.Println("DDC:", response.DDC)
	fmt.Println("Data Archived:", response.DataArchived)

	return &response, nil
}
