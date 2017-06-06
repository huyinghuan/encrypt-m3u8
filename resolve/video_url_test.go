package resolve

import "testing"
import "fmt"

func TestVideoURL(t *testing.T) {
	cdn, _ := GetCDNURL("3925408", "pad", "1")
	if cdn == "" {
		t.Fail()
	} else {
		fmt.Println(cdn)
	}

}
