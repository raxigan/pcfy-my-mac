package install_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/raxigan/macos-pc-mode/install"
	"io"
	"os"
	"testing"
)

func TestInstallWarpAlfredPC(t *testing.T) {

	os.Args = []string{"script_name", "--params=params/alfred-warp-pc.yml"}
	wd, _ := os.Getwd()
	i := install.RunInstaller(wd+"/homedir", install.MockCommander{})

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-warp-alfred-pc.json"
	equal, _ := areFilesEqual(actual, expected)
	if !equal {
		copyFile(i.KarabinerConfigFile(), karabinerTestInvalidConfig(i, "warp", "alfred", "pc"))
		copyFile(karabinerTestDefaultConfig(i), i.KarabinerConfigFile())
		removeFiles(i.KarabinerConfigBackupFile())
		t.Fatalf("Files %s and %s are not equal", actual, expected)
	}

	removeFiles(i.KarabinerConfigBackupFile())
	removeFiles(karabinerTestInvalidConfig(i, "warp", "alfred", "pc"))
}

func TestInstallItermSpotlightMac(t *testing.T) {

	os.Args = []string{"script_name", "--params=params/spotlight-iterm-mac.yml"}
	wd, _ := os.Getwd()
	i := install.RunInstaller(wd+"/homedir", install.MockCommander{})

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-iterm-spotlight-mac.json"
	equal, _ := areFilesEqual(actual, expected)
	if !equal {
		copyFile(i.KarabinerConfigFile(), karabinerTestInvalidConfig(i, "iterm", "spotlight", "mac"))
		copyFile(karabinerTestDefaultConfig(i), i.KarabinerConfigFile())
		removeFiles(i.KarabinerConfigBackupFile())
		t.Fatalf("JSON files %s and %s are not equal", actual, expected)
	}

	//restore karabiner initial config
	removeFiles(i.KarabinerConfigBackupFile())
	removeFiles(karabinerTestInvalidConfig(i, "iterm", "spotlight", "mac"))
	copyFile(karabinerTestDefaultConfig(i), i.KarabinerConfigFile())
}

func TestInstallAllKeymaps(t *testing.T) {

	os.Args = []string{"script_name", "--params=params/all-ides.yml"}
	wd, _ := os.Getwd()
	i := install.RunInstaller(wd+"/homedir", install.MockCommander{})

	assertEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJ()), i.IdeDirs(install.IntelliJ())[0])
	assertEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJ()), i.IdeDirs(install.IntelliJ())[1])
	assertEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJCE()), i.IdeDirs(install.IntelliJCE())[0])
	assertEqual(t, "../configs/"+i.SourceKeymap(install.GoLand()), i.IdeDirs(install.GoLand())[0])
	assertEqual(t, "../configs/"+i.SourceKeymap(install.Fleet()), i.IdeDirs(install.Fleet())[0])

	removeFiles(i.KarabinerConfigBackupFile())
	removeFiles(karabinerTestInvalidConfig(i, "iterm", "spotlight", "mac"))
	copyFile(karabinerTestDefaultConfig(i), i.KarabinerConfigFile())
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

func karabinerTestDefaultConfig(i install.Installation) string {
	return i.KarabinerConfigDir() + "/karabiner-default.json"
}

func karabinerTestInvalidConfig(i install.Installation, terminal, appLauncher, keyboardType string) string {
	return i.KarabinerConfigDir() + fmt.Sprintf("/karabiner-invalid-%s-%s-%s.json", terminal, appLauncher, keyboardType)
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

	return nil
}
