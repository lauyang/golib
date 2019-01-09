package utils

import (
	"testing"
)

func Test_APHash(t *testing.T) {
	var str [6]string = [6]string{"hello world", "hi", "hi1", "no", "hi11", "h111"}

	for _, s := range str {
		hash := APHash([]byte(s))
		t.Log(s, ":", hash)
	}
}
