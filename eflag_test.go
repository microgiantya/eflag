package eflag

import (
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	BoolAlias         bool
	Int64Alias        int64
	Float64Alias      float64
	StringAlias       string
	TimeDurationAlias time.Duration
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

func Test_parseFromFlagSet_positive(t *testing.T) {
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
			err:          ErrAlreadyParsed,
			expectedJSON: `{"Struct":{"String":"mystructstring","Bool":false,"Float64":3,"Int64":88,"TimeDuration":3000000000,"NoSetString":""},"String":"mystring","Float64":2,"Int64":1,"TimeDuration":5000000000,"Bool":true}`,
		},
		{
			name: "positive alias",
			args: parseFromFlagSetArgs{
				t: &struct {
					Struct struct {
						StringAlias       *StringAlias       `efName:"string"`
						BoolAlias         *BoolAlias         `efName:"bool"`
						Float64Alias      *Float64Alias      `efName:"float64"`
						Int64             *Int64Alias        `efName:"int64"`
						TimeDurationAlias *TimeDurationAlias `efName:"time-duration"`
						NoSetString       *StringAlias
					} `efName:"struct"`
					StringAlias       StringAlias       `efName:"string"`
					Float64Alias      Float64Alias      `efName:"float64"`
					Int64Alias        Int64Alias        `efName:"int64"`
					TimeDurationAlias TimeDurationAlias `efName:"time-duration"`
					BoolAlias         BoolAlias         `efName:"bool"`
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
			err:          ErrAlreadyParsed,
			expectedJSON: `{"Struct":{"StringAlias":"mystructstring","BoolAlias":true,"Float64Alias":0,"Int64":0,"TimeDurationAlias":0,"NoSetString":""},"StringAlias":"mystring","Float64Alias":2,"Int64Alias":1,"TimeDurationAlias":0,"BoolAlias":true}`,
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
