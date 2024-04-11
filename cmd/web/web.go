package web

import (
	"encoding/json"
	"log"
)

type NonceKey string

const NonceValue = NonceKey("nonce")

type HxHeader struct {
	XCSRFToken string `json:"X-CSRF-Token"`
	HXDate     string `json:"HX-Date"`
}

func ConvertHeaderToJSON(token string) string {
	header := HxHeader{XCSRFToken: token}
	jsonData, err := json.Marshal(header)
	if err != nil {
		log.Println(err)
	}

	return string(jsonData)
}
