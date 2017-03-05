package conf

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	jsonOrig     = "JSON"
	flagsOrig    = "Flags"
	envOrig      = "Environment"
	defaultsOrig = "Defaults"

	manf = "man:flags"
	optf = "opt:flags"
	manj = "man:json"
	optj = "opt:json"
	mane = "man:env"
	opte = "opt:env"
	mand = "man:defaults"
	optd = "opt:defaults"
)

func TestLoadFromFlags(t *testing.T) {
	options := map[string]Option{
		"man": Option{Desc: "mandatory item", Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{Options: options}

	config, origin, err := loader.load([]string{"-man", manf, "-opt", optf}, sampleFlagsHandler)
	if err != nil {
		t.Fatalf("Unexpected error loading configurations from flags: %s", err)
	}

	expectedConfig := map[string]string{"man": manf, "opt": optf}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Error("Configurations don't match when loaded from flags")
		t.Errorf("\nActual  : %#v", config)
		t.Errorf("\nExpected: %#v", expectedConfig)
	}

	expectedOrigin := map[string]string{"man": flagsOrig, "opt": flagsOrig}
	if !reflect.DeepEqual(origin, expectedOrigin) {
		t.Error("Origins don't match when loaded from flags")
		t.Errorf("\nActual  : %#v", origin)
		t.Errorf("\nExpected: %#v", expectedOrigin)
	}
}

func TestLoadFromJSON(t *testing.T) {
	content := fmt.Sprintf(`{ "man": "%s", "opt": "%s" }`, manj, optj)
	jsonFile := createFile(t, content)
	defer func() {
		if err := os.Remove(jsonFile); err != nil {
			t.Fatalf("Unexpected error deleting temporary file: %s", err)
		}
	}()

	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{
		Options: options,
		JSONKey: "conf",
	}

	config, origin, err := loader.load([]string{"-conf", jsonFile}, sampleFlagsHandler)
	if err != nil {
		t.Fatalf("Unexpected error loading configurations from JSON file: %s", err)
	}

	expectedConfig := map[string]string{"man": manj, "opt": optj}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Error("Configurations don't match when loaded from JSON file")
		t.Errorf("\nActual  : %#v", config)
		t.Errorf("\nExpected: %#v", expectedConfig)
	}

	expectedOrigin := map[string]string{"man": jsonOrig, "opt": jsonOrig}
	if !reflect.DeepEqual(origin, expectedOrigin) {
		t.Error("Origins don't match when loaded from JSON file")
		t.Errorf("\nActual  : %#v", origin)
		t.Errorf("\nExpected: %#v", expectedOrigin)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	if err := os.Setenv("man", mane); err != nil {
		t.Fatalf("Unexpected error setting environment variable: %s", err)
	}
	defer func() {
		if err := os.Unsetenv("man"); err != nil {
			t.Fatalf("Unexpected error unsetting environment variable: %s", err)
		}
	}()

	if err := os.Setenv("opt", opte); err != nil {
		t.Fatalf("Unexpected error setting environment variable: %s", err)
	}
	defer func() {
		if err := os.Unsetenv("opt"); err != nil {
			t.Fatalf("Unexpected error unsetting environment variable: %s", err)
		}
	}()

	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{Options: options}

	config, origin, err := loader.load(nil, sampleFlagsHandler)
	if err != nil {
		t.Fatalf("Unexpected error loading configurations from environment variables: %s", err)
	}

	expectedConfig := map[string]string{"man": mane, "opt": opte}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Error("Configurations don't match when loaded from environment variables")
		t.Errorf("\nActual  : %#v", config)
		t.Errorf("\nExpected: %#v", expectedConfig)
	}

	expectedOrigin := map[string]string{"man": envOrig, "opt": envOrig}
	if !reflect.DeepEqual(origin, expectedOrigin) {
		t.Error("Origins don't match when loaded from environment variables")
		t.Errorf("\nActual  : %#v", origin)
		t.Errorf("\nExpected: %#v", expectedOrigin)
	}
}

