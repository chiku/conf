package conf_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/chiku/conf"
)

const (
	json     = "JSON"
	flags    = "Flags"
	env      = "Environment"
	defaults = "Defaults"

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
	loader := &conf.MultiLoader{
		Args:      []string{"-man", manf, "-opt", optf},
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.Load()

	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], manf, "Expected mandatory flags config to be extracted")
	assertEqual(t, config["opt"], optf, "Expected optional flags config to be extracted")
	assertEqual(t, origin["man"], flags, "Expected mandatory config to be provided by flags")
	assertEqual(t, origin["opt"], flags, "Expected optional config to be provided by flags")
}

func TestLoadFromJSON(t *testing.T) {
	content := fmt.Sprintf(`{ "man": "%s", "opt": "%s"	}`, manj, optj)
	jsonFile := createFile(t, content)
	defer os.Remove(jsonFile)

	loader := &conf.MultiLoader{
		JSONKey:   "conf",
		Args:      []string{"-conf", jsonFile},
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.Load()

	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], manj, "Expected mandatory JSON config to be extracted")
	assertEqual(t, config["opt"], optj, "Expected optional JSON config to be extracted")
	assertEqual(t, origin["man"], json, "Expected mandatory config to be provided by JSON")
	assertEqual(t, origin["opt"], json, "Expected optional config to be provided by JSON")
}

func TestLoadFromEnvironment(t *testing.T) {
	err := os.Setenv("man", mane)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("man")
	err = os.Setenv("opt", opte)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("opt")

	loader := &conf.MultiLoader{
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
	}
	config, origin, err := loader.Load()

	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], mane, "Expected mandatory environment config to be extracted")
	assertEqual(t, config["opt"], opte, "Expected optional environment config to be extracted")
	assertEqual(t, origin["man"], env, "Expected mandatory config to be provided by environment")
	assertEqual(t, origin["opt"], env, "Expected optional config to be provided by environment")
}

func TestLoadFromDefaults(t *testing.T) {
	loader := &conf.MultiLoader{
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.Load()

	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], mand, "Expected mandatory defaults config to be extracted")
	assertEqual(t, config["opt"], optd, "Expected optional defaults config to be extracted")
	assertEqual(t, origin["man"], defaults, "Expected mandatory config to be provided by defaults")
	assertEqual(t, origin["opt"], defaults, "Expected optional config to be provided by defaults")
}

func TestLoadFromFlagsHasHighestPriority(t *testing.T) {
	content := fmt.Sprintf(`{ "man": "%s", "opt": "%s"	}`, manj, optj)
	jsonFile := createFile(t, content)
	defer os.Remove(jsonFile)

	err := os.Setenv("man", mane)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("man")
	err = os.Setenv("opt", opte)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("opt")

	loader := &conf.MultiLoader{
		JSONKey:   "conf",
		Args:      []string{"-conf", jsonFile, "-man", manf, "-opt", optf},
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.Load()

	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], manf, "Expected mandatory flags config to be extracted")
	assertEqual(t, config["opt"], optf, "Expected optional flags config to be extracted")
	assertEqual(t, origin["man"], flags, "Expected mandatory config to be provided by flags")
	assertEqual(t, origin["opt"], flags, "Expected optional config to be provided by flags")
}

func TestLoadFromJSONHasPriorityOverEnvironmentAndDefaults(t *testing.T) {
	content := fmt.Sprintf(`{ "man": "%s", "opt": "%s"	}`, manj, optj)
	jsonFile := createFile(t, content)
	defer os.Remove(jsonFile)

	err := os.Setenv("man", mane)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("man")
	err = os.Setenv("opt", opte)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("opt")

	loader := &conf.MultiLoader{
		JSONKey:   "conf",
		Args:      []string{"-conf", jsonFile},
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.Load()

	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], manj, "Expected mandatory JSON config to be extracted")
	assertEqual(t, config["opt"], optj, "Expected optional JSON config to be extracted")
	assertEqual(t, origin["man"], json, "Expected mandatory config to be provided by json")
	assertEqual(t, origin["opt"], json, "Expected optional config to be provided by json")
}

func TestLoadFromEnvironmentHasPriorityOverDefaults(t *testing.T) {
	err := os.Setenv("man", mane)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("man")
	err = os.Setenv("opt", opte)
	requireNoError(t, err, "Expected no error setting environment")
	defer os.Unsetenv("opt")

	loader := &conf.MultiLoader{
		JSONKey:   "conf",
		Mandatory: []string{"man"},
		Optional:  []string{"opt"},
		Defaults:  map[string]string{"man": mand, "opt": optd},
	}
	config, origin, err := loader.Load()

	requireNoError(t, err, "Expected no error loading conf")

	assertEqual(t, config["man"], mane, "Expected mandatory environment config to be extracted")
	assertEqual(t, config["opt"], opte, "Expected optional environment config to be extracted")
	assertEqual(t, origin["man"], env, "Expected mandatory config to be provided by environment")
	assertEqual(t, origin["opt"], env, "Expected optional config to be provided by environment")
}

func requireNoError(t *testing.T, err error, msg string) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		fmt.Printf("\t%v:%v: %s\n\n", fileBase, line, err.Error())
		t.FailNow()
	}
}

func assertEqual(t *testing.T, actual, expected, msg string) {
	if actual != expected {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		fmt.Printf("\t%v:%v: %#v != %#v\n\n", fileBase, line, actual, expected)
		t.Fail()
	}
}

func createFile(t *testing.T, content string) string {
	tmpfile, err := ioutil.TempFile("", "example")
	requireNoError(t, err, "Expected no error creating temporary file")
	_, err = tmpfile.Write([]byte(content))
	requireNoError(t, err, "Expected no error writing to temporary file")
	err = tmpfile.Close()
	requireNoError(t, err, "Expected no error closing temporary file")

	return tmpfile.Name()
}
