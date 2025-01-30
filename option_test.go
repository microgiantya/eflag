package eflag

import (
	"testing"
)

func Test_newOption(t *testing.T) {
	type args struct {
		options []option
	}
	testCaseList := []struct {
		name string
		args args
		want option
	}{
		{
			name: "empty input options",
			want: option(0),
		},
		{
			name: "not empty input options",
			args: args{options: []option{option(1)}},
			want: option(1),
		},
		{
			name: "not empty input options",
			args: args{options: []option{option(1), option(135)}},
			want: option(135),
		},
		{
			name: "not empty input options",
			args: args{options: []option{
				option(1),
				option(2),
				option(4),
				option(8),
				option(16),
				option(32),
				option(64),
				option(128),
			}},
			want: option(255),
		},
	}
	for _, testCase := range testCaseList {
		t.Run(testCase.name, func(t *testing.T) {
			if got := newOption(testCase.args.options...); got != testCase.want {
				t.Errorf("newOption() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func Test_option_isSet(t *testing.T) {
	type args struct {
		providedOption option
	}
	tests := []struct {
		name string
		tr   option
		args args
		want bool
	}{
		{
			name: "is set: option and provided the same",
			tr:   option(8),
			args: args{providedOption: option(8)},
			want: true,
		},
		{
			name: "is set: provided contains option",
			tr:   option(8),
			args: args{providedOption: option(255)},
			want: true,
		},
		{
			name: "is not set: provided not contains option",
			tr:   option(1),
			args: args{providedOption: option(10)},
			want: false,
		},
		{
			name: "is not set: option is 0",
			tr:   option(0),
			args: args{providedOption: option(10)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.isSet(tt.args.providedOption); got != tt.want {
				t.Errorf("option.isSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
