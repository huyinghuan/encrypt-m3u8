package utils

import (
	"fmt"
	"testing"
)

func TestGetDirname(t *testing.T) {
	dirname := GetDirname("/a/b/c?Asdasd")
	if dirname != "/a/b" {
		t.Fail()
		fmt.Println(dirname)
	}
}

func TestGetDirname_2(t *testing.T) {
	dirname := GetDirname("http://aww.com/a/b/c?Asdasd")
	if dirname != "http://aww.com/a/b" {
		t.Fail()
		fmt.Println(dirname)
	}
}
