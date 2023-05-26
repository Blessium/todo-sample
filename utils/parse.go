package utils

import (
    "strconv"
    "errors"
)

var (
    ErrUnsignedInt="id is not an unsigned integer"
    ErrNegativeId="id cannot be less than 0"
)

func StringToUint(id string) (uint, error) {
    i, err := strconv.Atoi(id)
    if err != nil {
        return 0, errors.New(ErrUnsignedInt)
    }
    if i < 0 {
        return 0, errors.New(ErrNegativeId)
    }
    return uint(i), nil
}
