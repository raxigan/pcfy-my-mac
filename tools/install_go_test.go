package main

import (
	"fmt"
	"os"
	"testing"
)

func TestInstall(t *testing.T) {

	pwd, _ := os.Getwd()
	curr := pwd + "/homedir"
	fmt.Println(pwd)
	fmt.Println(curr)

	os.Args = []string{"script_name", "--homedir=" + curr, "--terminal=warp", "--app-launcher=alfred", "--keyboard-type=mac"}
	main()
}
