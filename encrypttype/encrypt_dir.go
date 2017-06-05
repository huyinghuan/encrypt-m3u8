package encrypttype

import (
	"encry/encrypt"
	"encry/utils"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

func DecryptDir() {
	key := "0123456789123456"
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	startTime := time.Now()
	config, err := utils.ReadConfig(".yaml")
	if err != nil {
		log.Fatalln(err)
		return
	}
	sourceFiles, _ := ioutil.ReadDir(config.Download)
	var wg sync.WaitGroup
	wg.Add(len(sourceFiles))
	for _, file := range sourceFiles {
		fileName := file.Name()
		if strings.Index(fileName, ".ts") == -1 {
			continue
		}
		sourceFile := config.Download + "/" + fileName
		distFile := config.Decrypt + "/" + fileName
		fmt.Println(sourceFile, distFile)
		go func(sourceFile string, distFile string) {
			defer wg.Done()
			err := encrypt.CBCDecryptFile(sourceFile, distFile, key, iv)
			if err != nil {
				fmt.Println(err)
			}
		}(sourceFile, distFile)
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(startTime))
}

func EncryptDir() {
	key := "0123456789123456"
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	startTime := time.Now()
	config, err := utils.ReadConfig(".yaml")
	if err != nil {
		log.Fatalln(err)
		return
	}
	sourceFiles, _ := ioutil.ReadDir(config.Source)
	var wg sync.WaitGroup
	wg.Add(len(sourceFiles))
	for _, file := range sourceFiles {
		fileName := file.Name()
		sourceFile := config.Source + "/" + fileName
		distFile := config.Encrypt + "/" + fileName + ".cbc"
		go func(sourceFile string, distFile string) {
			defer wg.Done()
			encrypt.CBCEncryptFile(sourceFile, distFile, key, iv)
		}(sourceFile, distFile)
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(startTime))
}
