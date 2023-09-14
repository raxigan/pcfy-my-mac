package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/bitfield/script"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

const KarabinerConfigDir = ".config/karabiner"
const KarabinerConfig = KarabinerConfigDir + "/karabiner.json"

const branchName = "feature/installation_script"

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

	pwd, _ := os.Getwd()
	homeDirDefault, _ := os.UserHomeDir()

	homeDirFlagValue := flag.String("homedir", homeDirDefault, "Home directory path")
	appLauncherParam := flag.String("app-launcher", "", "Description for appLauncher")
	terminalParam := flag.String("terminal", "", "Description for terminalParam")
	kbTypeParam := flag.String("keyboard-type", "", "Description for terminalParam")

	flag.Parse()
	homeDir := *homeDirFlagValue

	validateFlagValue(*appLauncherParam, []string{Spotlight.String(), Launchpad.String(), Alfred.String()})
	validateFlagValue(*terminalParam, []string{Default.String(), iTerm.String(), Warp.String()})
	validateFlagValue(*kbTypeParam, []string{PC.String(), Mac.String()})

	if shouldBeInstalled("jq", "jq", true, true, false) {

	}

	if shouldBeInstalled("Karabiner-Elements", "Karabiner-Elements", false, true, true) {

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

		// do karabiner.json backup
		original := homeDir + "/" + KarabinerConfig
		//dest := homeDir + "/" + KarabinerConfigDir + "/karabiner-" + currentTime + ".json"
		dest := homeDir + "/" + KarabinerConfigDir + "/karabiner-new" + ".json"

		//fmt.Println(original)
		//fmt.Println(dest)

		script.Exec("cp " + original + " " + dest).Wait()

		//fmt.Println(pwd)

		// add karabiner profile

		//deleteProfileJq := "jq --arg PROFILE_NAME \"PC mode\" 'del(.profiles[] | select(.installAltTab == $PROFILE_NAME))' $KARABINER_CONFIG >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG"

		// delete existing profile
		oldProfileName := "PC mode"
		delete := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" 'del(.profiles[] | select(.installAltTab == \"%s\"))' %s >%s/INPUT.tmp && mv %s/INPUT.tmp %s", oldProfileName, oldProfileName, dest, pwd, pwd, dest)
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

		runWithOutput("clear")

		// add a flag for it
		var ideKeymaps []string
		prompt := &survey.MultiSelect{
			Message: "IDE keymaps to install:",
			Options: []string{"IntelliJ IDEA Ultimate", "PyCharm Community Edition"},
		}
		survey.AskOne(prompt, &ideKeymaps)

		if contains(ideKeymaps, "IntelliJ IDEA Ultimate") {
			installIdeKeymap("idea", "IntelliJ IDEA Ultimate")
		}
	}

	if shouldBeInstalled("Rectangle", "Rectangle", false, false, true) {
		runWithOutput("killall Rectangle")

		rectJson := "RectangleConfig.json"

		c := "cp " + pwd + "/../rectangle/" + rectJson + " \"" + homeDir + "/Library/Application Support/Rectangle/RectangleConfig.json\""
		cmdMkdir := "mkdir -p " + "\"" + homeDir + "/Library/Application Support/Rectangle\""
		runWithOutput(cmdMkdir)
		runWithOutput(c)

		runWithOutput("open -a Rectangle")
	}

	if shouldBeInstalled("Alt-Tab", "AltTab", false, false, true) {
		runWithOutput("killall AltTab")

		jsonName := "Settings.json"
		jsonFile := pwd + "/../alt-tab/" + jsonName

		fileContent, _ := os.ReadFile(jsonFile)

		var settings map[string]interface{}

		eee := json.Unmarshal(fileContent, &settings)

		if eee != nil {
			fmt.Println(eee)
		}

		altTabPlist := homeDir + "/Library/preferences/com.lwouis.alt-tab-macos.plist"
		fmt.Println(altTabPlist)

		for key, value := range settings {

			str := fmt.Sprintf("%v", value)
			sprintf := fmt.Sprintf("defaults write %s '%s' '%s'", altTabPlist, key, str)

			if key == "blacklist" {
				sprintf = fmt.Sprintf(`defaults write %s %s "'%s'"`, altTabPlist, key, strings.ReplaceAll(str, `"`, `\"`))
			}

			runWithOutput(sprintf)
		}

		runWithOutput("defaults read com.lwouis.alt-tab-macos")
		runWithOutput("open -a AltTab")
	}

	fmt.Println("SUCCESS")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func shouldBeInstalled(appName string, appFile string, isCommand bool, isRequired bool, isCask bool) bool {

	exists := false

	if isCommand {
		exists = commandExists(appName)
	} else {
		exists = fileExists("/Applications/" + appFile + ".app")
	}

	if exists {
		return true
	}

	runWithOutput("clear")
	installApp := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Do you want to install %s?", appName),
	}
	survey.AskOne(prompt, &installApp)

	if installApp {
		fmt.Println(fmt.Sprintf("Installing %s...", appName))

		brewCommand := fmt.Sprintf("brew install %s", strings.ToLower(appName))

		if isCask {
			brewCommand = fmt.Sprintf("brew install --cask %s", strings.ToLower(appName))
		}

		script.Exec(brewCommand).Stdout()
	} else {
		if isRequired {
			fmt.Println(fmt.Sprintf("%s is required to proceed. Quitting...", appName))
			os.Exit(0)
		}
	}

	return installApp
}

