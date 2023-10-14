package main

import (
	"flag"
	"fmt"
	"github.com/raxigan/pcfy-my-mac/cmd"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/install"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"os"
	"time"
)

var (
	version   string
	buildTime string
	commit    string
)

func main() {

	showVersion := flag.Bool("version", false, "Show version information")
	verbose := flag.Bool("verbose", false, "Enable verbose mode")
	showSampleYaml := flag.Bool("show-sample-yaml", false, "Show sample yaml config")
	paramsFile := flag.String("params", "", "Path to a YAML file containing installer parameters")
	flag.Parse()

	handleVersionFlag(showVersion)
	handleSampleYamlFlag(showSampleYaml)

	commander := install.NewDefaultCommander(*verbose)
	commander.Run("clear")
	params, err := param.CollectParams(*paramsFile)

	handleError(err, commander)
	handleError(cmd.Launch(
		install.DefaultHomeDir(),
		commander,
		install.DefaultTimeProvider{},
		params,
	), commander,
	)
}

func handleSampleYamlFlag(showSampleYaml *bool) {
	if *showSampleYaml {
		printSampleYaml()
		os.Exit(0)
	}
}

func handleVersionFlag(showVersion *bool) {
	if *showVersion {
		buildTime = time.Now().Format("2006-01-02 03:04:05 PM")
		fmt.Printf("Version: %s\nBuild Time: %s\nCommit: %s\n", version, buildTime, commit)
		os.Exit(0)
	}
}

func printSampleYaml() {
	yaml, _ := common.ReadFileFromEmbedFS("sample.yml")
	fmt.Println(fmt.Sprintf("This is a sample YAML-based config. Copy it, adjust and then use in --param flag.\n\n%s", yaml))
}

func handleError(err error, commander install.Commander) {
	if err != nil {
		commander.TryLog(install.ErrMsg, fmt.Sprintf("%s", err))
		os.Exit(1)
	}
}
