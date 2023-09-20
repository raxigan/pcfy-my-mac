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

	os.Args = []string{"script_name", "--homedir=" + curr, "--terminal=warp", "--app-launcher=alfred", "--keyboard-type=pc"}

	i := NewInstallation().install()

	actual := i.karabinerConfigFile()
	expected := pwd + "/expected/karabiner-expected-warp-alfred-pc.json"
	equal, _ := areFilesEqual(actual, expected)
	if !equal {
		copyFile(i.karabinerConfigFile(), i.karabinerTestInvalidConfig("warp", "alfred", "pc"))
		copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
		removeFiles(i.karabinerConfigBackupFile())
		t.Fatalf("Files %s and %s are not equal", actual, expected)
	}

	removeFiles(i.karabinerConfigBackupFile())
	removeFiles(i.karabinerTestInvalidConfig("warp", "alfred", "pc"))
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

func TestInstallAllKeymaps(t *testing.T) {

	pwd, _ := os.Getwd()
	curr := pwd + "/homedir"
	os.Args = []string{"script_name", "--homedir=" + curr, "--terminal=warp", "--app-launcher=alfred", "--keyboard-type=pc", "--ides=all"}
	i := NewInstallation().install()

	verifyKeymaps(t, i.sourceKeymap(IntelliJ()), i.ideDirs(IntelliJ())[0])
	verifyKeymaps(t, i.sourceKeymap(PyCharm()), i.ideDirs(PyCharm())[0])
	verifyKeymaps(t, i.sourceKeymap(GoLand()), i.ideDirs(GoLand())[0])
	verifyKeymaps(t, i.sourceKeymap(Fleet()), i.ideDirs(Fleet())[0])
}

func verifyKeymaps(t *testing.T, srcKeymap, destKeymap string) {

	keymapsEqual, _ := areFilesEqual(srcKeymap, destKeymap)

	if !keymapsEqual {
		t.Fatalf("Files %s are not equal", []string{srcKeymap, destKeymap})
	}

	removeFiles(destKeymap)
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
