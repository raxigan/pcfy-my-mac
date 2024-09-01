package test_utils

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func AssertEquals(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Fatalf("%s not equal to %s", actual, expected)
	}
}

func AssertFilesEqual(t *testing.T, actual, expected string) {

	AssertFileExists(t, actual)
	AssertFileExists(t, expected)

	actualFileContent := ReadFile(actual)
	expectedFileContent := ReadFile(expected)

	if strings.TrimSpace(actualFileContent) != strings.TrimSpace(expectedFileContent) {
		fmt.Println("=== ACTUAL:")
		fmt.Println(actualFileContent)
		fmt.Println("=== EXPECTED:")
		fmt.Println(expectedFileContent)
		t.Fatalf("Files %s are not equal", []string{actual, expected})
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
		t.Fatalf("Expected error, but there's none")
	}

	if !strings.Contains(err.Error(), Trim(expected)) {
		t.Fatalf(`EXPECTED: "%s" ACTUAL: "%s"`, expected, err)
	}
}

func AssertFileExists(t *testing.T, filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("File %s does not exist", filename)
		}
	}
}

func ReadFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(data)
}