func TestLoadFromDefaults(t *testing.T) {
	options := map[string]Option{
		"man": Option{Default: mand, Mandatory: true},
		"opt": Option{Default: optd},
	}
	loader := &MultiLoader{Options: options}

	config, origin, err := loader.load(nil, sampleFlagsHandler)
	if err != nil {
		t.Fatalf("Unexpected error loading configurations from defaults: %s", err)
	}

	expectedConfig := map[string]string{"man": mand, "opt": optd}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Error("Configurations don't match when loaded from defaults")
		t.Errorf("\nActual  : %#v", config)
		t.Errorf("\nExpected: %#v", expectedConfig)
	}

	expectedOrigin := map[string]string{"man": defaultsOrig, "opt": defaultsOrig}
	if !reflect.DeepEqual(origin, expectedOrigin) {
		t.Error("Origins don't match when loaded from defaults")
		t.Errorf("\nActual  : %#v", origin)
		t.Errorf("\nExpected: %#v", expectedOrigin)
	}
}

func TestLoadFromFlagsHasHighestPriority(t *testing.T) {
	content := fmt.Sprintf(`{ "man": "%s", "opt": "%s" }`, manj, optj)
	jsonFile := createFile(t, content)
	defer func() {
		if err := os.Remove(jsonFile); err != nil {
			t.Fatalf("Unexpected error deleting temporary file: %s", err)
		}
	}()

	if err := os.Setenv("man", mane); err != nil {
		t.Fatalf("Unexpected error setting environment variable: %s", err)
	}
	defer func() {
		if err := os.Unsetenv("man"); err != nil {
			t.Fatalf("Unexpected error unsetting environment variable: %s", err)
		}
	}()

	if err := os.Setenv("opt", opte); err != nil {
		t.Fatalf("Unexpected error setting environment variable: %s", err)
	}
	defer func() {
		if err := os.Unsetenv("opt"); err != nil {
			t.Fatalf("Unexpected error unsetting environment variable: %s", err)
		}
	}()

	options := map[string]Option{
		"man": Option{Default: mand, Mandatory: true},
		"opt": Option{Default: optd},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load([]string{"-conf", jsonFile, "-man", manf, "-opt", optf}, sampleFlagsHandler)
	if err != nil {
		t.Fatalf("Unexpected error loading configurations from flags, JSON, environment variable and defaults: %s", err)
	}

	expectedConfig := map[string]string{"man": manf, "opt": optf}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Error("Configurations don't match when loaded from flags, JSON, environment variable and defaults")
		t.Errorf("\nActual  : %#v", config)
		t.Errorf("\nExpected: %#v", expectedConfig)
	}

	expectedOrigin := map[string]string{"man": flagsOrig, "opt": flagsOrig}
	if !reflect.DeepEqual(origin, expectedOrigin) {
		t.Error("Origins don't match when loaded from flags, JSON, environment variable and defaults")
		t.Errorf("\nActual  : %#v", origin)
		t.Errorf("\nExpected: %#v", expectedOrigin)
	}
}

