package gofake

import "testing"

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
