package install

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(src, dst string) {
	sourceFile, _ := os.Open(src)
	defer sourceFile.Close()

	destFile, _ := os.Create(dst)
	defer destFile.Close()

	io.Copy(destFile, sourceFile)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func findMatchingDirs(basePath, namePrefix, subDir, fileName string) ([]string, error) {

	var result []string

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {

		if path != basePath && strings.HasPrefix(info.Name(), namePrefix) {
			if err != nil {
				return err
			}

			if fileExists(filepath.Join(basePath, info.Name())) {
				destDir := filepath.Join(path, subDir)
				destFilePath := filepath.Join(destDir, fileName)
				result = append(result, destFilePath)
			}
		}

		return nil
	})

	return result, err
}

func replaceWordInFile(path, oldWord, newWord string) error {

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	modifiedContent := strings.ReplaceAll(string(content), oldWord, newWord)

	err = os.WriteFile(path, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func TextFromFile(paramsFile string) (string, error) {
	d, e := os.ReadFile(paramsFile)

	if e != nil {
		return "", e
	}

	return string(d), nil
}
