package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestInstallWarpAlfredPC(t *testing.T) {

	os.Args = []string{"script_name", "--terminal=warp", "--app-launcher=alfred", "--keyboard-type=pc"}
	wd, _ := os.Getwd()
	i := runInstaller(wd+"/homedir", MockCommander{})

	actual := i.karabinerConfigFile()
	expected := "expected/karabiner-expected-warp-alfred-pc.json"
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

	os.Args = []string{"script_name", "--terminal=iterm", "--app-launcher=spotlight", "--keyboard-type=mac"}
	wd, _ := os.Getwd()
	i := runInstaller(wd+"/homedir", MockCommander{})

	actual := i.karabinerConfigFile()
	expected := "expected/karabiner-expected-iterm-spotlight-mac.json"
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

	os.Args = []string{"script_name", "--terminal=warp", "--app-launcher=alfred", "--keyboard-type=pc", "--ides=all"}
	wd, _ := os.Getwd()
	i := runInstaller(wd+"/homedir", MockCommander{})

	assertEqual(t, i.sourceKeymap(IntelliJ()), i.ideDirs(IntelliJ())[0])
	assertEqual(t, i.sourceKeymap(IntelliJ()), i.ideDirs(IntelliJ())[1])
	assertEqual(t, i.sourceKeymap(IntelliJCE()), i.ideDirs(IntelliJCE())[0])
	assertEqual(t, i.sourceKeymap(GoLand()), i.ideDirs(GoLand())[0])
	assertEqual(t, i.sourceKeymap(Fleet()), i.ideDirs(Fleet())[0])

	removeFiles(i.karabinerConfigBackupFile())
	removeFiles(i.karabinerTestInvalidConfig("iterm", "spotlight", "mac"))
	copyFile(i.karabinerTestDefaultConfig(), i.karabinerConfigFile())
}

func assertEqual(t *testing.T, srcKeymap, destKeymap string) {

	keymapsEqual, err := compareFilesBySHASum(srcKeymap, destKeymap)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

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

func computeSHA256(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func compareFilesBySHASum(file1, file2 string) (bool, error) {
	sha1, err := computeSHA256(file1)

	if err != nil {
		return false, err
	}

	sha2, err := computeSHA256(file2)

	if err != nil {
		return false, err
	}

	return sha1 == sha2, nil
}
