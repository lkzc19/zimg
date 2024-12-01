package config

import (
	"fmt"
	"testing"
)

func TestLongConfig(t *testing.T) {
	Load()
	for _, k := range Zimgrc.Keys() {
		v, _ := Get(k)
		fmt.Println(k, "\t", v)
	}
}

func TestFlush(t *testing.T) {
	Load()
	for _, k := range Zimgrc.Keys() {
		v, _ := Get(k)
		fmt.Println(k, "\t", v)
	}
	Set("github.bucket", "lkzc19")
	Flush()
}
