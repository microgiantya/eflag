package eflag

import (
	"flag"
	"fmt"
	"reflect"
)

var errDefaultMessage = `
Currently supported kinds (in term of [reflect] package) (including pointers and aliases):
	bool, string, int64, float64, struct.
A double pointer returns an error without the possibility of traversal.
Any other kinds return errors, but can be skipped with provided [OptionContinueOnUnknownKind] option.
`

func checkInput(t any) error {
	if reflect.ValueOf(t).Kind() != reflect.Pointer || reflect.ValueOf(t).Elem().Kind() != reflect.Struct {
		return ErrInvalidInput
	}
	return nil
}

func parseFromFlagSet(t any, options option, flagSet *flag.FlagSet, argumentList []string) error {
	if err := checkInput(t); err != nil {
		return fmt.Errorf("%w, %s", err, errDefaultMessage)
	}

	if err := parseToStruct(t, flagSet, options); err != nil {
		return fmt.Errorf("%w, %s", err, errDefaultMessage)
	}

	if err := flagSet.Parse(argumentList); err != nil {
		return fmt.Errorf("%w, %s", err, errDefaultMessage)
	}

	return nil
}

func Parse(t any, options option) error {
	if err := checkInput(t); err != nil {
		return err
	}
	if err := parseToStruct(t, flag.CommandLine, options); err != nil {
		return err
	}

	flag.Parse()
	return nil
}
