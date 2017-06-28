package service

import (
	"encry/utils"
	"fmt"
	"os"
)

func IsExistM3U8(vedioID string, terminalType string, resolution string) bool {
	config := utils.ReadConfig()
	dir := config.Originm3u8
	m3u8Filename := fmt.Sprintf(config.M3u8rule, terminalType, vedioID, resolution)
	m3u8FilePath = fmt.Sprintf(dir, m3u8Filename)
	if _, err := os.Stat(m3u8FilePath); os.IsExist(err) {
		return true
	}
	return false
}
