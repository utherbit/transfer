package main

import (
	"go/token"
	"reflect"
	"testing"
)

func TestParsePosition(t *testing.T) {
	type args struct {
		ref string
	}
	tests := []struct {
		name    string
		args    args
		wantR   token.Position
		wantErr bool
	}{
		{"file:line", args{ref: "example/directory:42"}, token.Position{Filename: "example/directory", Line: 42}, false},
		{"file:line:col", args{ref: "example/directory:42:24"}, token.Position{Filename: "example/directory", Line: 42, Column: 24}, false},
		{"noline", args{ref: "example/directory"}, token.Position{Filename: "example/directory"}, true},
		{"empty", args{ref: ""}, token.Position{}, true},
		{"lineNaN", args{ref: "example/directory:NaN"}, token.Position{Filename: "example/directory"}, true},
		{"colNaN", args{ref: "example/directory:1:NaN"}, token.Position{Filename: "example/directory", Line: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := ParsePosition(tt.args.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGoRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("ParseGoRef() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
