package main

import (
	"io"
	"net/http"
)

func server() {

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()
		io.WriteString(w, values["url"][0])
	})
	http.ListenAndServe("localhost:8080", nil)
	/*
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
	*/
}
