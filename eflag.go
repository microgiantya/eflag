// Package eflag allows command line flags to be bound to structure fields
// by setting the structure field tags 'efName' and 'efUsage'.
// For the actual parsing, the standard [flag] package is used without overriding
// its behavior.
// If the 'efName' tag exists for a nested struct type field, it will be
// added as a prefix to the fields with not empty 'efName' tag of this struct
// (except struct type).
package eflag

import (
	"flag"
	"net/http"
	"os"
	"reflect"

	json "github.com/goccy/go-json"
)

var parsed []byte

// Parse accepts not nil pointer to struct as a first argument.
// Supported field types (including pointers and custom types) is
// bool, string, int64, float64, struct and [time.Duration].
// Any nil pointer types will be intialized.
// Fields of type pointer to pointer cause an error.
// For now options is only for future compatibility.
//
// Available options is:
// WithEnv - reads the env variable with lower priority.
// The env variable name is created from the flag name by converting
// all letters to uppercase, replacing dashes with underscores,
// and adding the "APP_" prefix.
//
// WithColor - highlights the flag name and, if enables WithEnv, the environment
// variable name in the help.
func Parse(t any, options ...option) error {
	return parseWithFlagSet(flag.CommandLine, os.Args[1:], t, options...)
}

func parseWithFlagSet(flagSet *flag.FlagSet, argumentList []string, t any, options ...option) error {
	if flagSet.Parsed() {
		return errWrap(ErrAlreadyParsed)
	}

	if err := checkInput(t); err != nil {
		return errWrap(err)
	}

	option := newOption(options...)

	if err := parseToStruct(t, flagSet, option, ""); err != nil {
		return errWrap(err)
	}

	if err := flagSet.Parse(argumentList); err != nil {
		return errWrap(err)
	}

	parsed, _ = json.Marshal(t)
	return nil
}

// Handler write json representation of provided struct.
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(parsed)
	}
}

func checkInput(t any) error {
	if reflect.ValueOf(t).Kind() != reflect.Pointer || reflect.ValueOf(t).Elem().Kind() != reflect.Struct {
		return ErrInvalidInput
	}
	return nil
}
