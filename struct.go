package eflag

import (
	"flag"
	"reflect"
)

const (
	_flagName  = "efName"
	_flagUsage = "efUsage"
)

func parseToStruct(t any, flagSet *flag.FlagSet, options option, namespace string) error {
	var val = reflect.ValueOf(t).Elem()
	var typ = val.Type()

	for i := range typ.NumField() {
		var fieldValue = val.Field(i)
		if !typ.Field(i).IsExported() {
			continue
		}

		var crr = carrier{
			efName:  typ.Field(i).Tag.Get(_flagName),
			efUsage: typ.Field(i).Tag.Get(_flagUsage),
		}

		switch typ.Field(i).Type.Kind() {
		case reflect.Pointer:
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(typ.Field(i).Type.Elem()).Elem().Addr())
			}
			crr.ptr = fieldValue.Interface()
			crr.uptr = fieldValue.UnsafePointer()
			crr.value = fieldValue.Elem().Interface()
		default:
			crr.ptr = fieldValue.Addr().Interface()
			crr.uptr = fieldValue.Addr().UnsafePointer()
			crr.value = fieldValue.Interface()
		}

		if err := parseToStructFiled(crr, flagSet, options, namespace); err != nil {
			return err
		}
	}
	return nil
}
