package test_utils

import (
	"os"
	"path/filepath"
	"strings"
)

func RemoveFiles(paths ...string) {
	for _, path := range paths {
		os.Remove(path)
	}
}

func RemoveFilesWithExt(path, ext string) {
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(filePath, ext) {
			os.Remove(filePath)
		}
		return nil
	})
}

func Trim(yaml string) string {
	return strings.TrimSpace(strings.ReplaceAll(yaml, "\t", ""))
}
