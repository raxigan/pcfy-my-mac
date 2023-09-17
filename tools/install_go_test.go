package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestInstallWarpAlfredPC(t *testing.T) {

	pwd, _ := os.Getwd()
	curr := pwd + "/homedir"

	os.Args = []string{"script_name", "--homedir=" + curr, "--terminal=warp", "--app-launcher=alfred", "--keyboard-type=pc", "--ides=intellij"}

	i := NewInstallation().install()

	actual := i.karabinerConfigFile()
	expected := pwd + "/expected/karabiner-expected-warp-alfred-pc.json"
	if !compareJSONFiles(t, actual, expected) {
		copyFile(i.karabinerConfigFile(), i.karabinerTestInvalidConfig("warp", "alfred", "pc"))
		copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
		removeFile(i.karabinerConfigBackupFile())
		t.Fatalf("JSON files %s and %s are not equal", actual, expected)
	}

	//restore karabiner initial config
	removeFile(i.karabinerConfigBackupFile())
	removeFile(i.karabinerTestInvalidConfig("warp", "alfred", "pc"))
	copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
}

func TestInstallItermSpotlightMac(t *testing.T) {

	pwd, _ := os.Getwd()
	curr := pwd + "/homedir"

	os.Args = []string{"script_name", "--homedir=" + curr, "--terminal=iterm", "--app-launcher=spotlight", "--keyboard-type=mac"}

	i := NewInstallation().install()

	actual := i.karabinerConfigFile()
	expected := pwd + "/expected/karabiner-expected-iterm-spotlight-mac.json"
	if !compareJSONFiles(t, actual, expected) {
		copyFile(i.karabinerConfigFile(), i.karabinerTestInvalidConfig("iterm", "spotlight", "mac"))
		copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
		removeFile(i.karabinerConfigBackupFile())
		t.Fatalf("JSON files %s and %s are not equal", actual, expected)
	}

	//restore karabiner initial config
	removeFile(i.karabinerConfigBackupFile())
	removeFile(i.karabinerTestInvalidConfig("iterm", "spotlight", "mac"))
	copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
}

func (i Installation) karabinerTestDefaultConfig() string {
	return i.karabinerConfigDir() + "/karabiner-default.json"
}

func (i Installation) karabinerTestInvalidConfig(terminal, appLauncher, keyboardType string) string {
	return i.karabinerConfigDir() + fmt.Sprintf("/karabiner-invalid-%s-%s-%s.json", terminal, appLauncher, keyboardType)
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
