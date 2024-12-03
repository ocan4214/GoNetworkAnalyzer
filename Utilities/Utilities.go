package CommonUtils

import (
	"fmt"
	"strings"
	"time"
	"unsafe"
)

func GetCurrentDate() string {
	d := time.Now()

	return fmt.Sprintln("%v/%v/%v %v:%v:%v", d.Day(), d.Month(), d.Year(), d.Hour(), d.Minute(), d.Second())
}

func DecodePort(port uint32) string {
	res := ""
	if IsLittleEndian() {
		res = strings.Join(ConvertU16ToLittleEndianStringSlice(port), "")
	} else {
		res = strings.Join(ConvertU16ToBigEndianStringSlice(port), "")
	}
	return res
}

func DecodeIP(IP uint32) string {
	res := ""
	if IsLittleEndian() {
		res = strings.Join(ConvertU32ToLittleEndianStringSlice(IP), ".")
	} else {
		res = strings.Join(ConvertU32ToBigEndianStringSlice(IP), ".")
	}
	return res
}

func ConvertU16ToLittleEndianStringSlice(addr uint32) []string {
	strslice := make([]string, 0)

	strslice = append(strslice, fmt.Sprintf("%d", addr>>8&255))
	strslice = append(strslice, fmt.Sprintf("%d", addr&255))

	return strslice
}

func ConvertU32ToLittleEndianStringSlice(addr uint32) []string {
	strslice := make([]string, 0)

	strslice = append(strslice, fmt.Sprintf("%d", addr>>24&255))
	strslice = append(strslice, fmt.Sprintf("%d", addr>>16&255))
	strslice = append(strslice, fmt.Sprintf("%d", addr>>8&255))
	strslice = append(strslice, fmt.Sprintf("%d", addr&255))

	return strslice
}

func ConvertU32ToBigEndianStringSlice(addr uint32) []string {
	return ReverseStringSlice(ConvertU16ToLittleEndianStringSlice(addr))
}

func ReverseStringSlice(ss []string) []string {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
	return ss
}

func ConvertU16ToBigEndianStringSlice(addr uint32) []string {
	return ReverseStringSlice(ConvertU16ToLittleEndianStringSlice(addr))
}

func IsLittleEndian() bool {
	var i int = 0xABCD
	ptr := unsafe.Pointer(&i)

	fmt.Printf("%v", *(*byte)(ptr))
	if *(*byte)(ptr) == 0xCD {
		return true
	} else {
		return false

	}

}
