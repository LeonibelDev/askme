package functions

import (
	"math/rand"
	"time"
)

func RandomNumber() int {
	rand.Seed(time.Now().Unix())
	min, max := 11111, 99999
	RandomInRange := rand.Intn(max-min+1) + min
	return RandomInRange
}
