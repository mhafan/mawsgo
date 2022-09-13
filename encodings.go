package mawsgo

import (
	"encoding/base64"

	"github.com/google/uuid"
)

//
func MAWSUUID() string {
	//
	return uuid.NewString()
}

//
func Base64EncodeString(str *string) string {
	//
	return base64.StdEncoding.EncodeToString([]byte(*str))
}

//
func Base64DecodeString(str *string) []byte {
	//
	p, _ := base64.StdEncoding.DecodeString(*str)

	//
	return p
}
