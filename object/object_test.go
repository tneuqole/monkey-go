package object

import "testing"

func TestStringHashKey(t *testing.T) {
	h1 := &String{Value: "Hello World"}
	h2 := &String{Value: "Hello World"}
	d1 := &String{Value: "My name is Bob"}
	d2 := &String{Value: "My name is Bob"}

	if h1.HashKey() != h2.HashKey() {
		t.Errorf("strings with same content have different hash key")
	}

	if d1.HashKey() != d2.HashKey() {
		t.Errorf("strings with same content have different hash key")
	}

	if h1.HashKey() == d2.HashKey() {
		t.Errorf("strings with different content have same hash key")
	}

}

func TestIntHashKey(t *testing.T) {
	h1 := &Integer{Value: 1}
	h2 := &Integer{Value: 1}
	d1 := &Integer{Value: 2}
	d2 := &Integer{Value: 2}

	if h1.HashKey() != h2.HashKey() {
		t.Errorf("ints with same val have different hash key")
	}

	if d1.HashKey() != d2.HashKey() {
		t.Errorf("ints with same val have different hash key")
	}

	if h1.HashKey() == d2.HashKey() {
		t.Errorf("ints with different vals have same hash key")
	}

}

func TestBooleanHashKey(t *testing.T) {
	h1 := &Boolean{Value: true}
	h2 := &Boolean{Value: true}
	d1 := &Boolean{Value: false}
	d2 := &Boolean{Value: false}

	if h1.HashKey() != h2.HashKey() {
		t.Errorf("bools with same val have different hash key")
	}

	if d1.HashKey() != d2.HashKey() {
		t.Errorf("bools with same val have different hash key")
	}

	if h1.HashKey() == d2.HashKey() {
		t.Errorf("bools with different vals have same hash key")
	}

}
