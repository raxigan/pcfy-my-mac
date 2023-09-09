package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/bitfield/script"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
)

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExistsInHome(filename string) bool {
	usr, err := user.Current()
	if err != nil {
		return false
	}

	fullPath := filepath.Join(usr.HomeDir, filename)

	return fileExists(fullPath)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

const KarabinerConfigDir = ".config/karabiner"
const KarabinerConfig = KarabinerConfigDir + "/karabiner.json"

func makeSurvey(s MySurvey) string {
	if s.flagValue == "" {
		script.Exec("clear").Stdout()
		fmt.Println("App launcher...")

		appLauncher := ""

		prompt := &survey.Select{
			Message: s.description,
			Options: append(s.options, "None\n", "Quit"),
		}

		appLauncher = strings.TrimSpace(appLauncher)

		survey.AskOne(prompt, &appLauncher, survey.WithValidator(survey.Required))

		if appLauncher == "Quit" {
			fmt.Println("Quitting...")
			os.Exit(0)
		}

		return strings.ToLower(strings.TrimSpace(appLauncher))
	}

	return strings.ToLower(s.flagValue)
}

func main() {

	if !commandExists("brew") {
		fmt.Println("brew not installed. Exiting...")
	}

	if commandExists("jq") {
		fmt.Println("jq installed")
	} else {
		fmt.Println("jq not installed")
		script.Exec("brew install jq").Stdout()
	}

	if fileExistsInHome(KarabinerConfig) {
		fmt.Println("karabiner installed")
	} else {
		fmt.Println("Karabiner-Elements is not installed. Do you want to install karabiner-elements? [Y/n]")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		response := strings.TrimSpace(strings.ToLower(input))

		switch response {
		case "y", "yes":
			fmt.Println("yes")
			script.Exec("brew install --cask karabiner-elements")
		case "n", "no":
			fmt.Println("Karabiner-Elements required. Returning.")
			os.Exit(1)
		default:
			fmt.Println("Invalid input")
		}
	}

	appLauncherParam := flag.String("app-launcher", "", "Description for appLauncher")
	terminalParam := flag.String("terminal", "", "Description for terminalParam")
	kbTypeParam := flag.String("keyboard-type", "", "Description for terminalParam")

	flag.Parse()

	validateFlagValue(*appLauncherParam, []string{Spotlight.String(), Launchpad.String(), Alfred.String()})
	validateFlagValue(*terminalParam, []string{Default.String(), iTerm.String(), Warp.String()})
	validateFlagValue(*kbTypeParam, []string{PC.String(), Mac.String()})

	appLauncherSurvey := MySurvey{
		flagValue:   *appLauncherParam,
		description: "App Launcher:",
		options:     []string{Spotlight.String(), Launchpad.String(), Alfred.String()},
	}

	kbTypeSurvey := MySurvey{
		flagValue:   *kbTypeParam,
		description: "Your external keyboard type:",
		options:     []string{PC.String(), Mac.String()},
	}

	terminalSurvey := MySurvey{
		flagValue:   *terminalParam,
		description: "What is your terminal of choice:",
		options:     []string{Default.String(), iTerm.String(), Warp.String()},
	}

	app := makeSurvey(appLauncherSurvey)
	kbType := makeSurvey(kbTypeSurvey)
	term := makeSurvey(terminalSurvey)

	//currentTime := time.Now().Format("02-01-2006-15:04:05")
	//fmt.Println(currentTime)

	pwd, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()

	// do karabiner.json backup
	original := homeDir + "/" + KarabinerConfig
	//dest := homeDir + "/" + KarabinerConfigDir + "/karabiner-" + currentTime + ".json"
	dest := homeDir + "/" + KarabinerConfigDir + "/karabiner-new" + ".json"

	//fmt.Println(original)
	//fmt.Println(dest)

	script.Exec("cp " + original + " " + dest).Wait()

	//fmt.Println(pwd)

	// add karabiner profile

	//deleteProfileJq := "jq --arg PROFILE_NAME \"PC mode\" 'del(.profiles[] | select(.name == $PROFILE_NAME))' $KARABINER_CONFIG >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG"

	// delete existing profile
	oldProfileName := "PC mode"
	delete := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" 'del(.profiles[] | select(.name == \"%s\"))' %s >%s/INPUT.tmp && mv %s/INPUT.tmp %s", oldProfileName, oldProfileName, dest, pwd, pwd, dest)
	cmd1 := exec.Command("/bin/bash", "-c", delete)
	cmd1.Run()

	// add new karabiner profile
	cmdStr := fmt.Sprintf("jq '.profiles += $profile' %s --slurpfile profile %s/karabiner-elements-profile.json --indent 4 >%s/INPUT.tmp && mv %s/INPUT.tmp %s", dest, pwd, pwd, pwd, dest)
	cmd2 := exec.Command("/bin/bash", "-c", cmdStr)
	cmd2.Run()

	switch app {
	case "spotlight":
		fmt.Println("Applying spotlight rules...")
		applyRules("spotlight-rules.json", dest, pwd)
	case "launchpad":
		fmt.Println("Applying launchpad rules...")
		applyRules("launchpad-rules.json", dest, pwd)
	case "alfred":
		fmt.Println("Applying alfred rules...")
		applyRules("alfred-rules.json", dest, pwd)
	default:
		fmt.Println("Value is not A, B, or C")
	}

	switch kbType {
	case "pc":
		fmt.Println("Applying pc keyboard rules...")
	case "mac":
		fmt.Println("Applying mac keyboard rules...")
		prepareForMacKeyboard(dest, pwd)
	default:
		fmt.Println("Value is not A, B, or C")
	}

	switch term {
	case "default":
		fmt.Println("Applying apple terminal rules...")
		applyRules("terminal-rules.json", dest, pwd)
	case "iterm":
		fmt.Println("Applying iterm rules...")
		applyRules("iterm-rules.json", dest, pwd)
	case "warp":
		fmt.Println("Applying warp rules...")
		applyRules("warp-rules.json", dest, pwd)
	default:
		fmt.Println("Value is not A, B, or C")
	}
}

func applyRules(file string, karabinerConfig string, pwd string) {
	newProfile := "PC mode GOLANG"
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '(.profiles[] | select(.name == \"%s\") | .complex_modifications.rules) += $rules[].rules' %s --slurpfile rules %s/../karabiner-elements/%s --indent 4 >%s/INPUT.tmp && mv %s/INPUT.tmp %s", newProfile, newProfile, karabinerConfig, pwd, file, pwd, pwd, karabinerConfig)
	cmd1 := exec.Command("/bin/bash", "-c", jq)
	cmd1.Run()
}

func prepareForMacKeyboard(karabinerConfig string, pwd string) {
	newProfile := "PC mode GOLANG"
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '.profiles |= map(if .name == \"%s\" then walk(if type == \"object\" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' %s --indent 4 >%s/INPUT.tmp && mv %s/INPUT.tmp %s", newProfile, newProfile, karabinerConfig, pwd, pwd, karabinerConfig)
	cmd1 := exec.Command("/bin/bash", "-c", jq)
	cmd1.Run()
}

func validateFlagValue(value string, validValues []string) {

	v := toLowerSlice(validValues)

	if !slices.Contains(v, strings.ToLower(value)) {
		fmt.Println("Invalid flag value: " + value + ", valid values: " + strings.Join(v, ", "))
		os.Exit(1)
	}
}

func toLowerSlice(slice []string) []string {
	for i, s := range slice {
		slice[i] = strings.ToLower(s)
	}
	return slice
}