func TestLoadFromJSONHasPriorityOverEnvironmentAndDefaults(t *testing.T) {
	content := fmt.Sprintf(`{ "man": "%s", "opt": "%s" }`, manj, optj)
	jsonFile := createFile(t, content)
	defer func() {
		if err := os.Remove(jsonFile); err != nil {
			t.Fatalf("Unexpected error deleting temporary file: %s", err)
		}
	}()

	if err := os.Setenv("man", mane); err != nil {
		t.Fatalf("Unexpected error setting environment variable: %s", err)
	}
	defer func() {
		if err := os.Unsetenv("man"); err != nil {
			t.Fatalf("Unexpected error unsetting environment variable: %s", err)
		}
	}()

	if err := os.Setenv("opt", opte); err != nil {
		t.Fatalf("Unexpected error setting environment variable: %s", err)
	}
	defer func() {
		if err := os.Unsetenv("opt"); err != nil {
			t.Fatalf("Unexpected error unsetting environment variable: %s", err)
		}
	}()

	options := map[string]Option{
		"man": Option{Default: mand, Mandatory: true},
		"opt": Option{Default: optd},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load([]string{"-conf", jsonFile}, sampleFlagsHandler)
	if err != nil {
		t.Fatalf("Unexpected error loading configurations from JSON, environment variable and defaults: %s", err)
	}

	expectedConfig := map[string]string{"man": manj, "opt": optj}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Error("Configurations don't match when loaded from JSON, environment variable and defaults")
		t.Errorf("\nActual  : %#v", config)
		t.Errorf("\nExpected: %#v", expectedConfig)
	}

	expectedOrigin := map[string]string{"man": jsonOrig, "opt": jsonOrig}
	if !reflect.DeepEqual(origin, expectedOrigin) {
		t.Error("Origins don't match when loaded from JSON, environment variable and defaults")
		t.Errorf("\nActual  : %#v", origin)
		t.Errorf("\nExpected: %#v", expectedOrigin)
	}
}

func TestLoadFromEnvironmentHasPriorityOverDefaults(t *testing.T) {
	if err := os.Setenv("man", mane); err != nil {
		t.Fatalf("Unexpected error setting environment variable: %s", err)
	}
	defer func() {
		if err := os.Unsetenv("man"); err != nil {
			t.Fatalf("Unexpected error unsetting environment variable: %s", err)
		}
	}()

	if err := os.Setenv("opt", opte); err != nil {
		t.Fatalf("Unexpected error setting environment variable: %s", err)
	}
	defer func() {
		if err := os.Unsetenv("opt"); err != nil {
			t.Fatalf("Unexpected error unsetting environment variable: %s", err)
		}
	}()

	options := map[string]Option{
		"man": Option{Default: mand, Mandatory: true},
		"opt": Option{Default: optd},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load(nil, sampleFlagsHandler)
	if err != nil {
		t.Fatalf("Unexpected error loading configurations from environment variable and defaults: %s", err)
	}

	expectedConfig := map[string]string{"man": mane, "opt": opte}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Error("Configurations don't match when loaded from environment variable and defaults")
		t.Errorf("\nActual  : %#v", config)
		t.Errorf("\nExpected: %#v", expectedConfig)
	}

	expectedOrigin := map[string]string{"man": envOrig, "opt": envOrig}
	if !reflect.DeepEqual(origin, expectedOrigin) {
		t.Error("Origins don't match when loaded from environment variable and defaults")
		t.Errorf("\nActual  : %#v", origin)
		t.Errorf("\nExpected: %#v", expectedOrigin)
	}
}

func TestFlagParseError(t *testing.T) {
	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{Options: options}

	config, origin, err := loader.load([]string{"-many", "-opty"}, sampleFlagsHandler)
	if expectedMsg := "conf.Load: error parsing flags: "; !strings.Contains(err.Error(), expectedMsg) {
		t.Error("Invalid error message for improper flags")
		t.Errorf("Actual       : %q", err)
		t.Errorf("Expected part: %q", expectedMsg)
	}

	if len(config) != 0 || len(origin) != 0 {
		t.Error("Unexpected invalid values for improper flags")
		t.Errorf("Config: %#v", config)
		t.Errorf("Origin: %#v", origin)
	}
}

func TestJSONFileReadError(t *testing.T) {
	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load([]string{"-conf", "file-does-not-exist"}, sampleFlagsHandler)
	if expectedMsg := "conf.Load: error reading JSON file: "; !strings.Contains(err.Error(), expectedMsg) {
		t.Error("Invalid error message for missing JSON file")
		t.Errorf("Actual       : %q", err)
		t.Errorf("Expected part: %q", expectedMsg)
	}

	if len(config) != 0 || len(origin) != 0 {
		t.Error("Unexpected invalid values for missing JSON file")
		t.Errorf("Config: %#v", config)
		t.Errorf("Origin: %#v", origin)
	}
}

