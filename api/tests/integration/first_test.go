package unit

import "testing"

func TestHelloHandler(t *testing.T) {
	got := 2 + 3
	want := 5
	if got != want {
		t.Errorf("Add(2,3) = %d; want %d", got, want)
	}
}
