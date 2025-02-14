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

func parseToStructFiled(crr carrier, flagSet *flag.FlagSet, option option, namespace string) error {
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

	flagName := getFlagName(namespace, crr.efName)
	flagNameC := getFlagNameColor(namespace, crr.efName, option)
	flagUsage := getUsage(crr.efUsage, namespaceAdapt(namespace)+crr.efName, option)

	switch kind {
	case reflect.Bool:
		flagSet.BoolVar(
			(*bool)(crr.uptr),
			flagNameC,
			getValueBool(crr.value, option, flagName),
			flagUsage,
		)
	case reflect.Int64:
		switch crr.value.(type) {
		case time.Duration:
			flagSet.DurationVar(
				(*time.Duration)(crr.uptr),
				flagNameC,
				getValueDuration(reflect.ValueOf(crr.value).Interface().(time.Duration), option, flagName), //nolint:all // no error returned here.
				flagUsage,
			)
			return nil
		default:
		}
		flagSet.Int64Var(
			(*int64)(crr.uptr),
			flagNameC,
			getValueInt64(reflect.ValueOf(crr.value).Int(), option, flagName),
			flagUsage,
		)
	case reflect.Float64:
		flagSet.Float64Var(
			(*float64)(crr.uptr),
			flagNameC,
			getValueFloat64(reflect.ValueOf(crr.value).Float(), option, flagName),
			flagUsage,
		)
	case reflect.String:
		flagSet.StringVar(
			(*string)(crr.uptr),
			flagNameC,
			getValueString(reflect.ValueOf(crr.value).String(), option, flagName),
			flagUsage,
		)
	case reflect.Struct:
		if err := parseToStruct(crr.ptr, flagSet, option, namespaceAdapt(namespace)+crr.efName); err != nil {
			return err
		}
	default:
		return ErrUnknownKind
	}

	return nil
}
