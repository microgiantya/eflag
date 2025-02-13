package constvalues_test

import (
	"reflect"
	"testing"

	"git.zonatelecom.ru/e.artemev/eflag/pkg/constvalues"
	"git.zonatelecom.ru/e.artemev/eflag/pkg/constvalues/testtypes"
)

func TestGetConstValuesByType(t *testing.T) {
	type args struct {
		t any
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "exported typed string",
			args: args{testtypes.OneEx},
			want: []string{"four", "one", "three", "two"},
		},
		{
			name: "unexported typed string",
			args: args{testtypes.OneUn},
			want: []string{"four", "one", "three", "two"},
		},
		{
			name: "exported string",
			args: args{testtypes.One},
			want: []string(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := constvalues.GetConstValuesByType(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConstValuesByType() = %v, want %v", got, tt.want)
			}
		})
	}
}
