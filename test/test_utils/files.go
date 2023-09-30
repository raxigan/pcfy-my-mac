package test_utils

import (
	"io"
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

func CopyFile(src, dst string) {
	sourceFile, _ := os.Open(src)
	defer sourceFile.Close()

	destFile, _ := os.Create(dst)
	defer destFile.Close()

	io.Copy(destFile, sourceFile)
}

func Trim(yaml string) string {
	return strings.TrimSpace(strings.ReplaceAll(yaml, "\t", ""))
}
