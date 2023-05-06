package algorism

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	print(add(5, 3))
}

func TestSub(t *testing.T) {
	print(sub(5, 3))
}

func TestMulti(t *testing.T) {
	print(multi(5, 3))
}

func TestDiv(t *testing.T) {
	quotient, remainder := div(321, 123)
	fmt.Printf("%d, %d", quotient, remainder)
}