package conversion

import (
	"reflect"
	"testing"
)

func TestToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Test 1",
			args: args{s: "Hello, World!"},
			want: []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{b: []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}},
			want: "Hello, World!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.args.b); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
