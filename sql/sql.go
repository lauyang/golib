package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Query(db *sql.DB, sql string, val interface{}) error {
	var tagMap map[string]int
	var tp, tps reflect.Type
	var n, i int
	var err error
	var ret reflect.Value
	// 检测val参数是否为我们所想要的参数
	tp = reflect.TypeOf(val)
	if reflect.Ptr != tp.Kind() {
		return errors.New("is not pointer")
	}

	if reflect.Slice != tp.Elem().Kind() {
		return errors.New("is not slice pointer")
	}

	tp = tp.Elem()
	tps = tp.Elem()
	if reflect.Struct != tps.Kind() {
		return errors.New("is not struct slice pointer")
	}

	tagMap = make(map[string]int)
	n = tps.NumField()
	for i = 0; i < n; i++ {
		tag := tps.Field(i).Tag.Get("sql")
		if len(tag) > 0 {
			tagMap[tag] = i + 1
		}
	}

	// 执行查询
	ret, err = queryAndReflect(db, sql, tagMap, tp)
	if nil != err {
		return err
	}

	// 返回结果
	reflect.ValueOf(val).Elem().Set(ret)

	return nil
}

// 查询并构建返回
func queryAndReflect(db *sql.DB,
	sql string,
	tagMap map[string]int,
	tpSlice reflect.Type) (reflect.Value, error) {
	var ret reflect.Value
	// 执行sql语句
	rows, err := db.Query(sql)
	if nil != err {
		return reflect.Value{}, err
	}

	defer rows.Close()
	// 开始枚举结果
	cols, err := rows.Columns()
	if nil != err {
		return reflect.Value{}, err
	}

	ret = reflect.MakeSlice(tpSlice, 0, 50)
	// 构建接收队列
	scan := make([]interface{}, len(cols))
	row := make([]interface{}, len(cols))
	for r := range row {
		scan[r] = &row[r]
	}

	for rows.Next() {
		feild := reflect.New(tpSlice.Elem()).Elem()
		// 取得结果

		err = rows.Scan(scan...)
		// 开始遍历结果
		for i := 0; i < len(cols); i++ {
			n := tagMap[cols[i]] - 1
			if n < 0 {
				continue
			}

			switch feild.Type().Field(n).Type.Kind() {
			case reflect.Bool:
				if nil != row[i] {
					feild.Field(n).SetBool("false" != string(row[i].([]byte)))
				} else {
					feild.Field(n).SetBool(false)
				}
			case reflect.String:
				if nil != row[i] {
					feild.Field(n).SetString(string(row[i].([]byte)))
				} else {
					feild.Field(n).SetString("")
				}
			case reflect.Float32:
				fallthrough
			case reflect.Float64:
				if nil != row[i] {
					v, e := strconv.ParseFloat(string(row[i].([]byte)), 0)
					if nil == e {
						feild.Field(n).SetFloat(v)
					}
				} else {
					feild.Field(n).SetFloat(0)
				}
			case reflect.Int8:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Int64:
				fallthrough
			case reflect.Int:
				if nil != row[i] {
					byRow, ok := row[i].([]byte)
					if ok {
						v, e := strconv.ParseInt(string(byRow), 10, 64)
						if nil == e {
							feild.Field(n).SetInt(v)
						}
					} else {
						v, e := strconv.ParseInt(fmt.Sprint(row[i]), 10, 64)
						if nil == e {
							feild.Field(n).SetInt(v)
						}
					}
				} else {
					feild.Field(n).SetInt(0)
				}
			case reflect.Uint8:
				fallthrough
			case reflect.Uint16:
				fallthrough
			case reflect.Uint32:
				fallthrough
			case reflect.Uint64:
				fallthrough
			case reflect.Uint:
				if nil != row[i] {
					byRow, ok := row[i].([]byte)
					if ok {
						v, e := strconv.ParseUint(string(byRow), 10, 64)
						if nil == e {
							feild.Field(n).SetUint(v)
						}
					} else {
						v, e := strconv.ParseUint(fmt.Sprint(row[i]), 10, 64)
						if nil == e {
							feild.Field(n).SetUint(v)
						}
					}
				} else {
					feild.Field(n).SetUint(0)
				}
			}
		}

		ret = reflect.Append(ret, feild)
	}

	return ret, nil
}
