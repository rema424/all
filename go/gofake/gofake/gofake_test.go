package gofake

import (
	"reflect"
	"testing"
)

func Test_makeFileContents(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"", args{"animal"}, "import animal\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeFileContents(tt.args.pkg)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeFileContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[]byte(makeFileContents()) = %v, want %v", got, tt.want)
			}
		})
	}
}
