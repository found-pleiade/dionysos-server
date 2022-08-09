package utils

import (
	"errors"

	"github.com/lib/pq"
)

func RemoveUintFromSlice(slice pq.Int64Array, value uint) (pq.Int64Array, error) {
	for i, v := range slice {
		if v == int64(value) {
			return append(slice[:i], slice[i+1:]...), nil
		}
	}
	// If the value was not found, return the original slice and throw error
	return slice, errors.New("value not found in slice")
}