func TestJSONFileParseError(t *testing.T) {
	content := "bad-json"
	jsonFile := createFile(t, content)
	defer func() {
		if err := os.Remove(jsonFile); err != nil {
			t.Fatalf("Unexpected error deleting temporary file: %s", err)
		}
	}()

	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load([]string{"-conf", jsonFile}, sampleFlagsHandler)
	if expectedMsg := "conf.Load: json: "; !strings.Contains(err.Error(), expectedMsg) {
		t.Error("Invalid error message for malformed JSON file")
		t.Errorf("Actual       : %q", err)
		t.Errorf("Expected part: %q", expectedMsg)
	}

	if len(config) != 0 || len(origin) != 0 {
		t.Error("Unexpected invalid values for malformed JSON file")
		t.Errorf("Config: %#v", config)
		t.Errorf("Origin: %#v", origin)
	}
}

func TestMissingMandatoryConfigError(t *testing.T) {
	options := map[string]Option{
		"man":  Option{Default: mand, Mandatory: true},
		"man2": Option{Mandatory: true},
		"man3": Option{Mandatory: true},
		"opt":  Option{Default: optd},
		"opt2": Option{},
		"opt3": Option{},
	}
	loader := &MultiLoader{Options: options}

	config, origin, err := loader.load(nil, sampleFlagsHandler)
	if expectedMsg := "conf.Load: missing mandatory configurations: man2, man3"; err.Error() != expectedMsg {
		t.Error("Invalid error message for missing mandatory configurations")
		t.Errorf("Actual  : %q", err)
		t.Errorf("Expected: %q", expectedMsg)
	}

	if len(config) != 0 || len(origin) != 0 {
		t.Error("Unexpected invalid values for missing mandatory configurations")
		t.Errorf("Config: %#v", config)
		t.Errorf("Origin: %#v", origin)
	}
}

func TestLoaderInterface(t *testing.T) {
	interfaceType := reflect.TypeOf((*Loader)(nil)).Elem()
	implements := reflect.TypeOf(&MultiLoader{}).Implements(interfaceType)
	if !implements {
		t.Error("MultiLoader does not implement Loader")
	}
}

func TestFuzzError1(t *testing.T) {
	options := map[string]Option{
		"5Ò劯YņëHƋ訖玲薯ŀ":        Option{Mandatory: true},
		"CƱ屼=ðȡ":              Option{Mandatory: true},
		"E聻阑l":                Option{Mandatory: true},
		"FŧÒ簠}ZĀi>2鯢鎗觡ǲ":      Option{Mandatory: true},
		"r刍ĵsJ":               Option{Mandatory: true},
		"効谄縫BɈ璻)隽Ld":          Option{Mandatory: true},
		"固[飳Ɛ茞燂Yi衮ɼO榲\u00adȾ": Option{Mandatory: true},
		"鮀ȯIÏ忞":               Option{Mandatory: true},
		"Rl:O":                Option{},
		"^4uſǖʈƩʟǑȶªIƙǨ鋜": Option{},
		"e郊Ɔ鏬挋眖筎:ûǽǬ鴜Ȃ":   Option{},
		"i莝á沷俜ƦǱ缘Ín痐U":    Option{},
	}
	loader := &MultiLoader{Options: options, JSONKey: "Ĺ"}

	config, origin, err := loader.load(nil, sampleFlagsHandler)

	if err == nil {
		t.Error("Unexpected success for invalid configuration")
	}

	if len(config) != 0 || len(origin) != 0 {
		t.Error("Unexpected invalid values for invalid configuration")
		t.Errorf("Config: %#v", config)
		t.Errorf("Origin: %#v", origin)
	}
}

func sampleFlagsHandler(flags *flag.FlagSet) {
	flags.SetOutput(ioutil.Discard)
}
