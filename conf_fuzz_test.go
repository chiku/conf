// +build fuzz

package conf

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strconv"
	"testing"

	"github.com/google/gofuzz"
)

func TestConfFuzzRandom(t *testing.T) {
	maxSteps := 15000000
	if maxStensEnv, err := strconv.Atoi(os.Getenv("MAX_STEPS")); err == nil && maxStensEnv != 0 {
		maxSteps = maxStensEnv
	}

	notifyStep := maxSteps / 1000

	var options map[string]Option
	var jsonKey string
	var args []string

	var loader MultiLoader
	var config, origin map[string]string
	var err error

	dump := func() {
		t.Errorf("config: %#v", config)
		t.Errorf("origin: %#v", origin)
		t.Errorf("err: %#v", err)
		t.Errorf("%#v\n", loader)
	}

	f := fuzz.New()
	defer func() {
		if e := recover(); e != nil {
			t.Errorf("panic\n")
			dump()
			fmt.Printf("%v\n%v", err, string(debug.Stack()))
			t.FailNow()
		}
	}()

	for i := 0; i <= maxSteps; i++ {
		f.Fuzz(&options)
		f.Fuzz(&jsonKey)
		f.Fuzz(&args)

		options := map[string]Option{}

		loader = MultiLoader{Options: options, JSONKey: jsonKey}

		config, origin, err = loader.load(args, sampleFlagsHandler)

		if err != nil && (len(config) > 0 || len(origin) > 0) {
			t.Errorf("error present but output not empty\n")
			dump()
			t.FailNow()
		}

		if i%notifyStep == 0 {
			fmt.Fprintf(os.Stdout, "\r%0.1f%% ", float64(i)/float64(maxSteps)*100.0)
		}
	}
	fmt.Fprintf(os.Stderr, "\nComplete...\n")
}

func suppressedUsage(flags *flag.FlagSet) func() {
	return func() {
		flags.SetOutput(ioutil.Discard)
	}
}
