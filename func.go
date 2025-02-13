package eflag

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func namespaceAdapt(namespace string) string {
	if namespace != "" {
		namespace += "-"
	}
	return namespace
}

func getFlagName(namespace, flagName string) string {
	return namespaceAdapt(namespace) + flagName
}

func getFlagNameColor(namespace, flagName string, option option) string {
	flagName = namespaceAdapt(namespace) + flagName
	if WithColor.isSet(option) {
		flagName = color.RedString(flagName)
	}
	return flagName
}

func getUsage(flagUsage, flagName string, option option) string {
	if flagUsage == "" {
		flagUsage = "set " + flagName
	}

	// if WithColor.isSet(option) {
	// 	flagUsage = color.GreenString(flagUsage)
	// }

	if WithEnv.isSet(option) {
		flagUsage = "env: " + getEnvColor(flagName, option) + "\n" + flagUsage
	}
	return flagUsage
}

func getEnvColor(s string, option option) string {
	env := "APP_" + strings.ReplaceAll(strings.ToUpper(s), "-", "_")
	if WithColor.isSet(option) {
		return color.CyanString(env)
	}
	return env
}

func getEnv(s string) string {
	return "APP_" + strings.ReplaceAll(strings.ToUpper(s), "-", "_")
}

func getValueBool(t any, option option, flagName string) bool {
	v := reflect.ValueOf(t).Bool()
	if WithEnv.isSet(option) {
		if value, ok := os.LookupEnv(getEnv(flagName)); ok {
			if valueBool, err := strconv.ParseBool(value); err == nil {
				v = valueBool
			}
		}
	}
	return v
}

func getValueDuration(t any, option option, flagName string) time.Duration {
	v := time.Duration(reflect.ValueOf(t).Int())
	if WithEnv.isSet(option) {
		if value, ok := os.LookupEnv(getEnv(flagName)); ok {
			if valueBool, err := time.ParseDuration(value); err == nil {
				v = valueBool
			}
		}
	}
	return v
}

func getValueInt64(t any, option option, flagName string) int64 {
	v := reflect.ValueOf(t).Int()
	if WithEnv.isSet(option) {
		if value, ok := os.LookupEnv(getEnv(flagName)); ok {
			if valueBool, err := strconv.ParseInt(value, 10, 64); err == nil {
				v = valueBool
			}
		}
	}
	return v
}

func getValueFloat64(t any, option option, flagName string) float64 {
	v := reflect.ValueOf(t).Float()
	if WithEnv.isSet(option) {
		if value, ok := os.LookupEnv(getEnv(flagName)); ok {
			if valueBool, err := strconv.ParseFloat(value, 64); err == nil {
				v = valueBool
			}
		}
	}
	return v
}

func getValueString(t any, option option, flagName string) string {
	v := reflect.ValueOf(t).String()
	if WithEnv.isSet(option) {
		if value, ok := os.LookupEnv(getEnv(flagName)); ok {
			v = value
		}
	}
	return v
}
