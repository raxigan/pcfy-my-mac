package main

import (
	"flag"
	"fmt"
	"github.com/raxigan/pcfy-my-mac/cmd"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/install"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"log"
	"os"
)

func main() {

	verbose := flag.Bool("verbose", false, "Enable verbose mode")
	paramsFile := flag.String("params", "", "YAML file with installer parameters")
	flag.Parse()

	fileParams := param.FileParams{}

	if *paramsFile != "" {
		yamlStr, err := common.TextFromFile(*paramsFile)

		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		fileParams, _ = param.CollectYamlParams(yamlStr)
	}

	params := param.CollectParams(fileParams)

	handleError(
		cmd.RunInstaller(install.DefaultHomeDir(), install.NewDefaultCommander(*verbose), install.DefaultTimeProvider{}, params),
	)
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}
