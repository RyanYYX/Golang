package algorism

import "testing"

func TestBF(t *testing.T) {
	print(BF("hello world", "world"))
}

func TestRK(t *testing.T) {
	print(RK("wkldjdhvsdjkfhcinuqbhw", "jdhvsdjkfhcinuqbhw"))
}

func TestKMP(t *testing.T) {
	print(KMP("abababcabcababcabcd", "ababcabcd"))
}
