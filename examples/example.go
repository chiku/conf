package main

import "fmt"
import "os"
import "github.com/chiku/conf"

func main() {
	loader := conf.MultiLoader{
		JSONKey:   "shr",
		Mandatory: []string{"foo", "bar"},
		Optional:  []string{"baz", "qux"},
		Args:      os.Args[1:],
		Defaults:  map[string]string{"foo": "default foo"},
	}

	config, origin, err := loader.Load()

	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("configuration: %#v\n", config)
	fmt.Printf("origin: %#v\n", origin)
}
