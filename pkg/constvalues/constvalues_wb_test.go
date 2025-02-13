package constvalues

import (
	"reflect"
	"testing"
	"time"

	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
)

func Test_getRootType(t *testing.T) {
	type args struct {
		t any
	}
	tests := []struct {
		args args
		want reflect.Type
		name string
	}{
		{
			name: "pointer",
			args: args{lo.ToPtr("")},
			want: reflect.TypeOf(""),
		},
		{
			name: "string",
			args: args{""},
			want: reflect.TypeOf(""),
		},
		{
			name: "time duration",
			args: args{lo.ToPtr(time.Second)},
			want: reflect.TypeOf(time.Second),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRootType(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRootType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkKind(t *testing.T) {
	type args struct {
		typ reflect.Type
	}
	tests := []struct {
		args args
		name string
		want bool
	}{
		{
			name: "string",
			args: args{reflect.TypeOf("")},
			want: true,
		},
		{
			name: "int64",
			args: args{reflect.TypeOf(int64(1))},
			want: true,
		},
		{
			name: "float64",
			args: args{reflect.TypeOf(float64(2))},
			want: true,
		},
		{
			name: "bool",
			args: args{reflect.TypeOf(false)},
			want: true,
		},
		{
			name: "time duration",
			args: args{reflect.TypeOf(time.Second)},
			want: true,
		},
		{
			name: "int32",
			args: args{reflect.TypeOf(int32(3))},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkKind(tt.args.typ); got != tt.want {
				t.Errorf("checkKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadPackage(t *testing.T) {
	type args struct {
		packagePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "existing package",
			args:    args{"os"},
			wantErr: false,
		},
		{
			name:    "not existing package",
			args:    args{"osasdfasdfasaf"},
			wantErr: true,
		},
		{
			name:    "empty package path (self path)",
			args:    args{""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := loadPackage(tt.args.packagePath)
			if len(got.Errors) != 0 && !tt.wantErr {
				t.Errorf("loadPackage() errors len = %v, wantErr %v", len(got.Errors), tt.wantErr)
			}
			if len(got.Errors) == 0 && tt.wantErr {
				t.Errorf("loadPackage() errors len = %v, wantErr %v", len(got.Errors), tt.wantErr)
			}
		})
	}
}

func Test_getConstValueList(t *testing.T) {
	type args struct {
		pkg      *packages.Package
		typeName string
	}
	tests := []struct {
		want map[string][]string
		args args
		name string
	}{
		{
			name: "os int exported values",
			args: args{pkg: loadPackage("os"), typeName: "int"},
			want: map[string][]string{
				"int": {"0", "1", "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConstValueList(tt.args.pkg, tt.args.typeName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConstValueList() = %v, want %v", got, tt.want)
			}
		})
	}
}
