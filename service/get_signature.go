package service

import "net/http"

func GetSignature(req *http.Request) string {
	return "123456"
}
