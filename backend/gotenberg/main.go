package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {

	// URL of the Gotenberg service.
	url := "http://localhost:3000/forms/chromium/convert/html"

	
	// Open the HTML file to upload.
	htmlFile, err := os.Open("gotenberg/temp.html")
	if err != nil {
		fmt.Println("Error opening HTML file:", err)
		return
	}
	defer htmlFile.Close()

	// Create a buffer to store the form data.
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form field for the HTML file.
	htmlPart, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Copy the HTML file content to the form field.
	_, err = io.Copy(htmlPart, htmlFile)
	if err != nil {
		fmt.Println("Error copying file content:", err)
		return
	}

	// Close the form data writer.
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing form writer:", err)
		return
	}

	// Create an HTTP POST request with the form data.
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error performing HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code.
	if resp.StatusCode == http.StatusOK {
		// Save the resulting PDF to a file.
		pdfFile, err := os.Create("gotenberg/my.pdf")
		if err != nil {
			fmt.Println("Error creating PDF file:", err)
			return
		}
		defer pdfFile.Close()

		_, err = io.Copy(pdfFile, resp.Body)
		if err != nil {
			fmt.Println("Error saving PDF:", err)
			return
		}

		fmt.Println("PDF conversion successful. Result saved to my.pdf")
	} else {
		fmt.Printf("HTTP request failed with status code: %d\n", resp.StatusCode)
	}
}
