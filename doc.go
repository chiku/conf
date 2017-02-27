// Package conf is for extracting application configuration.
// It uses configuration from command-line arguments, JSON file,
// environment variable or default value.

// example.go

// import (
// 	"fmt"

// 	"github.com/chiku/conf"
// )

// func main() {
// 	options := map[string]conf.Option{
// 		"foo": conf.Option{Desc: "a description for foo", Default: "default foo", Mandatory: true},
// 		"bar": conf.Option{Mandatory: true},
// 		"baz": conf.Option{Desc: "a description for baz"},
// 		"qux": conf.Option{},
// 	}

// 	loader := conf.MultiLoader{
// 		JSONKey: "shr",
// 		Options: options,
// 		Usage:   "Example application",
// 	}

// 	config, origin, err := loader.Load()

// 	if err != nil {
// 		fmt.Printf("error: %s\n", err)
// 		return
// 	}

// 	fmt.Printf("configuration: %#v\n", config)
// 	fmt.Printf("origin: %#v\n", origin)
// }

// usage

// go build -o example example.go
// ./example -foo fooval -bar barval -shr file.json

// doc.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2016-2017. All rights reserved
// License::   MIT

package conf
