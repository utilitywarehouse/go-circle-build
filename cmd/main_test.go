package main

import "testing"

func TestMessage(t *testing.T) {
	mf := Message()
	m := "Hello Circle"
	if mf != m {
		t.Errorf("Expected %v, recieved %v instead", m, mf)
	}
}
