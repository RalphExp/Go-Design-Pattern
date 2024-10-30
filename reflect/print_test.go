package format

import (
	"os"
	"reflect"
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
		{struct{}{}, "^\\{\\}$"},
		{struct {
			A int
			B int
		}{}, "^\\{A:0, B:0\\}$"},
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
	Display("Struct", struct {
		A int
		B int
	}{A: 9, B: 22})

	Display("Array", [3]int{100, 200, 300})
	Display("MapString", map[string]int{"A": 1, "B": 2, "C": 3})
	Display("MapArray", map[[3]int]int{{1, 2, 3}: 1, {2, 3, 4}: 2, {4, 5, 6}: 3})
	Display("MapArray", map[struct {
		A int
		B int
	}]int{{A: 2, B: 3}: 0})

	var i interface{} = 3
	Display("i", i)
	Display("&i", &i)

	type Movie struct {
		Title, Subtitle string
		Year            int
		Color           bool
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}
	Display("strangelove", strangelove)
	Display("os.Stderr", os.Stderr)
	Display("rV", reflect.ValueOf(os.Stderr))
}
