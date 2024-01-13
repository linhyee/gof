package help

import (
	"strconv"
)

// Atoi string to int
func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// Atoi64 string to int64
func Atoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 0)
	return i
}

// ParseInt
func ParseInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 0)
	return i
}

// ParseUint
func ParseUint(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 0)
	return i
}

// Itoa int to string
func Itoa(i int) string {
	return strconv.Itoa(i)
}

// I64toa int64 to string
func I64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}

// FormatInt
func FormatInt(i int64) string {
	return strconv.FormatInt(i, 10)
}

// FormatUint
func FormatUint(i uint64) string {
	return strconv.FormatUint(i, 10)
}
