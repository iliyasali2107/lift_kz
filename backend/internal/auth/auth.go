package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"mado/internal"
)

const (
	block   = "Блок данных на подпись"
	baseUrl = "https://sigex.kz"
)

func PreparationStep() (egovMobileLink *string, qrSigner *internal.QRSigningClientCMS, nonce *string) {

	body, _ := json.Marshal(map[string]interface{}{})

	response, err := http.Post(baseUrl+"/api/auth", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error getting аутентификация, подготовительный этап:", err)
		return nil, nil, nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Server returned status '%d: %s\n'", response.StatusCode, response.Status)
		return nil, nil, nil
	}

	var responseJSON map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJSON)
	if err != nil {
		fmt.Println("decode err:", err)
		return nil, nil, nil
	}
	fmt.Println(responseJSON)

	if nonce, ok := responseJSON["nonce"].(string); ok {
		qrSigner := internal.NewQRSigningClientCMS("Тестовое подписание", false, baseUrl)
		err = qrSigner.AddDataToSign([]string{block, block, block}, nonce, nil, true)
		if err != nil {
			fmt.Println("Could not read file: ", err)
			return nil, nil, nil
		}

		// Register QR signing
		dataURL, err := qrSigner.RegisterQRSinging()
		if err != nil {
			fmt.Println("RegisterQRSinging Error:", err)
			return nil, nil, nil
		}
		fmt.Println("First man RegisterQRSinging dataURL: ", dataURL)

		eGovMobileLaunchLink := qrSigner.GetEGovMobileLaunchLink()
		fmt.Println("For log-in eGov Mobile Launch Link:", eGovMobileLaunchLink) //todo pull to front
		//todo response front
		return &eGovMobileLaunchLink, qrSigner, &nonce
	} else {
		fmt.Println("err did't have nonce in resp:", err)
		return nil, nil, nil
	}

}

func GetNonceSignature(qrSigner *internal.QRSigningClientCMS) *string {
	signatures, err := qrSigner.GetSignatures(nil)
	if err != nil {
		fmt.Println("GetSignatures Error:", err)
		return nil
	}

	return &signatures[0]
}
