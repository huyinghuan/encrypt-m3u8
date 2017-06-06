package utils

import (
	"log"
	"net/url"
	"strings"
)

//GetDirname cd go获取文件夹名
func GetDirname(uri string) string {
	urlObj, err := url.Parse(uri)
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	urlPathArr := strings.Split(urlObj.Path, "/")
	urlPathArr = urlPathArr[:len(urlPathArr)-1]
	if urlObj.Scheme == "" {
		return strings.Join(urlPathArr, "/")
	}
	return urlObj.Scheme + "://" + urlObj.Host + strings.Join(urlPathArr, "/")
}
