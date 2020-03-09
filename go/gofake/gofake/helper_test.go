package gofake

import (
	"os"
	"testing"
)

func Test_makeGenFilePath(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{"xxx.go"}, "xxx_gen.go"},
		{"", args{"xxx/yyy.go"}, "xxx/yyy_gen.go"},
		{"", args{"/xxx/yyy.go"}, "/xxx/yyy_gen.go"},
		{"", args{"./xxx/yyy.go"}, "./xxx/yyy_gen.go"},
		{"", args{"yyy.go.go"}, "yyy.go_gen.go"},
		{"", args{"あああ.go"}, "あああ_gen.go"},
		{"", args{"😄😄😄.go"}, "😄😄😄_gen.go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeGenFilePath(tt.args.filePath); got != tt.want {
				t.Errorf("makeGenFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isFileExists(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"xxx.txt"}, true},
		{"", args{"./xxx.txt"}, true},
		{"", args{"yyy.txt"}, false},
		{"", args{"www/zzz.txt"}, true},
	}

	if _, err := os.OpenFile("xxx.txt", os.O_CREATE, 0777); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("www"); os.IsNotExist(err) {
		if err := os.Mkdir("www", 0777); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.OpenFile("www/zzz.txt", os.O_CREATE, 0777); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("yyy.txt"); !os.IsNotExist(err) {
		if err := os.Remove("yyy.txt"); err != nil {
			t.Fatal(err)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFileExists(tt.args.filePath); got != tt.want {
				t.Errorf("isFileExists() = %v, want %v", got, tt.want)
			}
		})
	}

	if err := os.RemoveAll("www"); err != nil {
		t.Error(err)
	}
	if err := os.Remove("xxx.txt"); err != nil {
		t.Error(err)
	}
}
