package format

import (
	"regexp"
	"testing"
)

func TestFormat(t *testing.T) {
	c := []struct {
		V interface{}
		S string
	}{
		{3.14159265358979323846264, "3.141592653589793"},
		{"hello", "^\"hello\"$"},
		{make(chan int), "^chan int 0x\\w+$"},
		{make(map[string]int), "map\\[string\\]int 0x\\w+$"},
		{true, "^true$"},
		{false, "^false$"},
		{struct{}{}, "^struct \\{\\} value$"},
		{struct {
			A int
			B int
		}{}, "^struct \\{ A int; B int \\} value$"},
		{func() {}, "^func\\(\\) 0x\\w+$"},
	}

	for _, item := range c {
		s := Any(item.V)
		if b, _ := regexp.Match(item.S, []byte(s)); b != true {
			t.Fatalf("format error: regexp %s does not match %s\n", item.S, s)
		} else {
			t.Log(s)
		}
	}
}

func TestDisplay(t *testing.T) {
	Display("S", struct {
		A int
		B int
	}{A: 9, B: 22})

	Display("A", [3]int{100, 200, 300})
	Display("MSI", map[string]int{"A": 1, "B": 2, "C": 3})
	Display("MAI", map[[3]int]int{{1, 2, 3}: 1, {2, 3, 4}: 2, {4, 5, 6}: 3})
}
