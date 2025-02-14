package eflag

import (
	"encoding/json"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	BoolCustom         bool
	Int64Custom        int64
	Float64Custom      float64
	StringCustom       string
	TimeDurationCustom time.Duration
)

func Test_parseFromFlagSet_negative(t *testing.T) {
	type parseFromFlagSetArgs struct {
		t            any
		argumentList []string
		options      option
	}

	testCases := []struct {
		err  error
		name string
		args parseFromFlagSetArgs
	}{
		{
			name: "invalid input",
			args: parseFromFlagSetArgs{
				t: "",
			},
			err: ErrInvalidInput,
		},
		{
			name: "unknown kind",
			args: parseFromFlagSetArgs{
				t: &struct {
					Struct struct {
						Int32 *int32 `efName:"int32"`
					} `efName:"struct"`
					Int32 int32
				}{},
			},
			err: ErrUnknownKind,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := parseWithFlagSet(
				flag.NewFlagSet(time.Now().String(), flag.ExitOnError),
				testCase.args.argumentList,
				testCase.args.t,
				testCase.args.options,
			)
			if !assert.ErrorIs(t, err, testCase.err) {
				t.Errorf("parseFromFlagSet() want error: (%v) actual error: (%v)", testCase.err, err)
			}
		})
	}
}
func Test_parseFromFlagSet_negative_alreadyParsed(t *testing.T) {
	type parseFromFlagSetArgs struct {
		t            any
		argumentList []string
		options      option
	}

	testCases := []struct {
		err  error
		name string
		args parseFromFlagSetArgs
	}{
		{
			name: "already parsed",
			args: parseFromFlagSetArgs{
				t: &struct{}{},
			},
			err: ErrAlreadyParsed,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			flagSet := flag.NewFlagSet(time.Now().String(), flag.ExitOnError)
			err := parseWithFlagSet(
				flagSet,
				testCase.args.argumentList,
				testCase.args.t,
				testCase.args.options,
			)
			assert.Nil(t, err)
			err = parseWithFlagSet(
				flagSet,
				testCase.args.argumentList,
				testCase.args.t,
				testCase.args.options,
			)
			if !assert.ErrorIs(t, err, testCase.err) {
				t.Errorf("parseFromFlagSet() want error: (%v) actual error: (%v)", testCase.err, err)
			}
		})
	}
}

func Test_parseFromFlagSet_panic(t *testing.T) {
	type parseFromFlagSetArgs struct {
		t            any
		argumentList []string
		options      option
	}

	testCases := []struct {
		name string
		args parseFromFlagSetArgs
	}{
		{
			name: "flag redefined panic",
			args: parseFromFlagSetArgs{
				t: &struct {
					Flag     string `efName:"flag"`
					SameFlag string `efName:"flag"`
				}{},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("did not panic")
				}
			}()
			_ = parseWithFlagSet(
				flag.NewFlagSet(time.Now().String(), flag.ExitOnError),
				testCase.args.argumentList,
				testCase.args.t,
				testCase.args.options,
			)
		})
	}
}

