package format

import (
	"testing"
)

func TestFormat(t *testing.T) {
	t.Log(Any("hello"))
	t.Log(Any(3.14159265358979323846264))
	t.Log(Any(make(chan int)))
	t.Log(Any(3 + 4i))
	t.Log(Any(make(map[string]int)))
	t.Log(Any(true))
	t.Log(Any(false))
	t.Log(Any(make([]int, 10)))
	t.Log(Any(func() {}))
}
