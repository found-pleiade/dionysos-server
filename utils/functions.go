package utils

func RemoveUintFromSlice(s []uint, i uint) []uint {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
