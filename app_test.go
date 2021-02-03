package main

import (
	"reflect"
	"testing"

	"github.com/urfave/cli/v2"
)

func Test_getFlags(t *testing.T) {
	tests := []struct {
		name string
		want []cli.Flag
	}{
		{
			name: "Check Flags",
			want: getFlags(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFlags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}
