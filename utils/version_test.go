package utils

import "testing"

func TestCompareMiddleRange(t *testing.T) {
	v := "2.0.1.1"
	vlo := "2.0.0.2"
	vhi := "2.0.2.1"
	if !CompareMiddleRange(v, vlo, vhi) {
		t.Fail()
	}
}

func TestGreaterThanEq(t *testing.T) {
	v1 := "2.0.1.1"
	v2 := "2.0.1.2"

	t.Log(GreaterThanEq(v1, v2))
}

func TestLessThanEq(t *testing.T) {
	v1 := "2.0.1.1"
	v2 := "2.0.1.0"

	t.Log(LessThanEq(v1, v2))
}
