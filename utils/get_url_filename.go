package utils

import (
	"net/url"
	"strings"
)

//GetURLFilename 获取URL请求文件名
func GetURLFilename(uri string) string {
	urlObj, _ := url.Parse(uri)
	arr := strings.Split(urlObj.Path, "/")
	return arr[len(arr)-1]
}
