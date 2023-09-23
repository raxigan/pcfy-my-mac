package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/bitfield/script"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

//go:embed configs/*
var configs embed.FS

func copyFileFromEmbedFS(src, dst string) error {
	data, _ := fs.ReadFile(configs, src)
	return os.WriteFile(dst, data, 0644)
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

type Params struct {
	appLauncher  string
	terminal     string
	keyboardType string
	ides         *[]IDE
}

type Installation struct {
	homeDir    string
	currentDir string

	flagParams Params

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

	validateFlagValue("app-launcher", *appLauncherParam, []string{Spotlight.String(), Launchpad.String(), Alfred.String(), "none"})
	validateFlagValue("terminal", *terminalParam, []string{Default.String(), iTerm.String(), Warp.String(), "none"})
	validateFlagValue("keyboard-type", *kbTypeParam, []string{PC.String(), Mac.String(), "none"})
	validateFlagValue("ides", ids.value, append(IdeKeymapsFlags(), []string{"none", "all"}...))

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var ides *[]IDE

	if !ids.provided {
		ides = nil
	} else if ids.value == "all" {
		ides = &IDEKeymaps
	} else {

		var idesFromFlags []IDE

		for _, e := range strings.Split(ids.value, ",") {
			if e != "" {
				byFlag, _ := IdeKeymapByFlag(e)
				idesFromFlags = append(idesFromFlags, byFlag)
			}
		}

		ides = &idesFromFlags
	}

	return &Installation{
		homeDir:    homeDir,
		currentDir: pwd,
		flagParams: Params{
			appLauncher:  *appLauncherParam,
			terminal:     *terminalParam,
			keyboardType: *kbTypeParam,
			ides:         ides,
		},
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

func (i Installation) karabinerComplexModificationsDir() string {
	return i.homeDir + "/.config/karabiner/assets/complex_modifications"
}

func (i Installation) applicationSupportDir() string {
	return i.homeDir + "/Library/Application Support"
}

func (i Installation) preferencesDir() string {
	return i.homeDir + "/Library/preferences"
}

func (i Installation) sourceKeymap(ide IDE) string {
	return "configs/keymaps/" + ide.srcKeymapsFile
}

func makeSurvey(s MySurvey) string {
	script.Exec("clear").Stdout()

	appLauncher := ""

	prompt := &survey.Select{
		Message: s.message,
		Options: append(s.options, "None"),
	}

	appLauncher = strings.TrimSpace(appLauncher)

	survey.AskOne(prompt, &appLauncher, survey.WithValidator(survey.Required))

	if appLauncher == "Quit" {
		fmt.Println("Quitting...")
		os.Exit(0)
	}

	return strings.ToLower(strings.TrimSpace(appLauncher))
}

func makeMultiSelect(s survey.MultiSelect) []string {

	script.Exec("clear").Stdout()

	var appLauncher []string

	survey.AskOne(&s, &appLauncher)

	return appLauncher
}

func runInstaller() Installation {
	installation := NewInstallation()
	params := installation.collectParams()
	return installation.install(params)
}

func main() {
	runInstaller()
}

func (i Installation) collectParams() Params {

	app := i.flagParams.appLauncher
	term := i.flagParams.terminal
	kbType := i.flagParams.keyboardType
	var idesToInstall []IDE

	if shouldBeInstalled("jq", "jq", true, true, false) {

	}

	if shouldBeInstalled("Karabiner-Elements", "Karabiner-Elements", false, true, true) {

		appLauncherSurvey := MySurvey{
			message: "App Launcher:",
			options: []string{Spotlight.String(), Launchpad.String(), Alfred.String()},
		}

		kbTypeSurvey := MySurvey{
			message: "Your external keyboard layout:",
			options: []string{PC.String(), Mac.String()},
		}

		terminalSurvey := MySurvey{
			message: "What is your terminal of choice:",
			options: []string{Default.String(), iTerm.String(), Warp.String()},
		}

		if app == "" {
			app = makeSurvey(appLauncherSurvey)
		}
		if term == "" {
			kbType = makeSurvey(kbTypeSurvey)
		}

		if kbType == "" {
			term = makeSurvey(terminalSurvey)
		}

		run("clear")

		if i.flagParams.ides == nil {

			ideSurvey := survey.MultiSelect{
				Message: "IDE keymaps to install:",
				Options: IdeKeymapsSurveyOptions(),
				Help:    "help",
			}

			fullNames := makeMultiSelect(ideSurvey)

			for _, e := range fullNames {
				name, _ := IdeKeymapByFullName(e)
				idesToInstall = append(idesToInstall, name)
			}
		} else {
			idesToInstall = *i.flagParams.ides
		}
	}

	if shouldBeInstalled("Rectangle", "Rectangle", false, false, true) {
		// TODO remember decision and pass to install()
	}

	if shouldBeInstalled("Alt-Tab", "AltTab", false, false, true) {
		// TODO remember decision and pass to install()
	}

	defaults := survey.MultiSelect{
		Message: "Select additional options:",
		Options: []string{
			"Enable Home & End keys (recommended for PC keyboards)",
			"Use F1, F2, etc. keys as standard keys (recommended)",
			"Enable dock auto-hide - 2s delay (recommended)",
			"Change the Dock minimize animation to \"scale\" (recommended)",
			"Disable Spaces rearranging based on most recent use (recommended)",
			"Disable switching to a Space with open windows for the application (recommended)",
			"Enable displays having separated Spaces (recommended)",
			"Put the Dock on the left of the screen",
			"Show hidden files in Finder",
			"Show folders on top in Finder",
			"Show full POSIX path in Finder window title",
			"Shorten windows maximize animation",
			"Disable Mission Control",
		},
		Description: func(value string, index int) string {
			if index < 5 {
				return "Recommended"
			}
			return ""
		},
		Help:     "help",
		PageSize: 15,
	}

	makeMultiSelect(defaults)

	return Params{
		appLauncher:  app,
		terminal:     term,
		keyboardType: kbType,
		ides:         &idesToInstall,
	}
}

func (i Installation) install(params Params) Installation {

	run("killall Karabiner-Elements")

	// do karabiner.json backup
	original := i.karabinerConfigFile()
	backupDest := i.karabinerConfigBackupFile()

	script.Exec("cp " + original + " " + backupDest).Wait()

	// delete existing profile
	deleteProfileJqCmd := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" 'del(.profiles[] | select(.name == \"%s\"))' %s >tmp && mv tmp %s", i.profileName, i.profileName, i.karabinerConfigFile(), i.karabinerConfigFile())
	run(deleteProfileJqCmd)

	// add new karabiner profile
	copyFileFromEmbedFS("configs/karabiner/karabiner-profile.json", "tmp")
	addProfileJqCmd := fmt.Sprintf("jq '.profiles += $profile' %s --slurpfile profile tmp --indent 4 >INPUT.tmp && mv INPUT.tmp %s && rm tmp", i.karabinerConfigFile(), i.karabinerConfigFile())
	run(addProfileJqCmd)

	// rename the profile
	renameJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name == \"PROFILE_NAME\" then .name = \"%s\" else . end)' %s > tmp && mv tmp %s", i.profileName, i.karabinerConfigFile(), i.karabinerConfigFile())
	run(renameJqCmd)

	// unselect other profiles
	unselectJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name != \"%s\" then .selected = false else . end)' %s > tmp && mv tmp %s", i.profileName, i.karabinerConfigFile(), i.karabinerConfigFile())
	run(unselectJqCmd)

	applyRules(i, "main.json")
	applyRules(i, "finder.json")

	switch params.appLauncher {
	case "spotlight":
		fmt.Println("Applying spotlight rules...")
		applyRules(i, "spotlight.json")
	case "launchpad":
		fmt.Println("Applying launchpad rules...")
		applyRules(i, "launchpad.json")
	case "alfred":
		{
			shouldBeInstalled("Alfred", "Alfred 5", false, false, true)
			fmt.Println("Applying alfred rules...")
			applyRules(i, "alfred.json")

			dirs := findMatchingDirs(i.applicationSupportDir()+"/Alfred/Alfred.alfredpreferences/preferences/local", "", "hotkey", "prefs.plist")

			for _, e := range dirs {
				copyFileFromEmbedFS("configs/alfred5/prefs.plist", e)
			}
		}
	default:
		fmt.Println("Value is not A, B, or C")
	}

	switch params.keyboardType {
	case "pc":
		fmt.Println("Applying pc keyboard rules...")
	case "mac":
		fmt.Println("Applying mac keyboard rules...")
		prepareForExternalMacKeyboard(i)
	default:
	}

	switch params.terminal {
	case "default":
		fmt.Println("Applying apple terminal rules...")
		applyRules(i, "apple-terminal.json")
	case "iterm":
		fmt.Println("Applying iterm rules...")
		applyRules(i, "iterm.json")
	case "warp":
		fmt.Println("Applying warp rules...")
		applyRules(i, "warp.json")
	default:
	}

	// reformat using 2 spaces indentation
	run(fmt.Sprintf("jq '.' %s > tmp && mv tmp %s", i.karabinerConfigFile(), i.karabinerConfigFile()))

	run("open -a Karabiner-Elements")

	for _, ide := range *params.ides {
		i.installIdeKeymap(ide)
	}

	run("killall Rectangle")

	if shouldBeInstalled("Rectangle", "Rectangle", false, false, true) {

		rectanglePlist := i.preferencesDir() + "/com.knollsoft.Rectangle.plist"
		copyFileFromEmbedFS("configs/rectangle/Settings.xml", rectanglePlist)

		plutilCmd := fmt.Sprintf("plutil -convert binary1 %s", rectanglePlist)
		run(plutilCmd)

		run("defaults read com.knollsoft.Rectangle")
		run("open -a Rectangle")
	}

	if shouldBeInstalled("Alt-Tab", "AltTab", false, false, true) {
		run("killall AltTab")

		altTabPlist := i.preferencesDir() + "/com.lwouis.alt-tab-macos.plist"
		copyFileFromEmbedFS("configs/alt-tab/Settings.xml", altTabPlist)

		plutilCmd := fmt.Sprintf("plutil -convert binary1 %s", altTabPlist)
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

func (i Installation) installIdeKeymap(ide IDE) {

	if ide.requiresPlugin {
		cmd := fmt.Sprintf("open -na \"%s.app\" --args installPlugins com.intellij.plugins.xwinkeymap", ide.fullName)
		run(cmd)
	}

	destDirs := i.ideDirs(ide)

	for _, d := range destDirs {
		//err := copyFile(i.sourceKeymap(ide), d)
		err := copyFileFromEmbedFS(i.sourceKeymap(ide), d)
		if err != nil {
			fmt.Printf("Error copying to %s: %v\n", d, err)
		} else {
			fmt.Printf("Successfully copied to %s\n", d)
		}
	}
}

func findMatchingDirs(basePath, namePrefix, subDir, fileName string) []string {

	var result []string

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path != basePath && info.IsDir() && fileExists(filepath.Join(basePath, info.Name())) && strings.HasPrefix(info.Name(), namePrefix) {

			fmt.Println(info.Name())
			fmt.Println(path)

			destDir := filepath.Join(path, subDir)
			destFilePath := filepath.Join(destDir, fileName)
			result = append(result, destFilePath)
		}
		return nil
	})

	return result
}

func (i Installation) ideDirs(ide IDE) []string {
	return findMatchingDirs(i.homeDir+ide.parentDir, ide.dir, ide.keymapsDir, ide.destKeymapsFile)
}

func applyRules(i Installation, file string) {
	copyFileFromEmbedFS("configs/karabiner/"+file, i.karabinerComplexModificationsDir()+"/"+file)
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '(.profiles[] | select(.name == \"%s\") | .complex_modifications.rules) += $rules[].rules' %s --slurpfile rules %s/%s >tmp && mv tmp %s", i.profileName, i.profileName, i.karabinerConfigFile(), i.karabinerComplexModificationsDir(), file, i.karabinerConfigFile())
	run(jq)
}

func prepareForExternalMacKeyboard(i Installation) {
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '.profiles |= map(if .name == \"%s\" then walk(if type == \"object\" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' %s --indent 4 >tmp && mv tmp %s", i.profileName, i.profileName, i.karabinerConfigFile(), i.karabinerConfigFile())
	run(jq)
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

	// FIXME do not create base (IntelliJ dir) here
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

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
