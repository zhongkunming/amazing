package util

import (
	"strconv"
	"strings"
	"unsafe"
)

func TransByte(msg []byte) string {
	msgTemp := *(*string)(unsafe.Pointer(&msg))
	msgStr, _ := strconv.Unquote(strings.Replace(strconv.Quote(msgTemp), `\\u`, `\u`, -1))
	return msgStr
}

func Byte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
