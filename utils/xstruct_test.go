package utils

import "testing"

type User struct {
	Name string `json:"name"`
	Sex  int    `json:"sex"`
	Age  int    `json:"age"`
}

func TestStruct2Map(t *testing.T) {
	user := &User{
		Name: "aa",
		Sex:  1,
		Age:  20,
	}
	resMap := Struct2Map(user)

	t.Log(resMap)
}
