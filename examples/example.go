package main

import (
	"fmt"

	"github.com/chiku/conf"
)

func main() {
	loader := conf.MultiLoader{
		JSONKey:     "shr",
		Mandatory:   []string{"foo", "bar"},
		Optional:    []string{"baz", "qux"},
		Defaults:    map[string]string{"foo": "default foo"},
		Description: map[string]string{"foo": "a description for foo", "baz": "a description for baz"},
		Usage:       "Example application",
	}

	config, origin, err := loader.Load()

	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("configuration: %#v\n", config)
	fmt.Printf("origin: %#v\n", origin)
}
