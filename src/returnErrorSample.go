package main

import (
	"fmt"
	"math"
)

const (
	Eps = 1e-12
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %d",float64(e))
}

func Sqrt(x float64) (float64, error) {
	if (x < 0) {
		return x,ErrNegativeSqrt(x)
	}
	z := 1.0
	for ; math.Abs(z*z - x) > Eps; {
		z = z - (z*z - x)/(2*z)
	}
	return z,nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
