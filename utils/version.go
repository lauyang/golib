package utils

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	VERSION   string
	COMMIT    string
	BRANCH    string
	GOVERSION string
	BUILDTIME string
)

func String() string {
	return fmt.Sprintf("version:%s\nCOMMIT:%s\nBRANCH:%s\nGOVERSION:%s\nBUILDTIME:%s\n", VERSION, COMMIT, BRANCH, GOVERSION, BUILDTIME)
}

func CompareMiddleRange(v, vlo, vhi string) bool {
	sv, _ := strconv.Atoi(strings.Replace(v, ".", "", -1))
	svlo, _ := strconv.Atoi(strings.Replace(vlo, ".", "", -1))
	svhi, _ := strconv.Atoi(strings.Replace(vhi, ".", "", -1))
	fmt.Printf("v:%d, vlo:%d, vhi:%d\n", sv, svlo, svhi)
	return sv >= svlo && sv <= svhi
}

func GreaterThanEq(v, targetV string) bool {
	sv, _ := strconv.Atoi(strings.Replace(v, ".", "", -1))
	stv, _ := strconv.Atoi(strings.Replace(targetV, ".", "", -1))
	fmt.Printf("v:%d, stv:%d\n", sv, stv)
	return sv >= stv
}

func LessThanEq(v, targetV string) bool {
	sv, _ := strconv.Atoi(strings.Replace(v, ".", "", -1))
	stv, _ := strconv.Atoi(strings.Replace(targetV, ".", "", -1))
	fmt.Printf("v:%d, stv:%d\n", sv, stv)
	return sv <= stv
}
