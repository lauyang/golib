package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

/*
 *	字节数组转int
 */
func BytesToInt(b []byte) int {
	buf := bytes.NewBuffer(b)
	var x int
	binary.Read(buf, binary.BigEndian, &x)

	return int(x)
}

/*
 *	int转字节数组
 */
func IntToBytes(n int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)

	return buf.Bytes()
}

/*
 *	字节数组转int16
 */
func BytesToInt16(b []byte) int16 {
	buf := bytes.NewBuffer(b)
	var x int16
	binary.Read(buf, binary.BigEndian, &x)

	return int16(x)
}

/*
 *	int16转字节数组
 */
func Int16ToBytes(n int16) []byte {
	x := int16(n)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, x)

	return buf.Bytes()
}

/*
 *	字节数组转int32
 */
func BytesToInt32(b []byte) int32 {
	buf := bytes.NewBuffer(b)
	var x int32
	binary.Read(buf, binary.BigEndian, &x)

	return int32(x)
}

/*
 *	int32转字节数组
 */
func Int32ToBytes(n int32) []byte {
	x := int32(n)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, x)

	return buf.Bytes()
}

/*
 *	字节数组转int64
 */
func BytesToInt64(b []byte) int64 {
	buf := bytes.NewBuffer(b)
	var x int64
	binary.Read(buf, binary.BigEndian, &x)

	return int64(x)
}

/*
 *	int64转字节数组
 */
func Int64ToBytes(n int64) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)

	return buf.Bytes()
}

/*
 *	字节数组转uint64
 */
func BytesToUInt64(b []byte) uint64 {
	buf := bytes.NewBuffer(b)
	var x uint64
	binary.Read(buf, binary.BigEndian, &x)

	return uint64(x)
}

/*
 *	uint64转字节数组
 */
func UInt64ToBytes(n uint64) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)

	return buf.Bytes()
}

/*
 *	字节数组转bool
 */
func BytesToBool(b []byte) bool {
	buf := bytes.NewBuffer(b)
	var x bool
	binary.Read(buf, binary.BigEndian, &x)
	return x
}

/*
 *	bool转字节数组
 */
func BoolToBytes(x bool) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, x)

	return buf.Bytes()
}

/*
 *	struct转字节数组
 * 		i:要转换的对象
 */
func StructToBytes(i interface{}) ([]byte, error) {
	v := reflect.Indirect(reflect.ValueOf(i))
	if v.Kind() != reflect.Struct {
		return nil, errors.New("StructToBytes:data type error")
	}
	bufptr := bytes.NewBuffer(nil)
	err := writeBuf(v, bufptr)
	return bufptr.Bytes(), err

}

func writeBuf(v reflect.Value, bufptr *bytes.Buffer) error {
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Type().Kind() {
		case reflect.Struct:
			err := writeBuf(v.Field(i), bufptr)
			if err != nil {
				return err
			}
		case reflect.Bool:
			boolByte := []byte{0}
			if v.Field(i).Bool() {
				boolByte = []byte{1}
			}
			bufptr.Write(boolByte)
		case reflect.String:
			bufptr.WriteString(v.Field(i).String())
		case reflect.Slice:
			bufptr.Write(v.Field(i).Bytes())
		case reflect.Int:
			binary.Write(bufptr, binary.BigEndian, int32(v.Field(i).Int()))
		case reflect.Uint:
			binary.Write(bufptr, binary.BigEndian, uint32(v.Field(i).Uint()))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			binary.Write(bufptr, binary.BigEndian, v.Field(i).Interface())
		}
	}
	return nil
}

/*
 *	字节数组转struct
 *		i:用于接收的对象地址
 *		index:buf起始索引
 */
func BytesToStruct(buf []byte, i interface{}, index int) error {

	v := reflect.Indirect(reflect.ValueOf(i))
	fmt.Println(buf)
	_, err := readBuf(buf, &v, index)
	return err

}

func readBuf(buf []byte, v *reflect.Value, index int) (n int, err error) {

	elem_length := 0
	for i := 0; i < v.NumField(); i++ {
		elem_length = (int)(unsafe.Sizeof(v.Field(i).Type().Kind()))
		switch v.Field(i).Type().Kind() {
		case reflect.Struct:
			elem_length = 0
			temp_v := reflect.Indirect(reflect.ValueOf(v.Field(i).Interface()))
			_, err := readBuf(buf, &temp_v, index)
			if err != nil {
				return 0, err
			}
		case reflect.String:
			str_len := BytesToInt(buf[index : index+4])
			v.Field(i).SetString(string(buf[index+4 : index+4+str_len]))
		case reflect.Int, reflect.Int32:
			elem_length = 4
			v.Field(i).SetInt(int64(BytesToInt32(buf[index : index+elem_length])))
		case reflect.Int64:
			elem_length = 8
			v.Field(i).SetInt(BytesToInt64(buf[index : index+elem_length]))
		}
		index += elem_length

	}
	return index + 1, nil
}

/*
 *	判断主机序
 */
func isBigEndian() bool {
	var i int32 = 0x12345678
	var b byte = byte(i)
	if b == 0x12 {
		return true
	}

	return false
}

/*
 *	主机序转换成网络序
 */
func Htonl(n int32) []byte {
	buf := Int32ToBytes(n)
	if isBigEndian() {
		return buf
	} else {
		newBuf := make([]byte, 0, 4)
		newBuf = append(newBuf, buf[3], buf[2], buf[1], buf[0])
		return newBuf
	}
}

/*
 *	网络序转换成主机序
 */
func Ntohl(buf []byte) int32 {
	if isBigEndian() {
		return BytesToInt32(buf)
	} else {
		newBuf := make([]byte, 0, 4)
		newBuf = append(newBuf, buf[3], buf[2], buf[1], buf[0])
		return BytesToInt32(newBuf)
	}
}

/*
 *	删除slice中的一项
 */
func RemoveSliceElement(v interface{}, index int) interface{} {
	if reflect.TypeOf(v).Kind() != reflect.Slice {
		panic("wrong type")
	}
	slice := reflect.ValueOf(v)
	if index < 0 || index >= slice.Len() {
		panic("out of bounds")
	}
	prev := slice.Index(index)
	for i := index + 1; i < slice.Len(); i++ {
		value := slice.Index(i)
		prev.Set(value)
		prev = value
	}
	return slice.Slice(0, slice.Len()-1).Interface()
}
