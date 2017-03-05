package conf

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
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
	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], manf, "Expected mandatory flags config to be extracted")
	assertEqual(t, config["opt"], optf, "Expected optional flags config to be extracted")
	assertEqual(t, origin["man"], flagsOrig, "Expected mandatory config to be provided by flags")
	assertEqual(t, origin["opt"], flagsOrig, "Expected optional config to be provided by flags")
}

func TestLoadFromJSON(t *testing.T) {
	content := fmt.Sprintf(`{ "man": "%s", "opt": "%s" }`, manj, optj)
	jsonFile := createFile(t, content)
	defer os.Remove(jsonFile)

	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{
		Options: options,
		JSONKey: "conf",
	}

	config, origin, err := loader.load([]string{"-conf", jsonFile}, sampleFlagsHandler)
	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], manj, "Expected mandatory JSON config to be extracted")
	assertEqual(t, config["opt"], optj, "Expected optional JSON config to be extracted")
	assertEqual(t, origin["man"], jsonOrig, "Expected mandatory config to be provided by JSON")
	assertEqual(t, origin["opt"], jsonOrig, "Expected optional config to be provided by JSON")
}

func TestLoadFromEnvironment(t *testing.T) {
	err := os.Setenv("man", mane)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("man")
	err = os.Setenv("opt", opte)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("opt")

	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{Options: options}

	config, origin, err := loader.load(nil, sampleFlagsHandler)
	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], mane, "Expected mandatory environment config to be extracted")
	assertEqual(t, config["opt"], opte, "Expected optional environment config to be extracted")
	assertEqual(t, origin["man"], envOrig, "Expected mandatory config to be provided by environment")
	assertEqual(t, origin["opt"], envOrig, "Expected optional config to be provided by environment")
}

func TestLoadFromDefaults(t *testing.T) {
	options := map[string]Option{
		"man": Option{Default: mand, Mandatory: true},
		"opt": Option{Default: optd},
	}
	loader := &MultiLoader{Options: options}

	config, origin, err := loader.load(nil, sampleFlagsHandler)
	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], mand, "Expected mandatory defaults config to be extracted")
	assertEqual(t, config["opt"], optd, "Expected optional defaults config to be extracted")
	assertEqual(t, origin["man"], defaultsOrig, "Expected mandatory config to be provided by defaults")
	assertEqual(t, origin["opt"], defaultsOrig, "Expected optional config to be provided by defaults")
}

func TestLoadFromFlagsHasHighestPriority(t *testing.T) {
	content := fmt.Sprintf(`{ "man": "%s", "opt": "%s" }`, manj, optj)
	jsonFile := createFile(t, content)
	defer os.Remove(jsonFile)

	err := os.Setenv("man", mane)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("man")
	err = os.Setenv("opt", opte)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("opt")

	options := map[string]Option{
		"man": Option{Default: mand, Mandatory: true},
		"opt": Option{Default: optd},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load([]string{"-conf", jsonFile, "-man", manf, "-opt", optf}, sampleFlagsHandler)
	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], manf, "Expected mandatory flags config to be extracted")
	assertEqual(t, config["opt"], optf, "Expected optional flags config to be extracted")
	assertEqual(t, origin["man"], flagsOrig, "Expected mandatory config to be provided by flags")
	assertEqual(t, origin["opt"], flagsOrig, "Expected optional config to be provided by flags")
}

func TestLoadFromJSONHasPriorityOverEnvironmentAndDefaults(t *testing.T) {
	content := fmt.Sprintf(`{ "man": "%s", "opt": "%s" }`, manj, optj)
	jsonFile := createFile(t, content)
	defer os.Remove(jsonFile)

	err := os.Setenv("man", mane)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("man")
	err = os.Setenv("opt", opte)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("opt")

	options := map[string]Option{
		"man": Option{Default: mand, Mandatory: true},
		"opt": Option{Default: optd},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load([]string{"-conf", jsonFile}, sampleFlagsHandler)
	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], manj, "Expected mandatory JSON config to be extracted")
	assertEqual(t, config["opt"], optj, "Expected optional JSON config to be extracted")
	assertEqual(t, origin["man"], jsonOrig, "Expected mandatory config to be provided by json")
	assertEqual(t, origin["opt"], jsonOrig, "Expected optional config to be provided by json")
}

func TestLoadFromEnvironmentHasPriorityOverDefaults(t *testing.T) {
	err := os.Setenv("man", mane)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("man")
	err = os.Setenv("opt", opte)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("opt")

	options := map[string]Option{
		"man": Option{Default: mand, Mandatory: true},
		"opt": Option{Default: optd},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load(nil, sampleFlagsHandler)
	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], mane, "Expected mandatory environment config to be extracted")
	assertEqual(t, config["opt"], opte, "Expected optional environment config to be extracted")
	assertEqual(t, origin["man"], envOrig, "Expected mandatory config to be provided by environment")
	assertEqual(t, origin["opt"], envOrig, "Expected optional config to be provided by environment")
}

func TestFlagParseError(t *testing.T) {
	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{Options: options}

	config, origin, err := loader.load([]string{"-many", "-opty"}, sampleFlagsHandler)
	requireError(t, err, "Expected error loading conf with bad flags")
	assertContains(t, err.Error(), "conf.Load: error parsing flags: ", "Expected flag parse error message")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestJSONFileReadError(t *testing.T) {
	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load([]string{"-conf", "file-does-not-exist"}, sampleFlagsHandler)
	requireError(t, err, "Expected error loading conf with non-existing JSON file")
	assertContains(t, err.Error(), "conf.Load: error reading JSON file: ", "Expected JSON file read error message")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestJSONFileParseError(t *testing.T) {
	content := "bad-json"
	jsonFile := createFile(t, content)
	defer os.Remove(jsonFile)

	options := map[string]Option{
		"man": Option{Mandatory: true},
		"opt": Option{},
	}
	loader := &MultiLoader{Options: options, JSONKey: "conf"}

	config, origin, err := loader.load([]string{"-conf", jsonFile}, sampleFlagsHandler)
	requireError(t, err, "Expected error loading conf with malformed JSON file")
	assertContains(t, err.Error(), "conf.Load: json: ", "Expected JSON file parse error message")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
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
	requireError(t, err, "Expected error loading conf with missing mandatory configurations")
	assertEqual(t, err.Error(), "conf.Load: missing mandatory configurations: man2, man3", "Expected missing mandatory configurations")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestLoaderInterface(t *testing.T) {
	interfaceType := reflect.TypeOf((*Loader)(nil)).Elem()
	implements := reflect.TypeOf(&MultiLoader{}).Implements(interfaceType)
	assertEqual(t, implements, true, "Expected MultiLoader to be a Loader")
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
	requireError(t, err, "Expected error loading conf with empty configurations")
	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func sampleFlagsHandler(flags *flag.FlagSet) {
	flags.SetOutput(ioutil.Discard)
}
