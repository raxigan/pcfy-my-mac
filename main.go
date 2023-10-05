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
	"strings"
)

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose mode")
	showSampleYaml := flag.Bool("show-sample-yaml", false, "Show sample yaml config")
	paramsFile := flag.String("params", "", "Path to a YAML file containing installer parameters")
	flag.Parse()

	if *showSampleYaml {
		sampleYaml()
		os.Exit(0)
	}

	fileParams := param.FileParams{}

	if *paramsFile != "" {
		yamlStr, err := common.TextFromFile(*paramsFile)

		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		fileParams, err = param.CollectYamlParams(yamlStr)

		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}

	params := param.CollectParams(fileParams)

	handleError(
		cmd.Launch(install.DefaultHomeDir(), install.NewDefaultCommander(*verbose), install.DefaultTimeProvider{}, params),
	)
}

func sampleYaml() {
	yaml, _ := common.ReadFileFromEmbedFS("sample.yml")
	tabbedNewLine := "\n  "
	yaml = "  " + strings.ReplaceAll(yaml, "\n", tabbedNewLine)
	fmt.Println(fmt.Sprintf("\n  This is a sample YAML-based config. Copy it, adjust and then use in --param flag.\n\n%s", yaml))
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}
