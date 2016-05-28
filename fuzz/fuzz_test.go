package fuzz

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"github.com/chiku/conf"
	"github.com/google/gofuzz"
)

const (
	max        = 100000000
	notifyStep = 10000
)

func TestConfFuzzRandom(t *testing.T) {
	var jsonKey string
	var mandatory []string
	var optional []string
	var args []string
	var defaults map[string]string

	var loader conf.MultiLoader
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

	for i := 0; i <= max; i++ {
		f.Fuzz(&jsonKey)
		f.Fuzz(&mandatory)
		f.Fuzz(&optional)
		f.Fuzz(&args)
		f.Fuzz(&defaults)

		loader = conf.MultiLoader{
			JSONKey:   jsonKey,
			Mandatory: mandatory,
			Optional:  optional,
			Args:      args,
			Defaults:  defaults,
		}

		config, origin, err = loader.Load()

		if err != nil && (len(config) > 0 || len(origin) > 0) {
			t.Errorf("error present but output not empty\n")
			dump()
			t.FailNow()
		}

		if i%notifyStep == 0 {
			fmt.Fprintf(os.Stderr, "\r%0.1f%%", float64(i)/float64(max)*100.0)
		}
	}
	fmt.Fprintf(os.Stderr, "Complete...\n")
}
