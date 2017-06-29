package controller

import (
	"encry/encrypt"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/kataras/iris/context"
)

func Decrypt(ctx context.Context) {
	key := ctx.URLParam("key")
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	sourceFiles, _ := ioutil.ReadDir("/Users/hyh/Downloads/encrypt/encrypt/")
	var wg sync.WaitGroup
	wg.Add(len(sourceFiles))
	for _, file := range sourceFiles {
		fileName := file.Name()
		if strings.Index(fileName, ".ts") == -1 {
			continue
		}
		sourceFile := "/Users/hyh/Downloads/encrypt/encrypt/" + fileName
		distFile := "/Users/hyh/Downloads/encrypt/decrypt/" + fileName
		go func(sourceFile string, distFile string) {
			defer wg.Done()
			err := encrypt.CBCDecryptFile(sourceFile, distFile, key, iv)
			if err != nil {
				fmt.Println(err)
			}
		}(sourceFile, distFile)
	}
	wg.Wait()
}
