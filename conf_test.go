package conf

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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
	loader := &MultiLoader{
		Mandatory:   []string{"man"},
		Optional:    []string{"opt"},
		Description: map[string]string{"man": "mandatory item"},
	}
	config, origin, err := loader.load([]string{"-man", manf, "-opt", optf}, sampleUsage)

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

	loader := &MultiLoader{
		JSONKey:   "conf",
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.load([]string{"-conf", jsonFile}, sampleUsage)

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

	loader := &MultiLoader{
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.load(nil, sampleUsage)

	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], mane, "Expected mandatory environment config to be extracted")
	assertEqual(t, config["opt"], opte, "Expected optional environment config to be extracted")
	assertEqual(t, origin["man"], envOrig, "Expected mandatory config to be provided by environment")
	assertEqual(t, origin["opt"], envOrig, "Expected optional config to be provided by environment")
}

func TestLoadFromDefaults(t *testing.T) {
	loader := &MultiLoader{
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.load(nil, sampleUsage)

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

	loader := &MultiLoader{
		JSONKey:   "conf",
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.load([]string{"-conf", jsonFile, "-man", manf, "-opt", optf}, sampleUsage)

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

	loader := &MultiLoader{
		JSONKey:   "conf",
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.load([]string{"-conf", jsonFile}, sampleUsage)

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

	loader := &MultiLoader{
		JSONKey:   "conf",
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.load(nil, sampleUsage)

	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], mane, "Expected mandatory environment config to be extracted")
	assertEqual(t, config["opt"], opte, "Expected optional environment config to be extracted")
	assertEqual(t, origin["man"], envOrig, "Expected mandatory config to be provided by environment")
	assertEqual(t, origin["opt"], envOrig, "Expected optional config to be provided by environment")
}

func TestFlagParseError(t *testing.T) {
	loader := &MultiLoader{
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.load([]string{"-many", "-opty"}, sampleUsage)

	requireError(t, err, "Expected error loading conf with bad flags")
	assertContains(t, err.Error(), "conf.Load: error parsing flags: ", "Expected flag parse error message")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestJSONFileReadError(t *testing.T) {
	loader := &MultiLoader{
		JSONKey:   "conf",
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.load([]string{"-conf", "file-does-not-exist"}, sampleUsage)

	requireError(t, err, "Expected error loading conf with non-existing JSON file")
	assertContains(t, err.Error(), "conf.Load: error reading JSON file: ", "Expected JSON file read error message")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestJSONFileParseError(t *testing.T) {
	content := "bad-json"
	jsonFile := createFile(t, content)
	defer os.Remove(jsonFile)

	loader := &MultiLoader{
		JSONKey:   "conf",
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.load([]string{"-conf", jsonFile}, sampleUsage)

	requireError(t, err, "Expected error loading conf with malformed JSON file")
	assertContains(t, err.Error(), "conf.Load: error parsing JSON file: ", "Expected JSON file parse error message")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestMissingMandatoryConfigError(t *testing.T) {
	loader := &MultiLoader{
		Mandatory: []string{"man", "man2", "man3"},
		Optional:  []string{"opt", "opt2", "opt3"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.load(nil, sampleUsage)

	requireError(t, err, "Expected error loading conf with missing mandatory configurations")
	assertEqual(t, err.Error(), "conf.Load: missing mandatory configurations: man2, man3", "Expected missing mandatory configurations")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestFlagKeyCollisionsError(t *testing.T) {
	loader := &MultiLoader{
		JSONKey:   "shr",
		Mandatory: []string{"man", "man", "man1", "man1", "shr1", "shr2", "shr"},
		Optional:  []string{"opt", "opt", "opt1", "opt1", "shr1", "shr2", "shr"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.load(nil, sampleUsage)

	requireError(t, err, "Expected error loading conf with overlapping mandatory and optional configurations")
	assertEqual(t, err.Error(), "conf.Load: configuration keys are duplicated: mandatory(man, man1), optional(opt, opt1), mandatory+optional(shr1, shr2, shr), mandatory+jsonkey(shr), optional+jsonkey(shr)", "Expected overlapping configurations")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestUnknownDescriptionKeyError(t *testing.T) {
	loader := &MultiLoader{
		Mandatory:   []string{"man"},
		Optional:    []string{"opt"},
		Description: map[string]string{"man1": "man1 description", "opt": "opt description", "opt1": "opt1 description"},
		Defaults:    map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.load(nil, sampleUsage)

	requireError(t, err, "Expected error loading conf with unknown description")
	assertContains(t, err.Error(), "conf.Load: description keys are unknown: ", "Expected unknown descriptions")
	assertContains(t, err.Error(), "man1", "Expected unknown descriptions")
	assertContains(t, err.Error(), "opt1", "Expected unknown descriptions")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestEmptyKeyError(t *testing.T) {
	loader := &MultiLoader{
		Mandatory: []string{"man", ""},
		Optional:  []string{"opt", ""},
	}
	config, origin, err := loader.load(nil, sampleUsage)

	requireError(t, err, "Expected error loading conf with empty configurations")
	assertEqual(t, err.Error(), "conf.Load: empty keys exist: mandatory, optional", "Expected empty configuration error")

	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func TestFuzzError1(t *testing.T) {
	loader := &MultiLoader{
		JSONKey:   "Ĺ",
		Mandatory: []string{"5Ò劯YņëHƋ訖玲薯ŀ", "CƱ屼=ðȡ", "E聻阑l", "FŧÒ簠}ZĀi>2鯢鎗觡ǲ", "r刍ĵsJ", "Ĺ", "効谄縫BɈ璻)隽Ld", "固[飳Ɛ茞燂Yi衮ɼO榲\u00adȾ", "鮀ȯIÏ忞"},
		Optional:  []string{"Rl:O", "^4uſǖʈƩʟǑȶªIƙǨ鋜", "e郊Ɔ鏬挋眖筎:ûǽǬ鴜Ȃ", "i莝á沷俜ƦǱ缘Ín痐U"},
	}

	config, origin, err := loader.load(nil, sampleUsage)
	requireError(t, err, "Expected error loading conf with empty configurations")
	assertEqual(t, len(config), 0, "Expected configuration to not exist")
	assertEqual(t, len(origin), 0, "Expected origin to not exist")
}

func sampleUsage(flags *flag.FlagSet) func() {
	return func() {
		flags.SetOutput(ioutil.Discard)
	}
}
