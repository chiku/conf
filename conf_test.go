package conf_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/chiku/conf"
)

const (
	json     = "JSON"
	flags    = "Flags"
	env      = "Environment"
	defaults = "Defaults"
)

func TestLoadFromJSON(t *testing.T) {
	const man = "man:json"
	const opt = "opt:json"

	loader := &conf.MultiLoader{
		JSON: fmt.Sprintf(`{ "man": "%s", "opt": "%s"	}`, man, opt),
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.Load()

	if err != nil {
		t.Fatalf("Expected no error loading conf, but got: %v", err)
	}

	if config["man"] != man {
		t.Errorf(`Expected mandatory JSON config to be extracted, but '%v' != '%v'`, config["man"], man)
	}
	if config["opt"] != opt {
		t.Errorf(`Expected optional JSON config to be extracted, but '%v' != '%v'`, config["opt"], opt)
	}

	if origin["man"] != "JSON" {
		t.Errorf(`Expected mandatory config to be provided by JSON, but was provided by: '%v'`, origin["man"])
	}
	if origin["opt"] != "JSON" {
		t.Errorf(`Expected optional config to be provided by JSON, but was provided by: '%v'`, origin["opt"])
	}
}

func TestLoadFromFlags(t *testing.T) {
	const man = "man:flags"
	const opt = "opt:flags"

	loader := &conf.MultiLoader{
		Args:      []string{"-man", man, "-opt", opt},
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.Load()

	if err != nil {
		t.Fatalf("Expected no error loading conf, but got %v", err)
	}

	if config["man"] != man {
		t.Errorf(`Expected mandatory flags config to be extracted, but '%v' != '%v'`, config["man"], man)
	}
	if config["opt"] != opt {
		t.Errorf(`Expected optional flags config to be extracted, but '%v' != '%v'`, config["opt"], opt)
	}

	if origin["man"] != flags {
		t.Errorf(`Expected mandatory config to be provided by flags, but was provided by: '%v'`, origin["man"])
	}
	if origin["opt"] != flags {
		t.Errorf(`Expected optional config to be provided by flags, but was provided by: '%v'`, origin["opt"])
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	const man = "man:env"
	const opt = "opt:env"

	if err := os.Setenv("man", man); err != nil {
		t.Fatalf("Expected no error setting environment, but got: %v", err)
	}
	defer os.Unsetenv("man")
	if err := os.Setenv("opt", opt); err != nil {
		t.Fatalf("Expected no error setting environment, but got: %v", err)
	}
	defer os.Unsetenv("opt")

	loader := &conf.MultiLoader{
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.Load()

	if err != nil {
		t.Fatalf("Expected no error loading conf, but got %v", err)
	}

	if config["man"] != man {
		t.Errorf(`Expected mandatory environment config to be extracted, but '%v' != '%v'`, config["man"], man)
	}
	if config["opt"] != opt {
		t.Errorf(`Expected optional environment config to be extracted, but '%v' != '%v'`, config["opt"], opt)
	}

	if origin["man"] != env {
		t.Errorf(`Expected mandatory config to be provided by environment, but was provided by: '%v'`, origin["man"])
	}
	if origin["opt"] != env {
		t.Errorf(`Expected optional config to be provided by environment, but was provided by: '%v'`, origin["opt"])
	}
}

func TestLoadFromDefaults(t *testing.T) {
	const man = "man:defaults"
	const opt = "opt:defaults"

	loader := &conf.MultiLoader{
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
		Defaults: map[string]string{
			"man": man,
			"opt": opt,
		},
	}
	config, origin, err := loader.Load()

	if err != nil {
		t.Fatalf("Expected no error loading conf, but got %v", err)
	}

	if config["man"] != man {
		t.Errorf(`Expected mandatory defaults config to be extracted, but '%v' != '%v'`, config["man"], man)
	}
	if config["opt"] != opt {
		t.Errorf(`Expected optional defaults config to be extracted, but '%v' != '%v'`, config["opt"], opt)
	}

	if origin["man"] != defaults {
		t.Errorf(`Expected mandatory config to be provided by defaults, but was provided by: '%v'`, origin["man"])
	}
	if origin["opt"] != defaults {
		t.Errorf(`Expected optional config to be provided by defaults, but was provided by: '%v'`, origin["opt"])
	}
}
