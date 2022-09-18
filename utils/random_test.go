package utils

import (
	"fmt"
	"testing"
)

func TestRandomCharacter(t *testing.T) {
	var s string

	s = RandomCharacter(20)
	fmt.Println(s)

	s = RandomNumber(10)
	fmt.Println(s)
}
