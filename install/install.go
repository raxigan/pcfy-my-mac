package install

import (
	"errors"
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/bitfield/script"
	"github.com/raxigan/pcfy-my-mac/configs"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func copyFileFromEmbedFS(src, dst string) error {
	configs := &configs.Configs
	data, err := fs.ReadFile(configs, src)

	if err != nil {
		fail(err)
	}

	os.MkdirAll(filepath.Dir(dst), 0755)
	return os.WriteFile(dst, data, 0755)
}

type Params struct {
	appLauncher       string
	terminal          string
	keyboardType      string
	ides              *[]IDE
	additionalOptions *[]string
	blacklist         *[]string
}

type Installation struct {
	homeDir string
	params  Params

	profileName      string
	installationTime time.Time
	cmd              Commander
}

type FileParams struct {
	AppLauncher       string `yaml:"app-launcher"`
	Terminal          string
	KeyboardType      string `yaml:"keyboard-type"`
	Ides              *[]string
	AdditionalOptions *[]string `yaml:"additional-options"`
	Blacklist         *[]string
	Extra             map[string]string `yaml:",inline"`
}

func NewInstallation(homeDir string, commander Commander, yml *string) (*Installation, error) {

	var data []byte
	fp := FileParams{}

	if yml != nil {
		data = []byte(*yml)
	} else {

		paramsFile := flag.String("params", "", "YAML file with installer parameters")
		flag.Parse()

		if *paramsFile != "" {
			d, e := os.ReadFile(*paramsFile)

			if e != nil {
				return nil, e
			}

			data = d
		}
	}

	err := yaml.Unmarshal(data, &fp)
	if err != nil {
		return nil, err
	}

	if len(fp.Extra) > 0 {
		for field := range fp.Extra {
			return nil, errors.New("Unknown parameter: " + field)
		}
	}

	validateParamValue("app-launcher", fp.AppLauncher, []string{Spotlight.String(), Launchpad.String(), Alfred.String(), "None"})
	validateParamValue("terminal", fp.Terminal, []string{Default.String(), iTerm.String(), Warp.String(), "None"})
	validateParamValue("keyboard-type", fp.KeyboardType, []string{PC.String(), Mac.String(), "None"})
	validateParamValues("ides", fp.Ides, append(IdeKeymapOptions(), []string{"all"}...))
	// do not validate blacklist
	validateParamValues("additional-options", fp.AdditionalOptions, AdditionalOptions)

	var ides *[]IDE

	if fp.Ides == nil {
		ides = nil
	} else if slices.Contains(*fp.Ides, "all") {
		ides = &IDEKeymaps
	} else {

		var idesFromFlags []IDE

		for _, e := range *fp.Ides {
			if e != "" {
				byFlag, _ := IdeKeymapByFullName(e)
				idesFromFlags = append(idesFromFlags, byFlag)
			}
		}

		ides = &idesFromFlags
	}

	return &Installation{
		homeDir: homeDir,
		params: Params{
			appLauncher:       fp.AppLauncher,
			terminal:          fp.Terminal,
			keyboardType:      fp.KeyboardType,
			ides:              ides,
			additionalOptions: fp.AdditionalOptions,
			blacklist:         fp.Blacklist,
		},
		profileName:      "PC mode GOLANG",
		installationTime: time.Now(),
		cmd:              commander,
	}, nil
}

func (i Installation) KarabinerConfigDir() string {
	return i.homeDir + "/.config/karabiner"
}

func (i Installation) KarabinerConfigFile() string {
	return i.KarabinerConfigDir() + "/karabiner.json"
}

func (i Installation) KarabinerConfigBackupFile() string {
	currentTime := i.installationTime.Format("02-01-2006_15:04:05")
	return i.KarabinerConfigDir() + "/karabiner-" + currentTime + ".json"
}

func (i Installation) KarabinerComplexModificationsDir() string {
	return i.homeDir + "/.config/karabiner/assets/complex_modifications"
}

func (i Installation) applicationSupportDir() string {
	return i.homeDir + "/Library/Application Support"
}

func (i Installation) preferencesDir() string {
	return i.homeDir + "/Library/preferences"
}

func (i Installation) libraryDir() string {
	return i.homeDir + "/Library"
}

func (i Installation) SourceKeymap(ide IDE) string {
	return "keymaps/" + ide.srcKeymapsFile
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

func RunInstaller(homeDir string, commander Commander, yaml *string) (Installation, error) {
	installation, err := NewInstallation(homeDir, commander, yaml)

	if err != nil {
		return Installation{}, err
	}

	params := installation.collectParams()
	return installation.install(params)
}

func (i Installation) collectParams() Params {

	app := i.params.appLauncher
	term := i.params.terminal
	kbType := i.params.keyboardType
	var idesToInstall []IDE
	var blacklist []string

	if i.shouldBeInstalled("jq", "jq", true, false) {

	}

	if i.shouldBeInstalled("Karabiner-Elements", "/Applications/Karabiner-Elements.app", true, true) {

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
			term = makeSurvey(terminalSurvey)
		}

		if kbType == "" {
			kbType = makeSurvey(kbTypeSurvey)
		}

		i.run("clear")

		if i.params.ides == nil {

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
			idesToInstall = *i.params.ides
		}
	}

	if i.shouldBeInstalled("Rectangle", "/Applications/Rectangle.app", false, true) {
		// TODO remember decision and pass to install()
	}

	if i.shouldBeInstalled("Alt-Tab", "/Applications/AltTab.app", false, true) {
		msBlacklist := survey.MultiSelect{
			Message: "Select additional options:",
			Options: []string{
				"Enable Dock auto-hide (2s delay)",
				`Change Dock minimize animation to "scale"`,
				"Enable Home & End keys",
				"Show hidden files in Finder",
				"Show directories on top in Finder",
				"Show full POSIX paths in Finder",
			},
			Description: func(value string, index int) string {
				if index < 2 {
					return "Recommended"
				}
				return ""
			},
			Help:     "help",
			PageSize: 15,
		}

		if i.params.blacklist == nil {
			blacklist = makeMultiSelect(msBlacklist)
		} else {
			blacklist = *i.params.blacklist
		}
	}

	var options []string

	if i.params.additionalOptions == nil {

		ms := survey.MultiSelect{
			Message: "Select additional options:",
			Options: []string{
				"Enable Dock auto-hide (2s delay)",
				`Change Dock minimize animation to "scale"`,
				"Enable Home & End keys",
				"Show hidden files in Finder",
				"Show directories on top in Finder",
				"Show full POSIX paths in Finder",
			},
			Description: func(value string, index int) string {
				if index < 2 {
					return "Recommended"
				}
				return ""
			},
			Help:     "help",
			PageSize: 15,
		}

		options = makeMultiSelect(ms)
	} else {
		options = *i.params.additionalOptions
	}

	return Params{
		appLauncher:       app,
		terminal:          term,
		keyboardType:      kbType,
		ides:              &idesToInstall,
		additionalOptions: &options,
		blacklist:         &blacklist,
	}
}

func (i Installation) install(params Params) (Installation, error) {

	i.run("killall Karabiner-Elements")

	// do karabiner.json backup
	original := i.KarabinerConfigFile()
	backupDest := i.KarabinerConfigBackupFile()

	script.Exec("cp " + original + " " + backupDest).Wait()

	// delete existing profile
	deleteProfileJqCmd := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" 'del(.profiles[] | select(.name == \"%s\"))' %s >tmp && mv tmp %s", i.profileName, i.profileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
	i.run(deleteProfileJqCmd)

	// add new karabiner profile
	copyFileFromEmbedFS("karabiner/karabiner-profile.json", "tmp")
	addProfileJqCmd := fmt.Sprintf("jq '.profiles += $profile' %s --slurpfile profile tmp --indent 4 >INPUT.tmp && mv INPUT.tmp %s && rm tmp", i.KarabinerConfigFile(), i.KarabinerConfigFile())
	i.run(addProfileJqCmd)

	// rename the profile
	renameJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name == \"_PROFILE_NAME_\" then .name = \"%s\" else . end)' %s > tmp && mv tmp %s", i.profileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
	i.run(renameJqCmd)

	// unselect other profiles
	unselectJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name != \"%s\" then .selected = false else . end)' %s > tmp && mv tmp %s", i.profileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
	i.run(unselectJqCmd)

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
			i.shouldBeInstalled("Alfred", "/Applications/Alfred 5.app", false, true)
			fmt.Println("Applying alfred rules...")
			applyRules(i, "alfred.json")

			dirs := findMatchingDirs(i.applicationSupportDir()+"/Alfred/Alfred.alfredpreferences/preferences/local", "", "hotkey", "prefs.plist")

			for _, e := range dirs {
				copyFileFromEmbedFS("alfred5/prefs.plist", e)
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
	i.run(fmt.Sprintf("jq '.' %s > tmp && mv tmp %s", i.KarabinerConfigFile(), i.KarabinerConfigFile()))

	i.run("open -a Karabiner-Elements")

	for _, ide := range *params.ides {
		i.installIdeKeymap(ide)
	}

	i.run("killall Rectangle")

	if i.shouldBeInstalled("Rectangle", "/Applications/Rectangle.app", false, true) {

		rectanglePlist := i.preferencesDir() + "/com.knollsoft.Rectangle.plist"
		copyFileFromEmbedFS("rectangle/Settings.xml", rectanglePlist)

		plutilCmd := fmt.Sprintf("plutil -convert binary1 %s", rectanglePlist)
		i.run(plutilCmd)

		i.run("defaults read com.knollsoft.Rectangle.plist")
		i.run("open -a Rectangle")
	}

	if i.shouldBeInstalled("Alt-Tab", "/Applications/AltTab.app", false, true) {
		i.run("killall AltTab")

		altTabPlist := i.preferencesDir() + "/com.lwouis.alt-tab-macos.plist"
		copyFileFromEmbedFS("alt-tab/Settings.xml", altTabPlist)

		// set up blacklist
		bl := *params.blacklist

		var mappedStrings []string
		for _, s := range bl {
			mappedStrings = append(mappedStrings, fmt.Sprintf(`{"ignore":"0","bundleIdentifier":"%s","hide":"1"}`, s))
		}

		result := "[" + strings.Join(mappedStrings, ",") + "]"

		fmt.Println("Blacklist: " + result)

		replaceWordInFile(altTabPlist, "_BLACKLIST_", result)

		plutilCmd := fmt.Sprintf("plutil -convert binary1 %s", altTabPlist)
		i.run(plutilCmd)

		i.run("defaults read com.lwouis.alt-tab-macos.plist")
		i.run("open -a AltTab")
	}

	optionsMap := make(map[string]bool)
	for _, value := range *params.additionalOptions {
		optionsMap[value] = true
	}

	switch {
	case optionsMap["Enable Dock auto-hide (2s delay)"]:
		i.run("defaults write com.apple.dock autohide -bool true")
		i.run("defaults write com.apple.dock autohide-delay -float 2 && killall Dock")
	case optionsMap[`Change Dock minimize animation to "scale"`]:
		i.run(`defaults write com.apple.dock "mineffect" -string "scale" && killall Dock`)
	case optionsMap["Enable Home & End keys"]:
		fmt.Println("Enable Home & End keys...")
		copyFileFromEmbedFS("system/DefaultKeyBinding.dict", i.libraryDir()+"/KeyBindings/DefaultKeyBinding.dict")
	case optionsMap["Show hidden files in Finder"]:
		i.run("defaults write com.apple.finder AppleShowAllFiles -bool true")
	case optionsMap["Show directories on top in Finder"]:
		i.run("defaults write com.apple.finder _FXSortFoldersFirst -bool true")
	case optionsMap["Show full POSIX paths in Finder"]:
		i.run("defaults write com.apple.finder _FXShowPosixPathInTitle -bool true")
	}

	fmt.Println("SUCCESS")

	return i, nil
}

func (i Installation) shouldBeInstalled(appName string, appFile string, isRequired bool, isCask bool) bool {

	exists := i.cmd.exists(appFile)

	if exists {
		return true
	}

	i.run("clear")
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

func (i Installation) run(cmd string) {
	i.cmd.run(cmd)
}

func (i Installation) installIdeKeymap(ide IDE) {

	var destDirs []string

	if ide.multipleDirs {
		destDirs = i.IdeKeymapPaths(ide)
	} else {
		destDirs = []string{
			i.homeDir + "/" + ide.parentDir + "/" + ide.dir + "/" + ide.keymapsDir + "/" + ide.destKeymapsFile,
		}
	}

	if len(destDirs) == 0 {
		printColored(YELLOW, fmt.Sprintf("%s not found. Skipping...", ide.fullName))
		return
	}

	for _, d := range destDirs {
		err := copyFileFromEmbedFS(i.SourceKeymap(ide), d)
		checkError(err)
	}
}

func findMatchingDirs(basePath, namePrefix, subDir, fileName string) []string {

	var result []string

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {

		if path != basePath && strings.HasPrefix(info.Name(), namePrefix) {
			checkError(err)

			if fileExists(filepath.Join(basePath, info.Name())) {
				destDir := filepath.Join(path, subDir)
				destFilePath := filepath.Join(destDir, fileName)
				result = append(result, destFilePath)
			}
		}

		return nil
	})

	return result
}

func (i Installation) IdeKeymapPaths(ide IDE) []string {
	return i.IdesKeymapPaths([]IDE{ide})
}

func (i Installation) IdesKeymapPaths(ide []IDE) []string {

	var result []string

	for _, e := range ide {

		dirs := findMatchingDirs(i.homeDir+e.parentDir, e.dir, e.keymapsDir, e.destKeymapsFile)

		for _, e1 := range dirs {
			result = append(result, e1)
		}
	}

	return result
}

func applyRules(i Installation, file string) {
	copyFileFromEmbedFS("karabiner/"+file, i.KarabinerComplexModificationsDir()+"/"+file)
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '(.profiles[] | select(.name == \"%s\") | .complex_modifications.rules) += $rules[].rules' %s --slurpfile rules %s/%s >tmp && mv tmp %s", i.profileName, i.profileName, i.KarabinerConfigFile(), i.KarabinerComplexModificationsDir(), file, i.KarabinerConfigFile())
	i.run(jq)
}

func prepareForExternalMacKeyboard(i Installation) {
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '.profiles |= map(if .name == \"%s\" then walk(if type == \"object\" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' %s --indent 4 >tmp && mv tmp %s", i.profileName, i.profileName, i.KarabinerConfigFile(), i.KarabinerConfigFile())
	i.run(jq)
}

func validateParamValue(param, value string, validValues []string) {

	if value != "" {
		v := toLowerSlice(validValues)

		if !slices.Contains(v, strings.ToLower(value)) {
			fmt.Println("Invalid param '" + param + "' value '" + value + "', valid values:\n" + strings.Join(v, "\n"))
			os.Exit(1)
		}
	}
}

func validateParamValues(param string, values *[]string, validValues []string) {

	if values != nil && len(*values) != 0 {

		validMap := make(map[string]bool)
		for _, v := range validValues {
			validMap[v] = true
		}

		var invalidValues []string
		for _, val := range *values {
			if !validMap[val] {
				invalidValues = append(invalidValues, val)
			}
		}

		if len(invalidValues) != 0 {
			joined := strings.Join(invalidValues, ", ")
			fmt.Println("Invalid param '" + param + "' values '" + joined + "', valid values:\n" + strings.Join(validValues, "\n"))
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

func checkError(err error) {
	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	log.Fatalf("Error: %s", err)
}
