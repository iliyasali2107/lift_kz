package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mado/helpers"
	"time"

	"net/http"
	"strconv"
	// "mado/internal"
)

// //////////////////////////////
type MetaData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DocumentFile struct {
	Mime string `json:"mime"`
	Data string `json:"data"`
}

type Document struct {
	File        DocumentFile `json:"file"`
	DocumentXml string       `json:"documentXml"`
}

type DocumentToSign struct {
	ID       int        `json:"id"`
	NameRu   string     `json:"nameRu"`
	NameKz   string     `json:"nameKz"`
	NameEn   string     `json:"nameEn"`
	Meta     []MetaData `json:"meta"`
	Document Document   `json:"document"`
}

type ResponseGettingSignatureData struct {
	Message         string           `json:"messgage,omitempty"`
	SignMethod      string           `json:"signMethod"`
	DocumentsToSign []DocumentToSign `json:"documentsToSign"`
}

func (responseData *ResponseGettingSignatureData) GetSignaturesFromResponse(response *http.Response) ([]string, error) {
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server returned status '%d: %s'", response.StatusCode, response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	if responseData.Message != "" {
		return nil, errors.Join(fmt.Errorf("Error in geting signature"), fmt.Errorf(responseData.Message))
	}

	if responseData.SignMethod != "CMS_SIGN_ONLY" {
		return nil, fmt.Errorf("Invalid sign method in the response")
	}

	// Extract the signatures from the responseData.DocumentsToSign
	signatures := make([]string, 0, len(responseData.DocumentsToSign))
	for _, doc := range responseData.DocumentsToSign {
		signatures = append(signatures, doc.Document.File.Data)
	}
	fmt.Printf("DocumentsToSign ID: %d\n", responseData.DocumentsToSign[0].ID)
	return signatures, nil
}

type QRResponse struct {
	ExpireAt               int64  `json:"expireAt"`
	DataURL                string `json:"dataURL"`
	SignURL                string `json:"signURL"`
	EGovMobileLaunchLink   string `json:"eGovMobileLaunchLink"`
	EGovBusinessLaunchLink string `json:"eGovBusinessLaunchLink"`
	QRCode                 string `json:"qrCode"`
}

func (qr *QRResponse) printQRResponse(response *http.Response) {
	// Print the status code and status text
	fmt.Printf("Status: %d %s\n", response.StatusCode, response.Status)

	// Print all the response headers
	for key, value := range response.Header {
		fmt.Printf("%s: %s\n", key, value)
	}

	// Read and print the response body
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	defer response.Body.Close()

	var qrResponse QRResponse
	err = json.Unmarshal(bodyBytes, &qrResponse)
	if err != nil {
		fmt.Println("Error unmarshaling response body:", err)
		return
	}
	fmt.Println("ExpireAt:", qrResponse.ExpireAt)
	helpers.PrintTime(qrResponse.ExpireAt)
	fmt.Println("DataURL:", qrResponse.DataURL)
	fmt.Println("SignURL:", qrResponse.SignURL)
	fmt.Println("EGovMobileLaunchLink:", qrResponse.EGovMobileLaunchLink)
	fmt.Println("EGovBusinessLaunchLink:", qrResponse.EGovBusinessLaunchLink)
	fmt.Println("QRCode:", qrResponse.QRCode)
}

/*

type SendDataToSignResponse struct {
	ExpireAt int64  `json:"expireAt"`
	SignURL  string `json:"signURL"`
}

func (qc *QRSigningClientCMS) SendDataToSign(qrId string) (string, error) {
	response, err := http.Post(qrId, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned status '%d: %s'", response.StatusCode, response.Status)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}
	defer response.Body.Close()

	var sendDataToSignResponse SendDataToSignResponse
	err = json.Unmarshal(bodyBytes, &sendDataToSignResponse)
	if err != nil {
		fmt.Println("Error unmarshaling response body:", err)
		return "", err
	}

	return sendDataToSignResponse.SignURL, nil

}

*/

//*!for  POST /api - регистрация нового документа в системе

type DocumentRegistrationRequest struct {
	Title              string                   `json:"title"`
	Description        string                   `json:"description"`
	SignType           string                   `json:"signType,omitempty"`
	Signature          string                   `json:"signature"`
	EmailNotifications EmailNotificationOptions `json:"emailNotifications,omitempty"` //todo maybe add email notifications
	Settings           DocumentSettings         `json:"settings,omitempty"`
}

type EmailNotificationOptions struct {
	To []string `json:"to"`
}

type DocumentSettings struct {
	Private                   bool     `json:"private,omitempty"`
	SignaturesLimit           int      `json:"signaturesLimit,omitempty"`
	SwitchToPrivateAfterLimit bool     `json:"switchToPrivateAfterLimitReached,omitempty"`
	Unique                    []string `json:"unique,omitempty"`
	StrictSignersRequirements bool     `json:"strictSignersRequirements,omitempty"`
	SignersRequirements       []struct {
		IIN string `json:"iin"`
	} `json:"signersRequirements,omitempty"`
}

type DocumentRegistrationResponse struct {
	DocumentID string `json:"documentId"`
	SignID     int    `json:"signId"`
	Data       string `json:"data,omitempty"`

	// ErrorResponse fields
	Message   string `json:"message,omitempty"`
	RequestID int64  `json:"requestID,omitempty"`
}

func NewDocumentRegistrationRequest(title, description, signType, signature string, emailNotifications []string /* settings DocumentSettings */) DocumentRegistrationRequest {
	return DocumentRegistrationRequest{
		Title:       title,       //* must have
		Description: description, //* must have
		SignType:    signType,    //* must have default cms
		Signature:   signature,   //* must have
		EmailNotifications: EmailNotificationOptions{
			To: emailNotifications,
		},
		// Settings: settings,
	}
}

