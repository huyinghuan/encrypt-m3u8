package service

import (
	"encry/utils"
	"fmt"
	"log"
	"testing"
)

func TestExist(t *testing.T) {
	if e := utils.InitConfig(); e != nil {
		fmt.Println(e)
	}
	log.Println(IsExistM3U8("4", "3993398", "1"))
	log.Println(IsExistM3U8("4", "3993399", "1"))
}
