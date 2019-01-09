package balance

import (
	"fmt"
	"testing"
)

var s = []string{"a", "b", "c", "d", "e", "f"}

func TestBalance(t *testing.T) {
	Init(s)
	cm := make(map[string]int, 3)
	for i := 0; i < 6166; i++ {
		// fmt.Println(Next())
		cm[Next().host]++
	}
	fmt.Println(cm)
}

func BenchmarkBalance(t *testing.B) {
	Init(s)
	for i := 0; i < t.N; i++ {
		Next()
	}
}
