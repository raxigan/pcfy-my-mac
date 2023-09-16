package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestInstall(t *testing.T) {

	pwd, _ := os.Getwd()
	curr := pwd + "/homedir"
	fmt.Println(pwd)
	fmt.Println(curr)

	os.Args = []string{"script_name", "--homedir=" + curr, "--terminal=warp", "--app-launcher=alfred", "--keyboard-type=mac"}

	i := NewInstallation()
	i.install()

	fmt.Println(i.karabinerConfigDir())
	fmt.Println(i.karabinerConfigFile())
	fmt.Println(i.currentDir)

	actual := i.karabinerConfigFile()
	expected := pwd + "/expected/karabiner-expected.json"
	if !compareJSONFiles(t, actual, expected) {
		copyFile(i.karabinerConfigFile(), i.karabinerTestInvalidConfig())
		copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
		removeFile(i.karabinerConfigBackupFile())
		t.Fatalf("JSON files %s and %s are not equal", actual, expected)
	}

	//restore karabiner initial config
	removeFile(i.karabinerConfigBackupFile())
	removeFile(i.karabinerTestInvalidConfig())
	copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
}

func (p Installation) karabinerTestDefaultConfig() string {
	return p.karabinerConfigDir() + "/karabiner-default.json"
}

func (p Installation) karabinerTestInvalidConfig() string {
	return p.karabinerConfigDir() + "/karabiner-invalid.json"
}

func removeFile(name string) {
	err := os.Remove(name)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func readJSONFile(t *testing.T, path string) interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", path, err)
	}

	var parsed interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON from %s: %v", path, err)
	}
	return parsed
}

func compareJSONFiles(t *testing.T, pathA, pathB string) bool {
	a := readJSONFile(t, pathA)
	b := readJSONFile(t, pathB)
	return reflect.DeepEqual(a, b)
}

func copyFile(src, dst string) error {
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

	return destFile.Sync()
}
