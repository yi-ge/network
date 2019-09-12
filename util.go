package network

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
)

// ConvertToString .
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// SmartPrint .
func SmartPrint(i interface{}) {
	var kv = make(map[string]interface{})
	vValue := reflect.ValueOf(i)
	vType := reflect.TypeOf(i)
	for i := 0; i < vValue.NumField(); i++ {
		kv[vType.Field(i).Name] = vValue.Field(i)
	}
	fmt.Println("\n---SmartPrint---")
	for k, v := range kv {
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
	}
}

func hex2dot(value string) string {
	var err error

	value = strings.Replace(value, "0x", "", -1)
	ipInt := []int64{0, 0, 0, 0}
	ipStr := []string{"", "", "", ""}

	ipInt[0], err = strconv.ParseInt(value[0:2], 16, 64)
	ipStr[0] = strconv.FormatInt(ipInt[0], 10)

	ipInt[1], err = strconv.ParseInt(value[2:4], 16, 64)
	ipStr[1] = strconv.FormatInt(ipInt[1], 10)

	ipInt[2], err = strconv.ParseInt(value[4:6], 16, 64)
	ipStr[2] = strconv.FormatInt(ipInt[2], 10)

	ipInt[3], err = strconv.ParseInt(value[6:8], 16, 64)
	ipStr[3] = strconv.FormatInt(ipInt[3], 10)

	if err != nil {
		return ""
	}

	return strings.Join(ipStr, ".")
}
