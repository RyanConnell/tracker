package templates

import (
	"fmt"
)

func Mod(i, j int) int {
	return i % j
}

func DoubleDigits(n int) string {
	return fmt.Sprintf("%02d", n)
}
