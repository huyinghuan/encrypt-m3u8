package utils

import (
	"log"
	"testing"
)

func TestRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(RandString(16))
	}
}
