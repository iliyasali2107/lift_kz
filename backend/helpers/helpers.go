package helpers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"time"
)

func PrintTime(expireAt int64) {

	// Convert the Unix timestamp to a time.Time object
	seconds := expireAt / 1000
	nanoseconds := (expireAt % 1000) * 1000000
	tm := time.Unix(seconds, nanoseconds)

	// Format the time in a human-readable format
	timeString := tm.Format("2006-01-02 15:04:05 MST")
	fmt.Println("Human-Readable Time:", timeString)
}

func ReadPdf(filePath string) (string, []byte, error) {
	// filePath := "someFile.pdf" // Update this with the correct file path
	dataBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "", nil, err
	}

	// Convert the file contents to a Base64-encoded string
	return base64.StdEncoding.EncodeToString(dataBytes), dataBytes, nil
}

func ErrorHandlingWithRerurn(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		return
	}
}

func DecodeBase64ToPDF(base64String, outputPath string) error {
	// Decode the Base64 string to bytes
	decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return err
	}

	// Write the decoded bytes to the output file
	err = ioutil.WriteFile(outputPath, decodedBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
