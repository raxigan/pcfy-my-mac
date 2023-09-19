package main

import (
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/bitfield/script"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

type Installation struct {
	homeDir    string
	currentDir string

	appLauncher  string
	terminal     string
	keyboardType string
	ides         []string
	askForIdes   bool

	profileName      string
	installationTime time.Time
}

type idesFlag struct {
	provided bool
	value    string
}

func (sf *idesFlag) String() string {
	return sf.value
}

func (sf *idesFlag) Set(v string) error {
	sf.value = v
	sf.provided = true
	return nil
}

func NewInstallation() *Installation {

	pwd, _ := os.Getwd()
	homeDirDefault, _ := os.UserHomeDir()

	homeDirFlagValue := flag.String("homedir", homeDirDefault, "Home directory path")
	appLauncherParam := flag.String("app-launcher", "", "Description for appLauncher")
	terminalParam := flag.String("terminal", "", "Description for terminalParam")
	kbTypeParam := flag.String("keyboard-type", "", "Description for kbTypeParam")
	var ids idesFlag
	flag.Var(&ids, "ides", "Description for ides")

	flag.Parse()

	homeDir := *homeDirFlagValue

	validateFlagValue("app-launcher", *appLauncherParam, []string{Spotlight.String(), Launchpad.String(), Alfred.String()})
	validateFlagValue("terminal", *terminalParam, []string{Default.String(), iTerm.String(), Warp.String()})
	validateFlagValue("keyboard-type", *kbTypeParam, []string{PC.String(), Mac.String()})
	validateFlagValue("ides", ids.value, []string{
		strings.ToLower(IntelliJ().flag),
		strings.ToLower(PyCharm().flag),
		strings.ToLower(GoLand().flag),
		strings.ToLower(Fleet().flag),
	},
	)

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var ides []string

	askForIdes := true

	if !ids.provided {
		ides = []string{}
	} else {

		askForIdes = false

		var supportedIDEs map[string]string
		supportedIDEs = make(map[string]string)

		supportedIDEs[IntelliJ().flag] = IntelliJ().fullName
		supportedIDEs[PyCharm().flag] = PyCharm().fullName
		supportedIDEs[GoLand().flag] = GoLand().fullName
		supportedIDEs[Fleet().flag] = Fleet().fullName

		for _, e := range strings.Split(ids.value, ",") {
			if e != "" {
				ides = append(ides, supportedIDEs[e])
			}
		}
	}

	return &Installation{
		homeDir:          homeDir,
		currentDir:       pwd,
		appLauncher:      *appLauncherParam,
		terminal:         *terminalParam,
		keyboardType:     *kbTypeParam,
		ides:             ides,
		askForIdes:       askForIdes,
		profileName:      "PC mode GOLANG",
		installationTime: time.Now(),
	}
}

func (i Installation) karabinerConfigDir() string {
	return i.homeDir + "/.config/karabiner"
}

func (i Installation) karabinerConfigFile() string {
	return i.karabinerConfigDir() + "/karabiner.json"
}

func (i Installation) karabinerConfigBackupFile() string {
	currentTime := i.installationTime.Format("02-01-2023-15:04:05")
	return i.karabinerConfigDir() + "/karabiner-" + currentTime + ".json"
}

func (i Installation) applicationSupportDir() string {
	return i.homeDir + "/Library/Application Support"
}

func (i Installation) preferencesDir() string {
	return i.homeDir + "/Library/preferences"
}

func (i Installation) toolboxScriptsDir() string {
	return i.applicationSupportDir() + "/JetBrains/Toolbox/scripts"
}

const branchName = "feature/installation_script"

func makeSurvey(s MySurvey) string {
	script.Exec("clear").Stdout()

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

func makeMultiSelect(s MySurvey) []string {

	script.Exec("clear").Stdout()

	var appLauncher []string

	prompt := &survey.MultiSelect{
		Message: s.description,
		Options: s.options,
	}

	survey.AskOne(prompt, &appLauncher)

	return appLauncher
}

func main() {
	NewInstallation().install()
}

func (i Installation) install() Installation {

	if shouldBeInstalled("jq", "jq", true, true, false) {

	}

	if shouldBeInstalled("Karabiner-Elements", "Karabiner-Elements", false, true, true) {

		appLauncherSurvey := MySurvey{
			description: "App Launcher:",
			options:     []string{Spotlight.String(), Launchpad.String(), Alfred.String()},
		}

		kbTypeSurvey := MySurvey{
			description: "Your external keyboard type:",
			options:     []string{PC.String(), Mac.String()},
		}

		terminalSurvey := MySurvey{
			description: "What is your terminal of choice:",
			options:     []string{Default.String(), iTerm.String(), Warp.String()},
		}

		app := i.appLauncher
		term := i.terminal
		kbType := i.keyboardType

		if i.appLauncher == "" {
			app = makeSurvey(appLauncherSurvey)
		}
		if i.keyboardType == "" {
			kbType = makeSurvey(kbTypeSurvey)
		}

		if i.terminal == "" {
			term = makeSurvey(terminalSurvey)
		}

		run("killall Karabiner-Elements")

		// do karabiner.json backup
		original := i.karabinerConfigFile()
		backupDest := i.karabinerConfigBackupFile()

		script.Exec("cp " + original + " " + backupDest).Wait()

		// delete existing profile
		deleteProfileJqCmd := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" 'del(.profiles[] | select(.name == \"%s\"))' %s >%s/INPUT.tmp && mv %s/INPUT.tmp %s", i.profileName, i.profileName, i.karabinerConfigFile(), i.currentDir, i.currentDir, i.karabinerConfigFile())
		run(deleteProfileJqCmd)

		// add new karabiner profile
		addProfileJqCmd := fmt.Sprintf("jq '.profiles += $profile' %s --slurpfile profile %s/karabiner-elements-profile.json --indent 4 >%s/INPUT.tmp && mv %s/INPUT.tmp %s", i.karabinerConfigFile(), i.currentDir, i.currentDir, i.currentDir, i.karabinerConfigFile())
		run(addProfileJqCmd)

		// rename the profile
		renameJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name == \"PROFILE_NAME\" then .name = \"%s\" else . end)' %s > INPUT.tmp && mv INPUT.tmp %s", i.profileName, i.karabinerConfigFile(), i.karabinerConfigFile())
		run(renameJqCmd)

		// unselect other profiles
		unselectJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name != \"%s\" then .selected = false else . end)' %s > INPUT.tmp && mv INPUT.tmp %s", i.profileName, i.karabinerConfigFile(), i.karabinerConfigFile())
		run(unselectJqCmd)

		applyRules(i, "main-rules.json")
		applyRules(i, "finder-rules.json")

		switch app {
		case "spotlight":
			fmt.Println("Applying spotlight rules...")
			applyRules(i, "spotlight-rules.json")
		case "launchpad":
			fmt.Println("Applying launchpad rules...")
			applyRules(i, "launchpad-rules.json")
		case "alfred":
			{
				fmt.Println("Applying alfred rules...")
				applyRules(i, "alfred-rules.json")
				c1 := fmt.Sprintf("find '%s' -type d -name \"hotkey\" -exec cp %s {} \\;", i.applicationSupportDir()+"/Alfred/Alfred.alfredpreferences/preferences/local", i.currentDir+"/../alfred4/prefs.plist")
				run(c1)
			}
		default:
			fmt.Println("Value is not A, B, or C")
		}

		switch kbType {
		case "pc":
			fmt.Println("Applying pc keyboard rules...")
		case "mac":
			fmt.Println("Applying mac keyboard rules...")
			prepareForMacKeyboard(i, i.karabinerConfigFile(), i.currentDir)
		default:
			fmt.Println("Value is not A, B, or C")
		}

		switch term {
		case "default":
			fmt.Println("Applying apple terminal rules...")
			applyRules(i, "terminal-rules.json")
		case "iterm":
			fmt.Println("Applying iterm rules...")
			applyRules(i, "iterm-rules.json")
		case "warp":
			fmt.Println("Applying warp rules...")
			applyRules(i, "warp-rules.json")
		default:
			fmt.Println("Value is not A, B, or C")
		}

		run("open -a Karabiner-Elements")
		run("clear")

		var supportedIDEs map[string]IDE
		supportedIDEs = make(map[string]IDE)

		supportedIDEs[IntelliJ().fullName] = IntelliJ()
		supportedIDEs[PyCharm().fullName] = PyCharm()
		supportedIDEs[GoLand().fullName] = GoLand()
		supportedIDEs[Fleet().fullName] = Fleet()

		var ideOptions []string

		for _, value := range supportedIDEs {
			ideOptions = append(ideOptions, value.fullName)
		}

		idesToInstall := i.ides

		if i.askForIdes {
			ideSurvey := MySurvey{
				description: "IDE keymaps to install:",
				options:     ideOptions,
			}

			idesToInstall = makeMultiSelect(ideSurvey)
		}

		for _, name := range idesToInstall {
			installIdeKeymap(supportedIDEs[name], i)
		}
	}

	if shouldBeInstalled("Rectangle", "Rectangle", false, false, true) {
		run("killall Rectangle")

		xmlFile := i.currentDir + "/../rectangle/Settings.xml"
		rectanglePlist := i.preferencesDir() + "/com.knollsoft.Rectangle.plist"

		plutilCmd := fmt.Sprintf("plutil -convert binary1 -o %s %s", rectanglePlist, xmlFile)
		run(plutilCmd)

		run("defaults read com.knollsoft.Rectangle")
		run("open -a Rectangle")
	}

	if shouldBeInstalled("Alt-Tab", "AltTab", false, false, true) {
		run("killall AltTab")

		xmlFile := i.currentDir + "/../alt-tab/Settings.xml"
		altTabPlist := i.preferencesDir() + "/com.lwouis.alt-tab-macos.plist"

		plutilCmd := fmt.Sprintf("plutil -convert binary1 -o %s %s", altTabPlist, xmlFile)
		run(plutilCmd)

		run("defaults read com.lwouis.alt-tab-macos")
		run("open -a AltTab")
	}

	fmt.Println("SUCCESS")

	return i
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

	run("clear")
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

