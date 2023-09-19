package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestInstallWarpAlfredPC(t *testing.T) {

	pwd, _ := os.Getwd()
	curr := pwd + "/homedir"

	os.Args = []string{"script_name", "--homedir=" + curr, "--terminal=warp", "--app-launcher=alfred", "--keyboard-type=pc", "--ides=intellij"}

	i := NewInstallation().install()

	actual := i.karabinerConfigFile()
	expected := pwd + "/expected/karabiner-expected-warp-alfred-pc.json"
	equal, _ := areFilesEqual(actual, expected)
	if !equal {
		srcKeymap := i.currentDir + "/../keymaps/intellij-idea-ultimate.xml"
		destKeymap1 := i.applicationSupportDir() + "/JetBrains/IntelliJIdea2023.1/keymaps/intellij-idea-ultimate.xml"
		destKeymap2 := i.applicationSupportDir() + "/JetBrains/IntelliJIdea2023.2/keymaps/intellij-idea-ultimate.xml"

		keymapsEqual, _ := areFilesEqual(srcKeymap, destKeymap1, destKeymap2)

		if !keymapsEqual {
			t.Fatalf("Files %s are not equal", []string{srcKeymap, destKeymap1, destKeymap2})
		}

		removeFiles(destKeymap1, destKeymap2)
	}
}

func TestInstallItermSpotlightMac(t *testing.T) {

	pwd, _ := os.Getwd()
	curr := pwd + "/homedir"

	os.Args = []string{"script_name", "--homedir=" + curr, "--terminal=iterm", "--app-launcher=spotlight", "--keyboard-type=mac"}

	i := NewInstallation().install()

	actual := i.karabinerConfigFile()
	expected := pwd + "/expected/karabiner-expected-iterm-spotlight-mac.json"
	equal, _ := areFilesEqual(actual, expected)
	if !equal {
		copyFile(i.karabinerConfigFile(), i.karabinerTestInvalidConfig("iterm", "spotlight", "mac"))
		copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
		removeFiles(i.karabinerConfigBackupFile())
		t.Fatalf("JSON files %s and %s are not equal", actual, expected)
	}

	//restore karabiner initial config
	removeFiles(i.karabinerConfigBackupFile())
	removeFiles(i.karabinerTestInvalidConfig("iterm", "spotlight", "mac"))
	copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
}

func (i Installation) karabinerTestDefaultConfig() string {
	return i.karabinerConfigDir() + "/karabiner-default.json"
}

func (i Installation) karabinerTestInvalidConfig(terminal, appLauncher, keyboardType string) string {
	return i.karabinerConfigDir() + fmt.Sprintf("/karabiner-invalid-%s-%s-%s.json", terminal, appLauncher, keyboardType)
}

func removeFiles(paths ...string) error {
	for _, path := range paths {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to remove %s: %v", path, err)
		}
	}
	return nil
}

//func readJSONFile(t *testing.T, path string) interface{} {
//	data, err := os.ReadFile(path)
//	if err != nil {
//		t.Fatalf("Failed to read %s: %v", path, err)
//	}
//
//	var parsed interface{}
//	if err := json.Unmarshal(data, &parsed); err != nil {
//		t.Fatalf("Failed to parse JSON from %s: %v", path, err)
//	}
//	return parsed
//}
//
//func compareJSONFiles(t *testing.T, pathA, pathB string) bool {
//	a := readJSONFile(t, pathA)
//	b := readJSONFile(t, pathB)
//	return reflect.DeepEqual(a, b)
//}

func areFilesEqual(paths ...string) (bool, error) {
	if len(paths) < 2 {
		return false, errors.New("at least two paths are required for comparison")
	}

	reference, err := os.ReadFile(paths[0])
	if err != nil {
		return false, err
	}

	for _, path := range paths[1:] {
		content, err := os.ReadFile(path)
		if err != nil {
			return false, err
		}
		if !bytes.Equal(reference, content) {
			return false, nil
		}
	}

	return true, nil
}
