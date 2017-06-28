package service

import (
	"encry/utils"
	"fmt"
	"io/ioutil"
	"os"
)

func IsExistM3U8(terminalType string, vedioID string, resolution string) bool {
	config := utils.ReadConfig()
	dir := config.Originm3u8
	m3u8Filename := fmt.Sprintf(config.M3u8rule, terminalType, vedioID, resolution)
	m3u8FilePath := fmt.Sprintf(dir, m3u8Filename)
	if _, err := os.Stat(m3u8FilePath); err == nil {
		return true
	} else {
		return false
	}

}

func GetExistM3U8(terminalType string, vedioID string, resolution string) ([]byte, error) {
	config := utils.ReadConfig()
	dir := config.Originm3u8
	m3u8Filename := fmt.Sprintf(config.M3u8rule, terminalType, vedioID, resolution)
	m3u8FilePath := fmt.Sprintf(dir, m3u8Filename)
	return ioutil.ReadFile(m3u8FilePath)
}
