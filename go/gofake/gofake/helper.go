package gofake

import (
	"os"
	"path/filepath"
	"strings"
)

func makeGenFilePath(filePath string) string {
	ext := filepath.Ext(filePath)
	return filePath[:strings.LastIndex(filePath, ext)] + "_gen" + ext
}

func isFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
