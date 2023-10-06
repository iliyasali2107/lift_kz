package models

import (
	"encoding/json"
	"fmt"
)

// QRSigningError represents an error related to QR signing.
type QRSigningError struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func (e *QRSigningError) Error() string {
	return fmt.Sprintf("QRSigning Error: %s - %s", e.Message, e.Details)
}

type ErrorResponse struct {
	Message   string `json:"message,omitempty"`
	RequestID int64  `json:"requestID,omitempty"`
}

// Map of error messages to their human-readable versions
var errorMap = map[string]string{
	"Failed to build certificate chain":                      "Не удалось построить цепочку сертификатов.",
	"Failed to build digital document card":                  "Не удалось сформировать карточку электронного документа.",
	"Failed to parse JSON":                                   "Не удалось разобрать JSON запрос.",
	"Failed to parse digital document card":                  "Не удалось разобрать карточку электронного документа.",
	"Failed to parse signature":                              "Не удалось разобрать подпись.",
	"Failed to prepare archived data":                        "Ошибка подготовки архива.",
	"Feature is not available on current subscription level": "Функция недоступна на данном уровне подписки.",
	"HTTP Content-Length request header not set":             "Не указан HTTP заголовок Content-Length.",
	"Invalid API route":                                      "Некорректный маршрут.",
	"Invalid HTTP request headers":                           "Некорректный HTTP заголовок запроса.",
	"Invalid HTTP request method":                            "Некорректный HTTP метод.",
	"Invalid JSON request structure":                         "Некорректная структура JSON запроса.",
	"Invalid QR code logo format":                            "Неверный формат логотипа QR кода.",
	"Invalid QR code logo resolution":                        "Неподдерживаемое разрешение логотипа QR кода.",
	"Invalid QR code logo size":                              "Неверный размер логотипа QR кода.",
	"Invalid QR signing state":                               "Неверное состояние QR подписания.",
	"Invalid URL query parameter":                            "Некорректный URL параметр запроса.",
	"Invalid authentication nonce":                           "Неверный nonce аутентификации.",
	"Invalid certificate status":                             "Некорректный статус сертификата.",
	"Invalid document":                                       "Некорректный документ.",
	"Invalid document Title":                                 "Неверный заголовок документа.",
	"Invalid document identifier":                            "Некорректный идентификатор документа.",
	"Invalid file name":                                      "Некорректное имя файла.",
	"Invalid format of provided document settings":           "Некорректный формат объекта настроек документа.",
	"Invalid signature":                                      "Некорректная подпись.",
	"Invalid signature export format":                        "Некорректный формат экспорта.",
	"Invalid signature identifier":                           "Некорректный идентификатор подписи.",
	"Nonconforming order of elements in one of ASN.1 SET tags of the signature": "Неверный порядок элементов в одном из ASN.1 тегов SET подписи.",
	"Not enough TLS certificates":                   "Недостаточно TLS сертификатов.",
	"Not enough authorities":                        "Недостаточно полномочий.",
	"Not supported document type":                   "Неподдерживаемый тип документа.",
	"OCSP server problem":                           "Проблема с OCSP сервером.",
	"One of the provided IINs is in invalid format": "Некорректный формат одного из ИИН переданных в запросе.",
	"Access denied":                                 "Доступ запрещен.",
	"Already authenticated with TLS certificate":    "Аутентификация уже выполнена по TLS сертификату.",
	"Archive not found":                             "Архив не найден.",
	"Authentication required":                       "Требуется аутентификация.",
	"Bad signer certificate":                        "Плохой сертификат подписавшего.",
	"Certificate issuer is unknown":                 "Издатель сертификата неизвестен.",
	"Certificate signature algorithm not supported": "Алгоритм подписи сертификата не поддерживается.",
	"Digest algorithm not supported":                "Алгоритм хеширования не поддерживается.",
	"Digests values update required":                "Необходимо обновить хеши документа.",
	"Digital document card contains signatures that correspond to different registered documents": "В карточке электронного документа присутствуют подписи, зарегистрированные под разными электронными документами.",
	"Digital document card has signature in invalid format":                                       "В карточке электронного документа присутствует подпись не подходящего формата.",
	"Digital document card original document is not registered":                                   "Документ в карточке электронного документа не зарегистрирован на сервисе.",
	"Document data archival is started by another user":                                           "Архивирование данных выполняется другим пользователем.",
	"Document digests are already known":                                                          "Хеши документа уже известны.",
	"Document digests are not known":                                                              "Хеши документа не известны.",
	"Document does not have signatures in required format":                                        "Под документом не зарегистрировано ни одной подписи подходящего для данной операции формата.",
	"Document has more signatures than new limit":                                                 "Количество подписей под документом превышает устанавливаемое ограничение.",
	"Document has signatures in not supported by this operation format":                           "Под документом зарегистрированы подписи не подходящего для данной операции формата.",
	"Document not found":                "Документ не найден.",
	"Document signatures limit reached": "Превышено ограничение на количество подписей.",
	"Failed to archive document data":   "Ошибка архивирования документа.",
	"One of the provided TLS certificates can not be used for authentication":                     "Один из TLS сертификатов не может быть использован для аутентификации.",
	"One of the provided URL addresses is in invalid format":                                      "Некорректный формат одного из URL адресов.",
	"One of the provided authority OIDs is in invalid format":                                     "Некорректный формат одного из OID-ов.",
	"One of the provided certificate indices is invalid":                                          "Один из индексов сертификатов не верен.",
	"One of the provided notification email addresses is in invalid format":                       "Некорректный формат одного из адресов электронной почты для отправки уведомлений.",
	"Other endpoint connected":                                                                    "Подключена другая конечная точка.",
	"QR code versions lower than 11 are not allowed":                                              "Версии QR кодов ниже 11 не поддерживаются.",
	"QR signing operation is not active":                                                          "Операция QR подписания не активна.",
	"QR signing operation timeout":                                                                "Таймаут операции QR подписания.",
	"Request body is too large":                                                                   "Размер запроса слишком велик.",
	"Request field size is too large":                                                             "Размер поля слишком велик.",
	"Request rate limit reached":                                                                  "Достигнуто ограничение частоты запросов.",
	"Signature algorithm not supported":                                                           "Алгоритм подписи не поддерживается.",
	"Signature contains invalid OCSP data":                                                        "Подпись содержит поврежденные или некорректные данные OCSP.",
	"Signature contains invalid TSP time stamp":                                                   "Подпись содержит поврежденную или некорректную метку времени TSP.",
	"Signature does not conform to document settings requirements":                                "Подпись не удовлетворяет требованиям в настройках документа.",
	"Signature does not correspond to the document":                                               "Подпись не соответствует документу.",
	"Signature type is not supported":                                                             "Не поддерживаемый тип подписи.",
	"Signer certificate expired or not yet valid":                                                 "Срок действия сертификата уже истек, либо еще не наступил.",
	"Some of digital document card signatures are not registered":                                 "В карточке электронного документа присутствует не зарегистрированная на сервисе подпись.",
	"TLS certificate used for authentication was disabled":                                        "TLS сертификат, использованный для аутентификации, отключен.",
	"TSP server problem":                                                                          "Проблема с TSP сервером.",
	"The processed data size does not match the one passed in HTTP Content-Length request header": "Размер обработанных данных не равен размеру, переданному в HTTP заголовке Content-Length.",
	"This signature has already been submitted":                                                   "Данная подпись уже была зарегистрирована.",
	"Too many notification email addresses provided":                                              "Указано слишком много адресов электронной почты для рассылки уведомлений.",
	"Too much users use QR signing":                                                               "Слишком много пользователей выполняют QR подписание.",
	"Unexpected error":                                                                            "Непредвиденная ошибка.",
	"User does not represent an organization":                                                     "Пользователь не является представителем организации.",
	"eGov mobile data exchange failed":                                                            "Ошибка обмена данными с eGov mobile.",
}

func (errResp *ErrorResponse) ParseErrorResponse(jsonData []byte) (*ErrorResponse, error) {
	err := json.Unmarshal(jsonData, &errResp)
	if err != nil {
		return nil, err
	}

	return errResp, nil
}

// GetHumanReadableErrorMessage converts the given error message to its human-readable version
func (errResp *ErrorResponse) GetHumanReadableErrorMessageByResponse(response map[string]interface{}) string {
	if msg, ok := response["message"].(string); ok {
		if humanReadableMsg, exists := errorMap[msg]; exists {
			return humanReadableMsg
		}
		// If the error message is not found in the map, return the original error message
		return msg
	}
	// If the 'message' field does not exist in the response, return an empty string
	return ""
}

func (errResp *ErrorResponse) GetHumanReadableErrorMessage() string {
	if humanReadableMsg, exists := errorMap[errResp.Message]; exists {
		return humanReadableMsg
	}
	// If the error message is not found in the map, return the original error message
	return errResp.Message
}
