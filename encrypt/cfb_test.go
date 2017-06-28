package encrypt

import (
	"fmt"
	"log"
	"testing"
)

func TestCFBCryptString(t *testing.T) {
	str, err := CFBEncryptString([]byte("0123456789123456"), "helloworldhelloworldhelloworldhelloworld123456789")
	if err != nil {
		log.Println(err)
		t.Fail()
		return
	}
	fmt.Println(str)
	decryptStr, err := CFBDecryptString([]byte("0123456789123456"), str)
	if err != nil {
		log.Println(err)
		t.Fail()
		return
	}
	fmt.Println(decryptStr)

}
