package eflag

import (
	"flag"
	"reflect"
	"time"
)

func structFieldValidate(t any) error {
	if reflect.ValueOf(t).Elem().Kind() == reflect.Ptr {
		return ErrInvalidInput
	}
	return nil
}

func parseToStructFiled(crr carrier, flagSet *flag.FlagSet, options option, namespace string) error {
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

	if namespace != "" {
		namespace += "-"
	}

	switch kind {
	case reflect.Bool:
		flagSet.BoolVar(
			(*bool)(crr.uptr),
			namespace+crr.efName,
			reflect.ValueOf(crr.value).Bool(),
			crr.efUsage,
		)
	case reflect.Int64:
		switch crr.value.(type) {
		case time.Duration:
			flagSet.DurationVar(
				(*time.Duration)(crr.uptr),
				namespace+crr.efName,
				reflect.ValueOf(crr.value).Interface().(time.Duration),
				crr.efUsage,
			)
			return nil
		default:
		}
		flagSet.Int64Var(
			(*int64)(crr.uptr),
			namespace+crr.efName,
			reflect.ValueOf(crr.value).Int(),
			crr.efUsage,
		)
	case reflect.Float64:
		flagSet.Float64Var(
			(*float64)(crr.uptr),
			namespace+crr.efName,
			reflect.ValueOf(crr.value).Float(),
			crr.efUsage,
		)
	case reflect.String:
		flagSet.StringVar(
			(*string)(crr.uptr),
			namespace+crr.efName,
			reflect.ValueOf(crr.value).String(),
			crr.efUsage,
		)
	case reflect.Struct:
		if err := parseToStruct(crr.ptr, flagSet, options, crr.efName); err != nil {
			return err
		}
	default:
		if !OptionContinueOnUnknownKind.isSet(options) {
			return ErrKindIsNotSupported
		}
	}
	return nil
}
