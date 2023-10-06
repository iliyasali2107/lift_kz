package model

import (
	"mado/internal"
)

type AuthResponse struct {
	UserID                  string           `json:"userId"`
	BusinessID              string           `json:"businessId"`
	Email                   string           `json:"email"`
	Subject                 string           `json:"subject"`
	SubjectStructure        [][]Attribute    `json:"subjectStructure"`
	SubjectAltName          string           `json:"subjectAltName"`
	SubjectAltNameStructure []SubjectAltName `json:"subjectAltNameStructure"`
	SignAlgorithm           string           `json:"signAlgorithm"`
	KeyStorage              string           `json:"keyStorage"`
	PolicyIds               []string         `json:"policyIds"`
	ExtKeyUsages            []string         `json:"extKeyUsages"`
}

type Attribute struct {
	Oid        string `json:"oid"`
	Name       string `json:"name"`
	ValueInB64 bool   `json:"valueInB64"`
	Value      string `json:"value"`
}

type SubjectAltName struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AuthRequest struct {
	Nonce     *string `json:"nonce"`
	Signature *string `json:"signature"`
	External  bool    `json:"external"`
}

type LoginRequirements struct {
	// Context  context.Context              `json:"context"`
	QrSigner   *internal.QRSigningClientCMS `json:"qrsigner"`
	Nonce      *string                      `json:"nonce"`
	Is_manager *bool                        `json:"is_manager"`
}
