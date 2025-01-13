package eflag

import (
	"flag"
	"reflect"
)

func structFieldValidate(t any) error {
	if reflect.ValueOf(t).Elem().Kind() == reflect.Ptr {
		return ErrInvalidInput
	}
	return nil
}

func parseToStructFiled(crr carrier, flagSet *flag.FlagSet, options option) error {
	if err := structFieldValidate(crr.ptr); err != nil {
		return err
	}

	var val = reflect.ValueOf(crr.value)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	var kind = val.Kind()

	if crr.efName == "" && kind != reflect.Struct {
		return nil
	}

	switch kind {
	case reflect.Bool:
		flagSet.BoolVar(
			(*bool)(crr.uptr),
			crr.efName,
			reflect.ValueOf(crr.value).Bool(),
			crr.efUsage,
		)
	case reflect.Int64:
		flagSet.Int64Var(
			(*int64)(crr.uptr),
			crr.efName,
			reflect.ValueOf(crr.value).Int(),
			crr.efUsage,
		)
	case reflect.Float64:
		flagSet.Float64Var(
			(*float64)(crr.uptr),
			crr.efName,
			reflect.ValueOf(crr.value).Float(),
			crr.efUsage,
		)
	case reflect.String:
		flagSet.StringVar(
			(*string)(crr.uptr),
			crr.efName,
			reflect.ValueOf(crr.value).String(),
			crr.efUsage,
		)
	case reflect.Struct:
		if err := parseToStruct(crr.ptr, flagSet, options); err != nil {
			return err
		}
	default:
		if !OptionContinueOnUnknownKind.isSet(options) {
			return ErrKindIsNotSupported
		}
	}
	return nil
}
