package main

import (
	"reflect"
	"testing"
)

func Test_execute(t *testing.T) {

	type args struct {
		expression string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"OR operator", args{"foo(-(bar|boo))"}, []string{"foo-bar", "foo-boo"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := execute(tt.args.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ok := false
			for _, possibility := range tt.want {
				ok = reflect.DeepEqual(got.String(), possibility)
			}

			if !ok {
				t.Errorf("%v: execute() = %v", tt.name, got.String())
			}
		})
	}
}
