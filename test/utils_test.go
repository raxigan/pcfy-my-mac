package install_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func removeFiles(paths ...string) {
	for _, path := range paths {
		os.Remove(path)
	}
}

func removeFilesWithExt(path, ext string) {
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(filePath, ext) {
			os.Remove(filePath)
		}
		return nil
	})
}

func copyFile(src, dst string) {
	sourceFile, _ := os.Open(src)
	defer sourceFile.Close()

	destFile, _ := os.Create(dst)
	defer destFile.Close()

	io.Copy(destFile, sourceFile)
}
