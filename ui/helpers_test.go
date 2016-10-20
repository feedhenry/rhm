package ui

import "testing"

func TestPadLeft(t *testing.T) {
	res := PadLeft(8, " ", "test")
	if res != "    test" {
		t.Fatal("PadLeft got: '" + res + "', expected: '    test'")
	}

	res = PadLeft(3, " ", "test")
	if res != "test" {
		t.Fatal("PadLeft got: '" + res + "', expected: 'test'")
	}
}

func TestPadRight(t *testing.T) {
	res := PadRight(8, " ", "test")
	if res != "test    " {
		t.Fatal("PadRight got: '" + res + "', expected: 'test    '")
	}

	res = PadRight(3, " ", "test")
	if res != "test" {
		t.Fatal("PadRight got: '" + res + "', expected: 'test'")
	}
}

func TestPadCentered(t *testing.T) {
	res := PadCentered(8, " ", "test")
	if res != "  test  " {
		t.Fatal("PadCentered got: '" + res + "', expected: '  test  '")
	}
	res = PadCentered(9, " ", "test")
	if res != "   test  " {
		t.Fatal("PadCentered got: '" + res + "', expected: '   test  '")
	}

	res = PadCentered(3, " ", "test")
	if res != "test" {
		t.Fatal("PadCentered got: '" + res + "', expected: 'test'")
	}
}
