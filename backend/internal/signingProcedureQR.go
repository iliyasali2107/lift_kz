package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"mado/models"
)

// QRSigningClientCMS represents a client for QR signing with CMS.
type QRSigningClientCMS struct {
	Description            string
	Attach                 bool
	BaseURL                string
	Retries                int
	DocumentsToSign        []map[string]interface{}
	DataURL, SignURL       string
	QRCode                 string
	EGovMobileLaunchLink   string
	EGovBusinessLaunchLink string
}

// NewQRSigningClientCMS creates a new instance of QRSigningClientCMS.
/**
  * Конструктор.
  *
  //* @param {String} description описание подписываемых данных.
  *
  //? @param {Boolean} [attach = false] следует ли включить в подпись подписываемые данные. желательно false так как нет смысоа включать данные в подпись. потомушто потом ее невытащиш
  *
  *! @param {String} [baseUrl = 'https://sigex.kz'] базовый URL сервиса SIGEX. лутше быть константой
*/
func NewQRSigningClientCMS(description string, attach bool, baseURL string) *QRSigningClientCMS {
	return &QRSigningClientCMS{
		Description: description,
		Attach:      attach,
		BaseURL:     baseURL,
		Retries:     25,
	}
}

// RegisterQRSinging registers the QR signing procedure and returns the QR code as a base64 encoded string.
func (qc *QRSigningClientCMS) RegisterQRSinging() (string, error) {
	data := map[string]string{
		"description": qc.Description,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	fmt.Println("data: ", (dataBytes))
	fmt.Println("POST URL: ", qc.BaseURL+"/api/egovQr")

	response, err := http.Post(qc.BaseURL+"/api/egovQr", "application/json", bytes.NewReader(dataBytes))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned status '%d: %s'", response.StatusCode, response.Status)
	} else {
		//*	expireAt - момент истечения срока действия процедуры подписания ЭЦП через QR в миллисекундах с UNIX Epoch;
		//*? qrCode - строка, закодированное в base64 изображение в формате PNG с закодированным QR кодом для приложения eGov mobile;
		//*! eGovMobileLaunchLink - ссылка для запуска eGov mobile на том же устройстве;
		//* eGovBusinessLaunchLink - ссылка для запуска eGov Business на том же устройстве;
		//*? dataURL - ссылка для отправки данных на подписание, это полная ссылка, ее следует использовать “как есть” для POST /api/egovQr/{qrId} - отправка данных на подписание;
		//*! signURL - ссылка для получения сформированных подписей, это полная ссылка, ее следует использовать “как есть” для GET /api/egovQr/{qrId} - получение подписей.

	}

	var responseJSON map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJSON)
	if err != nil {
		return "", err
	}

	var ErrorResp models.ErrorResponse
	errMessage := ErrorResp.GetHumanReadableErrorMessageByResponse(responseJSON)
	if errMessage != "" {
		return "", fmt.Errorf(errMessage)
	}

	if qrCode, ok := responseJSON["qrCode"].(string); ok {
		qc.QRCode = qrCode
	} else {
		return "", fmt.Errorf("qrCode is not a valid string")
	}

	if dataURL, ok := responseJSON["dataURL"].(string); ok {
		qc.DataURL = dataURL
	} else {
		return "", fmt.Errorf("dataURL is not a valid string")
	}

	if signURL, ok := responseJSON["signURL"].(string); ok {
		qc.SignURL = signURL
	} else {
		return "", fmt.Errorf("signURL is not a valid string")
	}

	if eGovMobileLaunchLink, ok := responseJSON["eGovMobileLaunchLink"].(string); ok {
		qc.EGovMobileLaunchLink = eGovMobileLaunchLink
	} else {
		return "", fmt.Errorf("eGovMobileLaunchLink is not a valid string")
	}

	if eGovBusinessLaunchLink, ok := responseJSON["eGovBusinessLaunchLink"].(string); ok {
		qc.EGovBusinessLaunchLink = eGovBusinessLaunchLink
	} else {
		return "", fmt.Errorf("eGovBusinessLaunchLink is not a valid string")
	}
	return qc.DataURL, nil
	// return qc.QRCode, nil
}