func runWithOutput(cmd string) {
	run(cmd, true)
}

func runNoOutput(cmd string) {
	run(cmd, false)
}

func run(cmd string, out bool) {

	fmt.Println("Running: " + cmd)

	if out {
		output, err := exec.Command("/bin/bash", "-c", cmd).Output()

		if err != nil {
			fmt.Println("Error executing command: "+cmd+"\n", err)
		}

		fmt.Print(string(output))

	} else {
		exec.Command("/bin/bash", "-c", cmd)
	}

}

func installIdeKeymap(scriptName string, ideFullName string) {
	homeDir, _ := os.UserHomeDir()
	content, err := os.ReadFile(homeDir + "/Library/Application Support/JetBrains/Toolbox/scripts/" + scriptName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	r := regexp.MustCompile(`ch-0/(\d+\.\d+\.\d+)`)
	matches := r.FindSubmatch(content)
	if matches != nil {

		version := string(matches[1])

		fmt.Println("Installin XWin plugin for " + version + " " + ideFullName)

		cmd := fmt.Sprintf("open -na \"%s.app\" --args installPlugins com.intellij.plugins.xwinkeymap", ideFullName)

		exec.Command("/bin/bash", "-c", cmd)

		configs := homeDir + "/Library/Application Support/JetBrains"

		dirPath := homeDir + "/Library/Caches/Jetbrains"
		intelliJDir, err := findIntelliJDir(dirPath, version)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		if intelliJDir != "" {

			fmt.Printf("Found directory: %s\n", intelliJDir)

			keymapDir := configs + "/" + filepath.Base(intelliJDir)
			keymapFileName := strings.ReplaceAll(strings.ToLower(ideFullName), " ", "-")

			// if there is a local dir with keymaps, then take it from there
			cmd := fmt.Sprintf("curl --silent -o \"%s/%s.xml\" https://raw.githubusercontent.com/raxigan/macos-pc-mode/%s/keymaps/\"%s\".xml", keymapDir, keymapFileName, branchName, keymapFileName)
			cmd1 := exec.Command("/bin/bash", "-c", cmd)
			cmd1.Run()
		} else {
			fmt.Println("No matching directory found.")
		}

	} else {
		fmt.Println("Version not found.")
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

func findIntelliJDir(path string, version string) (string, error) {
	var resultDir string

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.Contains(info.Name(), "IntelliJ") {
			appInfoPath := filepath.Join(currentPath, ".appinfo")
			content, err := os.ReadFile(appInfoPath)
			if err == nil && strings.Contains(string(content), "app.build.number="+version) {
				resultDir = currentPath
				return fmt.Errorf("found")
			}
		}
		return nil
	})

	if err != nil && err.Error() == "found" {
		err = nil
	}

	return resultDir, err
}