func Test_parseFromFlagSet_positive_flag_only(t *testing.T) {
	type parseFromFlagSetArgs struct {
		t            any
		argumentList []string
		options      option
	}

	testCases := []struct {
		err          error
		name         string
		expectedJSON string
		args         parseFromFlagSetArgs
	}{
		{
			name: "positive",
			args: parseFromFlagSetArgs{
				t: &struct {
					Struct struct {
						String       *string        `efName:"string"`
						Bool         *bool          `efName:"bool"`
						Float64      *float64       `efName:"float64"`
						Int64        *int64         `efName:"int64"`
						TimeDuration *time.Duration `efName:"time-duration"`
						NoSetString  *string
					} `efName:"struct"`
					String       string        `efName:"string"`
					Float64      float64       `efName:"float64"`
					Int64        int64         `efName:"int64"`
					TimeDuration time.Duration `efName:"time-duration"`
					Bool         bool          `efName:"bool"`
				}{},
				argumentList: []string{
					"-string=mystring",
					"-bool",
					"-float64=2",
					"-int64=1",
					"-time-duration=5s",
					"-struct-string=mystructstring",
					"-struct-bool=false",
					"-struct-float64=3",
					"-struct-int64=88",
					"-struct-time-duration=3s",
				},
			},
			expectedJSON: `{"Struct":{"String":"mystructstring","Bool":false,"Float64":3,"Int64":88,"TimeDuration":3000000000,"NoSetString":""},"String":"mystring","Float64":2,"Int64":1,"TimeDuration":5000000000,"Bool":true}`,
		},
		{
			name: "positive custom",
			args: parseFromFlagSetArgs{
				t: &struct {
					Struct struct {
						StringCustom       *StringCustom       `efName:"string"`
						BoolCustom         *BoolCustom         `efName:"bool"`
						Float64Custom      *Float64Custom      `efName:"float64"`
						Int64              *Int64Custom        `efName:"int64"`
						TimeDurationCustom *TimeDurationCustom `efName:"time-duration"`
						NoSetString        *StringCustom
					} `efName:"struct"`
					StringCustom       StringCustom       `efName:"string"`
					Float64Custom      Float64Custom      `efName:"float64"`
					Int64Custom        Int64Custom        `efName:"int64"`
					TimeDurationCustom TimeDurationCustom `efName:"time-duration"`
					BoolCustom         BoolCustom         `efName:"bool"`
				}{},
				argumentList: []string{
					"-string=mystring",
					"-bool",
					"-float64=2",
					"-int64=1",
					"-struct-string=mystructstring",
					"-struct-bool", "false",
					"--struct-float64=3",
					"-struct-int64", "88",
				},
			},
			expectedJSON: `{"Struct":{"StringCustom":"mystructstring","BoolCustom":true,"Float64Custom":0,"Int64":0,"TimeDurationCustom":0,"NoSetString":""},"StringCustom":"mystring","Float64Custom":2,"Int64Custom":1,"TimeDurationCustom":0,"BoolCustom":true}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := parseWithFlagSet(
				flag.NewFlagSet("test", flag.ExitOnError),
				testCase.args.argumentList,
				testCase.args.t,
				testCase.args.options,
			)
			assert.Nil(t, err)
			actualJSON, _ := json.Marshal(testCase.args.t)
			assert.Equal(t, testCase.expectedJSON, string(actualJSON))
		})
	}
}

func Test_parseFromFlagSet_positive_with_env(t *testing.T) {
	type parseFromFlagSetArgs struct {
		t            any
		argumentList []string
		envList      map[string]string
		options      option
	}

	testCases := []struct {
		err          error
		name         string
		expectedJSON string
		args         parseFromFlagSetArgs
	}{
		{
			name: "positive",
			args: parseFromFlagSetArgs{
				t: &struct {
					Struct struct {
						String       *string        `efName:"string"`
						Bool         *bool          `efName:"bool"`
						Float64      *float64       `efName:"float64"`
						Int64        *int64         `efName:"int64"`
						TimeDuration *time.Duration `efName:"time-duration"`
						NoSetString  *string
					} `efName:"struct"`
					String       string        `efName:"string"`
					Float64      float64       `efName:"float64"`
					Int64        int64         `efName:"int64"`
					TimeDuration time.Duration `efName:"time-duration"`
					Bool         bool          `efName:"bool"`
				}{},
				argumentList: []string{
					"-string=mystring",
					"-bool",
					"-int64=1",
					"-time-duration=5s",
					"-struct-string=mystructstring",
					"-struct-bool=false",
					"-struct-float64=3",
					"-struct-int64=88",
				},
				envList: map[string]string{
					"APP_FLOAT64":              "45.6",
					"APP_STRUCT_TIME_DURATION": "10s",
					"APP_STRUCT_INT64":         "89",
				},
				options: WithEnv,
			},
			expectedJSON: `{"Struct":{"String":"mystructstring","Bool":false,"Float64":3,"Int64":88,"TimeDuration":10000000000,"NoSetString":""},"String":"mystring","Float64":45.6,"Int64":1,"TimeDuration":5000000000,"Bool":true}`,
		},
		{
			name: "positive custom",
			args: parseFromFlagSetArgs{
				t: &struct {
					Struct struct {
						StringCustom       *StringCustom       `efName:"string"`
						BoolCustom         *BoolCustom         `efName:"bool"`
						Float64Custom      *Float64Custom      `efName:"float64"`
						Int64Custom        *Int64Custom        `efName:"int64"`
						TimeDurationCustom *TimeDurationCustom `efName:"time-duration"`
						NoSetString        *StringCustom
					} `efName:"struct"`
					StringCustom       StringCustom       `efName:"string"`
					Float64Custom      Float64Custom      `efName:"float64"`
					Int64Custom        Int64Custom        `efName:"int64"`
					TimeDurationCustom TimeDurationCustom `efName:"time-duration"`
					BoolCustom         BoolCustom         `efName:"bool"`
				}{},
				argumentList: []string{
					"-string=mystring",
					"-bool",
					"-float64=2",
					"-int64=1",
					"-struct-string=mystructstring",
					"-struct-bool", "false",
					"--struct-float64=3",
					"-struct-int64", "88",
				},
			},
			expectedJSON: `{"Struct":{"StringCustom":"mystructstring","BoolCustom":true,"Float64Custom":0,"Int64Custom":0,"TimeDurationCustom":0,"NoSetString":""},"StringCustom":"mystring","Float64Custom":2,"Int64Custom":1,"TimeDurationCustom":0,"BoolCustom":true}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			setEnv(testCase.args.envList)
			defer unsetEnv(testCase.args.envList)
			err := parseWithFlagSet(
				flag.NewFlagSet("test", flag.ExitOnError),
				testCase.args.argumentList,
				testCase.args.t,
				testCase.args.options,
			)
			assert.Nil(t, err)
			actualJSON, _ := json.Marshal(testCase.args.t)
			assert.Equal(t, testCase.expectedJSON, string(actualJSON))
		})
	}
}

func setEnv(envList map[string]string) {
	for k, v := range envList {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}
}
func unsetEnv(envList map[string]string) {
	for k := range envList {
		if err := os.Unsetenv(k); err != nil {
			panic(err)
		}
	}
}
