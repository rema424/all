package gofake

import (
	"path/filepath"
	"strings"
)

func makeGenFilePath(filePath string) string {
	ext := filepath.Ext(filePath)
	return filePath[:strings.LastIndex(filePath, ext)] + "_gen" + ext
}
