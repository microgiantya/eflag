package eflag

import (
	"flag"
	"testing"
)

func Test_parseFromFlagSet(t *testing.T) {
	type args struct {
		t            any
		options      option
		flagSet      *flag.FlagSet
		argumentList []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseFromFlagSet(
				tt.args.t,
				tt.args.options,
				tt.args.flagSet,
				tt.args.argumentList,
			); (err != nil) != tt.wantErr {
				t.Errorf("parseFromFlagSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
