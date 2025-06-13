package utils

import (
	"crypto/rand"
	"math/big"
)

func RandomBool() bool {
	n, err := rand.Int(rand.Reader, big.NewInt(2)) // returns 0 or 1
	if err != nil {
		return false
	}

	if n.Int64() == 1 {
		return true
	}

	return false
}
