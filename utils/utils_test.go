package utils

import (
	"log"
	"testing"
)

func TestMD5(t *testing.T) {
	password := MD5("123")
	log.Println(password)
	t.Errorf("%s", password)
}