func (docRegReq *DocumentRegistrationRequest) RegisterDocument(baseURL string) (*DocumentRegistrationResponse, error) {
	// Create the request payload
	/*	requestBody := DocumentRegistrationRequest{
		Title:       "document title",
		Description: "document description",
		SignType:    "cms",
		Signature:   signature,
		//EmailNotifications: EmailNotificationOptions{
		//		To: []string{"user@example.com", "other@example.com"},
		//	},
		Settings: DocumentSettings{
			Private:                   false,
			SignaturesLimit:           0,
			SwitchToPrivateAfterLimit: false,
			Unique:                    []string{"iin"},
			StrictSignersRequirements: false,
			// SignersRequirements: []struct {
			// IIN string `json:"iin"`
			// }{
			// {IIN: "IIN112233445566"},
			// },
		},
	}*/

	// Marshal the request payload to JSON
	requestJSON, err := json.Marshal(docRegReq)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/api", baseURL) // Replace {id} with the appropriate ID in the actual URL
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	// Set the appropriate headers (Content-Type and Authorization if required)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var errResp ErrorResponse
	// Check the response status code for errors
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		errResp, err := errResp.ParseErrorResponse(bodyBytes)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("Error parsing error response: "), err)
		}
		return nil, fmt.Errorf("API error: %s", errResp.GetHumanReadableErrorMessage())
	}

	// Parse the response
	var response DocumentRegistrationResponse
	// err = json.Unmarshal(bodyBytes, &response)
	// if err != nil {
	// 	return nil, err
	// }
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	} else if err == nil && response.Message != "" {
		// If there's an error message, return the ErrorResponse
		return nil, errors.Join(fmt.Errorf("RequestID: %s%s", strconv.FormatInt(response.RequestID, 10), " Error msg: "), fmt.Errorf(response.Message))
	}

	fmt.Println("Document ID:", response.DocumentID)
	fmt.Println("Sign ID:", response.SignID)
	fmt.Println("Data:", response.Data) ///todo ETO POLE EMPTY

	return &response, nil
}

// *Registration HESH
// Define the structure for the response
type DocumentHashesResponse struct {
	DocumentID                       string                     `json:"documentId"`
	SignedDataSize                   int                        `json:"signedDataSize"`
	Digests                          map[string]string          `json:"digests"`
	EmailNotifications               *EmailNotificationResponse `json:"emailNotifications,omitempty"`
	AutomaticallyCreatedUserSettings *UserSettingsResponse      `json:"automaticallyCreatedUserSettings,omitempty"`
	DataArchived                     bool                       `json:"dataArchived"`

	// ErrorResponse fields
	Message   string `json:"message,omitempty"`
	RequestID int64  `json:"requestID,omitempty"`
}

type EmailNotificationResponse struct {
	Attached bool   `json:"attached"`
	Message  string `json:"message,omitempty"`
}

type UserSettingsResponse struct {
	UserID                    string `json:"userId"`
	EmailNotificationsEnabled bool   `json:"emailNotificationsEnabled"`
	Email                     string `json:"email"`
	ModifiedAt                int64  `json:"modifiedAt"`
}

// Function to perform the POST request and receive the response FixingDocumentHashes
func (docRes DocumentHashesResponse) FixingDocumentHashes(id string, document []byte, baseURL string) (*DocumentHashesResponse, error) {
	// Construct the URL for the specific document ID
	url := fmt.Sprintf("%s/api/%s/data", baseURL, id)

	// Create a new request with the document data
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(document))
	if err != nil {
		return nil, err
	}

	// Set the Content-Type header to application/octet-stream
	req.Header.Set("Content-Type", "application/octet-stream")

	// Set the Content-Length header with the correct size of the document data
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(document)))

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code for errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}
	var response DocumentHashesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	} else if err != nil || response.Message != "" {
		// If there's an error message, return the ErrorResponse
		return nil, errors.Join(fmt.Errorf("RequestID: %s%s", strconv.FormatInt(response.RequestID, 10), " Error msg: "), fmt.Errorf(response.Message))

	}

	/*
		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Parse the response body into the DocumentDataResponse struct
		var dataResponse DocumentDataResponse
		err = json.Unmarshal(body, &dataResponse)
		if err != nil {
			return nil, err
		}*/

	return &response, nil
}

type Survey struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	Status       bool       `json:"status"`
	Rka          string     `json:"rka"`
	RcName       string     `json:"rc_name"`
	Adress       string     `json:"adress"`
	Questions    []Question `json:"questions"`
	QuestionsStr []string   `json:"questions_str"`
	CreatedAt    time.Time  `json:"created_at"`
	UserId       int        `json:"user_id"`
	// Ids          []int      `json:"ids,omitempty"`
	Answers []Answer `json:"answers,omitempty"`
}

type UserSurvey struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Status    bool       `json:"status"`
	Adress    string     `json:"adress"`
	Questions []Question `json:"questions,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UserId    int        `json:"user_id"`
}

type Question struct {
	Id          int      `json:"id"`
	Description string   `json:"description"`
	Answers     []Answer `json:"answers"`
}

type Answer struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}
type SaveSurveyRequest struct {
	Id        int         `json:"id"`
	Questions []Question2 `json:"questions,omitempty"`
	UserId    int         `json:"user_id"`
}

type Question2 struct {
	Id       int `json:"id"`
	AnswerId int `json:"answer_id"`
}