func run(cmd string) {

	fmt.Println("Running: " + cmd)

	output, err := exec.Command("/bin/bash", "-c", cmd).Output()

	if err != nil {
		fmt.Println("Error executing command: "+cmd+"\n", err)
	}

	fmt.Print(string(output))
}

func installIdeKeymap(ide IDE, installation Installation) {

	if ide.requiresPlugin {
		cmd := fmt.Sprintf("open -na \"%s.app\" --args installPlugins com.intellij.plugins.xwinkeymap", ide.fullName)
		run(cmd)
	}

	sourceFile := installation.currentDir + "/../keymaps/" + ide.srcKeymapsFile

	err := filepath.Walk(installation.homeDir+ide.parentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasPrefix(info.Name(), ide.dir) {
			destDir := filepath.Join(path, ide.keymapsDir)
			destFilePath := filepath.Join(destDir, ide.destKeymapsFile)
			os.Mkdir(destDir, 0755)
			err := copyFile(sourceFile, destFilePath)
			if err != nil {
				fmt.Printf("Error copying to %s: %v\n", destFilePath, err)
			} else {
				fmt.Printf("Successfully copied to %s\n", destFilePath)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error while walking the path: %v\n", err)
	}
}

func applyRules(i Installation, file string) {
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '(.profiles[] | select(.name == \"%s\") | .complex_modifications.rules) += $rules[].rules' %s --slurpfile rules %s/../karabiner-elements/%s --indent 4 >%s/INPUT.tmp && mv %s/INPUT.tmp %s", i.profileName, i.profileName, i.karabinerConfigFile(), i.currentDir, file, i.currentDir, i.currentDir, i.karabinerConfigFile())
	cmd1 := exec.Command("/bin/bash", "-c", jq)
	cmd1.Run()
}

func prepareForMacKeyboard(i Installation, karabinerConfig string, pwd string) {
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '.profiles |= map(if .name == \"%s\" then walk(if type == \"object\" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' %s --indent 4 >%s/INPUT.tmp && mv %s/INPUT.tmp %s", i.profileName, i.profileName, karabinerConfig, pwd, pwd, karabinerConfig)
	cmd1 := exec.Command("/bin/bash", "-c", jq)
	cmd1.Run()
}

func validateFlagValue(flag, value string, validValues []string) {

	if value != "" {
		v := toLowerSlice(validValues)

		if !slices.Contains(v, strings.ToLower(value)) {
			fmt.Println("Invalid flag " + flag + " value: " + value + ", valid values: " + strings.Join(v, ", "))
			os.Exit(1)
		}
	}
}

func toLowerSlice(slice []string) []string {
	for i, s := range slice {
		slice[i] = strings.ToLower(s)
	}
	return slice
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
