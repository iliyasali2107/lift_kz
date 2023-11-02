package helpers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
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

func init() {

}

func CurrentDateModel() Date {
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		fmt.Println(err)
		return Date{}
	}

	currentTime := time.Now().In(loc)
	day := strconv.Itoa(currentTime.Day())
	monthInt := int(currentTime.Month())
	year := strconv.Itoa(currentTime.Year())
	hour := strconv.Itoa(currentTime.Hour())
	minute := strconv.Itoa(currentTime.Minute())

	months := make(map[int]string, 12)
	months[1] = "Январь"
	months[2] = "Февраль"
	months[3] = "Март"
	months[4] = "Апрель"
	months[5] = "Май"
	months[6] = "Июнь"
	months[7] = "Июль"
	months[8] = "Август"
	months[9] = "Сентябрь"
	months[10] = "Октябрь"
	months[11] = "Ноябрь"
	months[12] = "Декабрь"

	// Print the extracted values

	if len(minute) == 1 {
		minute = "0" + minute
	}

	if len(hour) == 1 {
		hour = "0" + hour
	}

	if len(day) == 1 {
		day = "0" + day
	}

	time := fmt.Sprintf("%s:%s", hour, minute)
	month := months[monthInt]
	year = year[2:]
	date := Date{}

	date.Day = day
	date.Month = month
	date.Time = time
	date.Year = year

	return date

}

type Date struct {
	Day   string
	Month string
	Time  string
	Year  string
}
