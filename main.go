package main

import (
	"flag"
	"fmt"
	"github.com/raxigan/pcfy-my-mac/install"
	"log"
	"os"
)

func main() {

	verbose := flag.Bool("verbose", false, "Enable verbose mode")
	paramsFile := flag.String("params", "", "YAML file with installer parameters")
	flag.Parse()

	fileParams := install.FileParams{}

	if *paramsFile != "" {
		yamlStr, err := install.TextFromFile(*paramsFile)

		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		fileParams, _ = install.CollectYamlParams(yamlStr)
	}

	params := install.CollectParams(fileParams)

	handleError(
		install.RunInstaller(install.DefaultHomeDir(), install.DefaultCommander{Verbose: *verbose}, install.DefaultTimeProvider{}, params),
	)
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}
