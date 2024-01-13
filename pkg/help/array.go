package help

// InArrayInt64
func InArrayInt64(array []int64, n int64) bool {
	for _, v := range array {
		if v == n {
			return true
		}
	}
	return false
}

// InArrayInt
func InArrayInt(array []int, n int) bool {
	for _, v := range array {
		if v == n {
			return true
		}
	}
	return false
}

// Iatosa int array to string array
func Iatosa(array []int) []string {
	var s []string
	for _, v := range array {
		s = append(s, Itoa(v))
	}
	return s
}

// Satoia string array to int array
func Satoia(array []string) []int {
	var i []int
	for _, v := range array {
		i = append(i, Atoi(v))
	}
	return i
}

// IntArrayShift
func Int64ArrayShift(array *[]int64, n int64) {
	*array = append([]int64{n}, (*array)...)
}
