package service

import "net/http"

func GetSignature(req *http.Request) string {
	return "123456"
}

func MatchSignature(req *http.Request, signature string) bool {
	return "123456" == signature
}
