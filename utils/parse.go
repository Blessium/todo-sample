package utils

import (
    "strconv"
    "errors"
)

var (
    IdParseError="id is not an unsigned integer"
    IdNegativeError="id cannot be negative"
)

func StringToUint(id string) (uint, error) {
    i, err := strconv.Atoi(id)
    if err != nil {
        return 0, errors.New(IdParseError)
    }
    if i < 0 {
        return 0, errors.New(IdNegativeError)
    }
    return uint(i), nil
}
