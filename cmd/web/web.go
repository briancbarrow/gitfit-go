package web

import (
	"encoding/json"
	"log"
)

type NonceKey string

const NonceValue = NonceKey("nonce")

type HxHeader struct {
	XCSRFToken string `json:"X-CSRFToken"`
}

func ConvertHeaderToJSON(token string) string {
	header := HxHeader{XCSRFToken: token}
	jsonData, err := json.Marshal(header)
	if err != nil {
		log.Println(err)
	}
	return string(jsonData)

}
