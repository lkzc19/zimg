package mapx

import (
	"fmt"
	"testing"
)

func TestSliceMap(t *testing.T) {
	m := NewSliceMap[string, string]()

	m.Set("a", "1")
	m.Set("b", "2")
	m.Set("c", "3")
	m.Set("d", "4")
	m.Set("c", "5")

	for _, k := range m.Keys() {
		v, _ := m.Get(k)
		fmt.Println(k, v)
	}
}
