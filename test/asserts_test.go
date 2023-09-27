package install_test

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func AssertFilesEqual(t *testing.T, path1, path2 string) {

	keymapsEqual, err := compareFilesBySHASum(path1, path2)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if !keymapsEqual {
		t.Fatalf("Files %s are not equal", []string{path1, path2})
	}
}

func AssertSlicesEqual(t *testing.T, slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		join1 := strings.Join(slice1, "\n")
		join2 := strings.Join(slice2, "\n")
		t.Fatalf("Slices not equal: \n\n=== ACTUAL:\n%s  \n\n=== EXPECTED:\n%s", join1, join2)
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			join1 := strings.Join(slice1, "\n")
			join2 := strings.Join(slice2, "\n")
			t.Fatalf("Slices not equal: \n\n=== ACTUAL:\n%s  \n\n=== EXPECTED:\n%s", join1, join2)
		}
	}

	return true
}

func AssertErrorContains(t *testing.T, err error, expected string) {
	if err == nil {
		t.Fatalf("")
	}

	if !strings.Contains(err.Error(), yaml(expected)) {
		t.Fatalf(`EXPECTED: "%s" ACTUAL: "%s"`, err, expected)
	}
}

func computeSHA256(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error: %s", err)
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
