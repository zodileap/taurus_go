package byteutil

import (
	"strconv"
	"strings"
)

func BytesToInt(b []byte) int {
	n, _ := strconv.Atoi(string(b))
	return n
}

func BytesToInt8(b []byte) int8 {
	n, _ := strconv.ParseInt(string(b), 10, 8)
	return int8(n)
}

func BytesToInt16(b []byte) int16 {
	n, _ := strconv.ParseInt(string(b), 10, 16)
	return int16(n)
}

func BytesToInt32(b []byte) int32 {
	n, _ := strconv.ParseInt(string(b), 10, 32)
	return int32(n)
}

func BytesToInt64(b []byte) int64 {
	n, _ := strconv.ParseInt(string(b), 10, 64)
	return n
}

func BytesToUint(b []byte) uint {
	n, _ := strconv.ParseUint(string(b), 10, 0)
	return uint(n)
}

func BytesToIntS(b []byte, sep string) []int {
	s := string(b)
	ss := strings.Split(s, sep)
	ns := make([]int, len(ss))
	for i, v := range ss {
		ns[i], _ = strconv.Atoi(v)
	}
	return ns
}

func BytesToInt8S(b []byte, sep string) []int8 {
	s := string(b)
	ss := strings.Split(s, sep)
	ns := make([]int8, len(ss))
	for i, v := range ss {
		n, _ := strconv.ParseInt(v, 10, 8)
		ns[i] = int8(n)
	}
	return ns
}

func BytesToInt16S(b []byte, sep string) []int16 {
	s := string(b)
	ss := strings.Split(s, sep)
	ns := make([]int16, len(ss))
	for i, v := range ss {
		n, _ := strconv.ParseInt(v, 10, 16)
		ns[i] = int16(n)
	}
	return ns
}

func BytesToInt32S(b []byte, sep string) []int32 {
	s := string(b)
	ss := strings.Split(s, sep)
	ns := make([]int32, len(ss))
	for i, v := range ss {
		n, _ := strconv.ParseInt(v, 10, 32)
		ns[i] = int32(n)
	}
	return ns
}

func BytesToInt64S(b []byte, sep string) []int64 {
	s := string(b)
	ss := strings.Split(s, sep)
	ns := make([]int64, len(ss))
	for i, v := range ss {
		ns[i], _ = strconv.ParseInt(v, 10, 64)
	}
	return ns
}

func BytesToUintS(b []byte, sep string) []uint {
	s := string(b)
	ss := strings.Split(s, sep)
	ns := make([]uint, len(ss))
	for i, v := range ss {
		n, _ := strconv.ParseUint(v, 10, 0)
		ns[i] = uint(n)
	}
	return ns
}

func BytesToUint8(b []byte) uint8 {
	n, _ := strconv.ParseUint(string(b), 10, 8)
	return uint8(n)
}
