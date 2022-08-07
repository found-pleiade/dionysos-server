package utils

import "errors"

func RemoveUintFromSlice(slice []uint, value uint) ([]uint, error) {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...), nil
		}
	}
	// If the value was not found, return the original slice and throw error
	return slice, errors.New("value not found in slice")
}
