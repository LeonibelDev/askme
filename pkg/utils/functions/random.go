package functions

import (
	"crypto/rand"
	"math/big"
)

func RandomNumber() (int, error) {
	min, max := 00001, 99999
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}
	return int(nBig.Int64()), nil
}
