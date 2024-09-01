package test_utils

import (
	"os"
	"path/filepath"
	"strings"
)

func RemoveFiles(paths ...string) error {
	for _, path := range paths {
		err := os.RemoveAll(path)

		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveFilesWithExt(path, ext string) {
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(filePath, ext) {
			os.Remove(filePath)
		}
		return nil
	})
}

func RemoveDirs(path, dir string) {
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.HasSuffix(filePath, dir) {
			os.Remove(filePath)
		}
		return nil
	})
}

func Trim(yaml string) string {
	return strings.TrimSpace(strings.ReplaceAll(yaml, "\t", ""))
}