/*
//*! AddDataToSign adds a block of data for signing.
* @param {String[]} names массив имен подписываемого блока данных на разных языках
* [ru, kk, en]. Массив должен сожердать как минимум одну строку, в этом случае она будет
* использоваться для всех языков.
*
* @param {String | ArrayBuffer} data данные, которые нужно подписать, в виде строки Base64 либо
* ArrayBuffer.
*
* @param {Object[]} [meta = []] опциональный массив объектов метаданных, содержащих поля
* `"name"` и `"value"` со строковыми значениями.
*
* @param {Boolean} [isPDF = false] опциональная подсказка для приложения eGov mobile помогающая
* ему лучше подобрать приложение для отображения данных перед подписанием.
*/
func (qc *QRSigningClientCMS) AddDataToSign(names []string, data string, meta []map[string]string, isPDF bool) error {
	if len(names) == 0 {
		msg := models.QRSigningError{
			Message: "Data for signing is not provided correctly.",
			Details: "At least one name for the data to be signed must be specified.",
		}
		return errors.New(msg.Error())
	}

	documentToSign := map[string]interface{}{
		"id":     len(qc.DocumentsToSign) + 1,
		"nameRu": names[0],
		"nameKz": names[1],
		"nameEn": names[2],
		// "meta":   meta,
		"document": map[string]interface{}{
			"file": map[string]interface{}{
				"mime": "@file/pdf",
				"data": data,
			},
		},
		// "documentXml": "<groupId>2</groupId>",
	}

	qc.DocumentsToSign = append(qc.DocumentsToSign, documentToSign)
	return nil
}

// GetEGovMobileLaunchLink returns the link to launch the signing procedure in eGov mobile.
func (qc *QRSigningClientCMS) GetEGovMobileLaunchLink() string {
	return qc.EGovMobileLaunchLink
}

// GetEGovBusinessLaunchLink returns the link to launch the signing procedure in eGov Business.
func (qc *QRSigningClientCMS) GetEGovBusinessLaunchLink() string {
	return qc.EGovBusinessLaunchLink
}

func (qc *QRSigningClientCMS) GetSignatures(dataSentCallback func()) ([]string, error) {

	if len(qc.DocumentsToSign) == 0 {
		return nil, fmt.Errorf("Данные на подписание предоставлены не корректно. Не зарегистрировано ни одного блока данных для подписания.")
	}

	// Отправка данных
	signMethod := "CMS_SIGN_ONLY"
	if qc.Attach {
		signMethod = "CMS_WITH_DATA"
	}
	data := map[string]interface{}{
		"signMethod":      signMethod,
		"documentsToSign": qc.DocumentsToSign,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var response *http.Response
	for i := 0; i < qc.Retries; i++ {

		// fmt.Println("qc.DataURL: ", qc.DataURL)
		// fmt.Println("signMethod: ", signMethod)
		// fmt.Println("dataBytes: ", string(dataBytes))

		response, err = http.Post(qc.DataURL, "application/json", bytes.NewReader(dataBytes))

		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Не удалось получить ответа от сервера %s.", qc.DataURL)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Сервер вернул статус '%d: %s'", response.StatusCode, response.Status)
	}

	var responseJSON map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJSON)
	if err != nil {
		return nil, err
	}
	// fmt.Println("after POST request response: ", responseJSON) //*?i am getting this: map[expireAt:1.690350980919e+12 signURL:https://sigex.kz/api/egovQr/dVoCifgX9iZGFiBB]

	var ErrorResp models.ErrorResponse
	errMessage := ErrorResp.GetHumanReadableErrorMessageByResponse(responseJSON)
	if errMessage != "" {
		return nil, errors.Join(fmt.Errorf("Error in sending data for signing"), fmt.Errorf(errMessage))
	}

	if dataSentCallback != nil {
		dataSentCallback()
	}

	qc.SignURL = responseJSON["signURL"].(string)

	// Получение подписей    //TODO split this into separate function
	for i := 0; i < qc.Retries; i++ {
		response, err = http.Get(qc.SignURL)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Не удалось получить ответа от сервера %s.", qc.SignURL)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Сервер вернул статус '%d: %s'", response.StatusCode, response.Status)
	}

	var responseData models.ResponseGettingSignatureData
	return responseData.GetSignaturesFromResponse(response)
}
