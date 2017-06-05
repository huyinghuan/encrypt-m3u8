package utils

import (
	"io/ioutil"
	"log"
	"net/http"
)

func GetStream(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		// handle err
		log.Fatalln(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Error: response status is %d", resp.StatusCode)
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
