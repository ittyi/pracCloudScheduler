package main

import "testing"

func AddHelper(x, y int) int {
	return x + y
}

func TestAddHelper(t *testing.T) {
	AssertEquals(t, 10, Add(1, 1))
	AssertEquals(t, 10, Add(1, 2))
	AssertEquals(t, 10, Add(1, 3))
}

func AssertEqualsHelper(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Unexpected int\nexpected:%d, actual:%d", expected, actual)
	}
}
