package eflag

import (
	"flag"
	"reflect"
)

func checkInput(t any) error {
	if reflect.ValueOf(t).Kind() != reflect.Pointer || reflect.ValueOf(t).Elem().Kind() != reflect.Struct {
		return ErrInvalidInput
	}
	return nil
}

func ParseFromFlagSet(t any, options option, flagSet *flag.FlagSet, argumentList []string) error {
	if err := checkInput(t); err != nil {
		return err
	}

	if err := parseToStruct(t, flagSet, options); err != nil {
		return err
	}

	if err := flagSet.Parse(argumentList); err != nil {
		return err
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
