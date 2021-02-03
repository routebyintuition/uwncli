package main

import (
	"testing"
)

type stubInputReader struct {
	Input string
}

func (ir stubInputReader) ReadInputSecure() (string, error) {
	return ir.Input, nil
}

func TestGetMibFromMB(t *testing.T) {
	type args struct {
		mb int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1GB",
			args: args{
				mb: 1000,
			},
			want: 953,
		},
		{
			name: "10GB",
			args: args{
				mb: 10000,
			},
			want: 9536,
		},
		{
			name: "100GB",
			args: args{
				mb: 100000,
			},
			want: 95367,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMibFromMB(tt.args.mb); got != tt.want {
				t.Errorf("GetMibFromMB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidUUID(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "VALID UUID",
			args: args{
				uuid: "E64AD4E5-C6E2-462B-881F-A02BD0CDD8BB",
			},
			want: true,
		},
		{
			name: "INALID UUID",
			args: args{
				uuid: "not a valid uuid",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidUUID(tt.args.uuid); got != tt.want {
				t.Errorf("IsValidUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInputStringValue(t *testing.T) {
	// ir := &stubInputReader{}

	type args struct {
		ir      InputReader
		message string
		minLen  int
		def     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "SHORT PASS",
			args: args{
				ir:      &stubInputReader{Input: "demo"},
				message: "enter password",
				minLen:  8,
				def:     "",
			},
			//want:    "test",
			wantErr: true,
		},
		{
			name: "GOOD PASS",
			args: args{
				ir:      &stubInputReader{Input: "demodemodemo"},
				message: "enter password",
				minLen:  8,
				def:     "",
			},
			want:    "demodemodemo",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInputStringValue(tt.args.ir, tt.args.message, tt.args.minLen, tt.args.def)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInputStringValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInputStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
