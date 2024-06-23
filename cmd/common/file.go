package common

import (
	"github.com/raxigan/pcfy-my-mac/assets"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func CopyFileFromEmbedFS(src, dst string) error {
	assets := &assets.Assets
	data, _ := fs.ReadFile(assets, src)
	os.MkdirAll(filepath.Dir(dst), 0755)
	return os.WriteFile(dst, data, 0755)
}

func ReadFileFromEmbedFS(src string) (string, error) {
	configs := &assets.Assets
	data, _ := fs.ReadFile(configs, src)
	return string(data), nil
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)

	if err != nil {
		return err
	}

	defer sourceFile.Close()

	destFile, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)

	if err != nil {
		return err
	}

	return nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func FindMatchingPaths(basePath, namePrefix, subDir, fileName string) ([]string, error) {

	var result []string

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {

		if path != basePath && strings.HasPrefix(info.Name(), namePrefix) {
			if err != nil {
				return err
			}

			if FileExists(filepath.Join(basePath, info.Name())) {
				destDir := filepath.Join(path, subDir)
				destFilePath := filepath.Join(destDir, fileName)
				result = append(result, destFilePath)
			}
		}

		return nil
	})

	return result, err
}

func ReplaceWordInFile(path, oldWord, newWord string) error {

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
