package mawsgo

import "encoding/base64"

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
