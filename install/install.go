package install

import (
	"errors"
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/bitfield/script"
	"github.com/raxigan/pcfy-my-mac/configs"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func copyFileFromEmbedFS(src, dst string) error {
	configs := &configs.Configs
	data, _ := fs.ReadFile(configs, src)
	os.MkdirAll(filepath.Dir(dst), 0755)
	return os.WriteFile(dst, data, 0755)
}

type Params struct {
	appLauncher       string
	terminal          string
	keyboardLayout    string
	ides              []IDE
	additionalOptions []string
	blacklist         []string
}

type Installation struct {
	homeDir HomeDir
	params  FileParams

	profileName      string
	installationTime time.Time
	cmd              Commander
}

type FileParams struct {
	AppLauncher       *string `yaml:"app-launcher"`
	Terminal          *string
	KeyboardLayout    *string `yaml:"keyboard-layout"`
	Ides              *[]string
	AdditionalOptions *[]string `yaml:"additional-options"`
	Blacklist         *[]string
	Extra             map[string]string `yaml:",inline"`
}

func collectYamlParams(yml *string) (FileParams, error) {

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
				return FileParams{}, e
			}

			data = d
		}
	}

	err := yaml.Unmarshal(data, &fp)
	if err != nil {
		return FileParams{}, err
	}

	if len(fp.Extra) > 0 {
		for field := range fp.Extra {
			return FileParams{}, errors.New("Unknown parameter: " + field)
		}
	}

	validationErr := validateAll(
		func() error {
			if fp.AppLauncher != nil {
				return validateParamValues("app-launcher", &[]string{*fp.AppLauncher}, []string{Spotlight.String(), Launchpad.String(), Alfred.String(), "None"})
			}

			return nil
		},
		func() error {
			if fp.Terminal != nil {
				return validateParamValues("terminal", &[]string{*fp.Terminal}, []string{Default.String(), iTerm.String(), Warp.String(), "None"})
			}

			return nil
		},
		func() error {
			if fp.KeyboardLayout != nil {
				return validateParamValues("keyboard-layout", &[]string{*fp.KeyboardLayout}, []string{PC.String(), Mac.String(), "None"})
			}

			return nil
		},
		func() error {
			return validateParamValues("ides", fp.Ides, append(IdeKeymapOptions(), []string{"all"}...))
		},
		func() error {
			return validateParamValues("additional-options", fp.AdditionalOptions, AdditionalOptions)
		},
	)

	if validationErr != nil {
		return FileParams{}, validationErr
	}

	return FileParams{
		AppLauncher:       fp.AppLauncher,
		Terminal:          fp.Terminal,
		KeyboardLayout:    fp.KeyboardLayout,
		Ides:              fp.Ides,
		AdditionalOptions: fp.AdditionalOptions,
		Blacklist:         fp.Blacklist,
	}, nil
}

func makeSurvey(s MySurvey) string {

	appLauncher := ""

	prompt := &survey.Select{
		Message: s.Message,
		Options: append(s.Options, "None"),
	}

	appLauncher = strings.TrimSpace(appLauncher)

	handleInterrupt(survey.AskOne(prompt, &appLauncher, survey.WithValidator(survey.Required)))

	return strings.ToLower(strings.TrimSpace(appLauncher))
}

func makeMultiSelect(s survey.MultiSelect) []string {
	var appLauncher []string
	handleInterrupt(survey.AskOne(&s, &appLauncher))
	return appLauncher
}

func RunInstaller(homeDir HomeDir, commander Commander, tp TimeProvider, yaml *string) error {
	fp, err := collectYamlParams(yaml)

	if err != nil {
		return err
	}

	installation := Installation{
		homeDir: homeDir,
		params: FileParams{
			AppLauncher:       fp.AppLauncher,
			Terminal:          fp.Terminal,
			KeyboardLayout:    fp.KeyboardLayout,
			Ides:              fp.Ides,
			AdditionalOptions: fp.AdditionalOptions,
			Blacklist:         fp.Blacklist,
		},
		profileName:      "PC mode GOLANG",
		installationTime: tp.Now(),
		cmd:              commander,
	}

	params := installation.collectParams()
	return installation.install(params)
}

func (i Installation) collectParams() Params {

	var app string
	var term string
	var kbType string
	var idesToInstall []IDE
	var blacklist []string

	appLauncherSurvey := MySurvey{
		Message: "App Launcher (will be available with Win(⊞)/Opt(⌥) key):",
		Options: []string{Spotlight.String(), Launchpad.String(), Alfred.String()},
	}

	terminalSurvey := MySurvey{
		Message: "What is your terminal of choice (will be available with Ctrl+Alt+T/Ctrl+Cmd+T shortcut):",
		Options: []string{Default.String(), iTerm.String(), Warp.String()},
	}

	kbTypeSurvey := MySurvey{
		Message: "Your external keyboard layout:",
		Options: []string{PC.String(), Mac.String()},
	}

	if i.params.AppLauncher == nil {
		app = makeSurvey(appLauncherSurvey)
	} else {
		app = *i.params.AppLauncher
	}

	if i.params.Terminal == nil {
		term = makeSurvey(terminalSurvey)
	} else {
		term = *i.params.Terminal
	}

	if i.params.KeyboardLayout == nil {
		kbType = makeSurvey(kbTypeSurvey)
	} else {
		kbType = *i.params.KeyboardLayout
	}

	if i.params.Ides == nil {

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

		if slices.Contains(*i.params.Ides, "all") {
			idesToInstall = IDEKeymaps
		} else {

			var idesFromFlags []IDE

			for _, e := range *i.params.Ides {
				if e != "" {
					byFlag, _ := IdeKeymapByFullName(e)
					idesFromFlags = append(idesFromFlags, byFlag)
				}
			}

			idesToInstall = idesFromFlags
		}
	}

	if i.params.Blacklist == nil {

		msBlacklist := survey.MultiSelect{
			Message: "Select apps to be blacklisted:",
			Options: []string{
				"Spotify",
				"Finder",
				"System Preferences",
			},
			Help: "help",
		}

		blacklist = makeMultiSelect(msBlacklist)
	} else {
		blacklist = *i.params.Blacklist
	}

	var options []string

	if i.params.AdditionalOptions == nil {

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
		options = *i.params.AdditionalOptions
	}

	return Params{
		appLauncher:       app,
		terminal:          term,
		keyboardLayout:    kbType,
		ides:              idesToInstall,
		additionalOptions: options,
		blacklist:         blacklist,
	}
}

func (i Installation) install(params Params) error {

	i.tryInstallDependencies()

	i.run("killall Karabiner-Elements")

	// do karabiner.json backup
	original := i.homeDir.KarabinerConfigFile()
	backupDest := i.homeDir.KarabinerConfigBackupFile(i.installationTime)

	script.Exec("cp " + original + " " + backupDest).Wait()

	// delete existing profile
	deleteProfileJqCmd := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" 'del(.profiles[] | select(.name == \"%s\"))' %s >tmp && mv tmp %s", i.profileName, i.profileName, i.homeDir.KarabinerConfigFile(), i.homeDir.KarabinerConfigFile())
	i.run(deleteProfileJqCmd)

	// add new karabiner profile
	copyFileFromEmbedFS("karabiner/karabiner-profile.json", "tmp")
	addProfileJqCmd := fmt.Sprintf("jq '.profiles += $profile' %s --slurpfile profile tmp --indent 4 >INPUT.tmp && mv INPUT.tmp %s && rm tmp", i.homeDir.KarabinerConfigFile(), i.homeDir.KarabinerConfigFile())
	i.run(addProfileJqCmd)

	// rename the profile
	renameJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name == \"_PROFILE_NAME_\" then .name = \"%s\" else . end)' %s > tmp && mv tmp %s", i.profileName, i.homeDir.KarabinerConfigFile(), i.homeDir.KarabinerConfigFile())
	i.run(renameJqCmd)

	// unselect other profiles
	unselectJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name != \"%s\" then .selected = false else . end)' %s > tmp && mv tmp %s", i.profileName, i.homeDir.KarabinerConfigFile(), i.homeDir.KarabinerConfigFile())
	i.run(unselectJqCmd)

	applyRules(i, "main.json")
	applyRules(i, "finder.json")

	switch strings.ToLower(params.appLauncher) {
	case "spotlight":
		fmt.Println("Applying spotlight rules...")
		applyRules(i, "spotlight.json")
	case "launchpad":
		fmt.Println("Applying launchpad rules...")
		applyRules(i, "launchpad.json")
	case "alfred":
		{
			if i.cmd.Exists("Alfred 4.app") || i.cmd.Exists("Alfred 5.app") {

				fmt.Println("Applying alfred rules...")
				applyRules(i, "alfred.json")

				dirs, err := findMatchingDirs(i.homeDir.ApplicationSupportDir()+"/Alfred/Alfred.alfredpreferences/preferences/local", "", "hotkey", "prefs.plist")

				if err != nil {
					return err
				}

				for _, e := range dirs {
					copyFileFromEmbedFS("alfred5/prefs.plist", e)
				}
			} else {
				printColored(YELLOW, fmt.Sprintf("Alfred app not found. Skipping..."))
			}
		}
	}

	switch strings.ToLower(params.keyboardLayout) {
	case "pc":
		fmt.Println("Applying pc keyboard rules...")
	case "mac":
		fmt.Println("Applying mac keyboard rules...")
		prepareForExternalMacKeyboard(i)
	}

	switch strings.ToLower(params.terminal) {
	case "default":
		fmt.Println("Applying apple terminal rules...")
		applyRules(i, "apple-terminal.json")
	case "iterm":
		if i.cmd.Exists("iTerm.app") {
			fmt.Println("Applying iterm rules...")
			applyRules(i, "iterm.json")
		} else {
			printColored(YELLOW, fmt.Sprintf("iTerm app not found. Skipping..."))
		}
	case "warp":
		{
			if i.cmd.Exists("Warp.app") {
				fmt.Println("Applying warp rules...")
				applyRules(i, "warp.json")
			} else {
				printColored(YELLOW, fmt.Sprintf("Warp app not found. Skipping..."))
			}
		}
	}

	// reformat using 2 spaces indentation
	i.run(fmt.Sprintf("jq '.' %s > tmp && mv tmp %s", i.homeDir.KarabinerConfigFile(), i.homeDir.KarabinerConfigFile()))

	i.run("open -a Karabiner-Elements")

	for _, ide := range params.ides {
		i.installIdeKeymap(ide)
	}

	i.run("killall Rectangle")

	rectanglePlist := i.homeDir.PreferencesDir() + "/com.knollsoft.Rectangle.plist"
	copyFileFromEmbedFS("rectangle/Settings.xml", rectanglePlist)

	plutilCmdRectangle := fmt.Sprintf("plutil -convert binary1 %s", rectanglePlist)
	i.run(plutilCmdRectangle)
	i.run("defaults read com.knollsoft.Rectangle.plist")
	i.run("open -a Rectangle")

	i.run("killall AltTab")

	altTabPlist := i.homeDir.PreferencesDir() + "/com.lwouis.alt-tab-macos.plist"
	copyFileFromEmbedFS("alt-tab/Settings.xml", altTabPlist)

	// set up blacklist

	var mappedStrings []string
	for _, s := range params.blacklist {
		mappedStrings = append(mappedStrings, fmt.Sprintf(`{"ignore":"0","bundleIdentifier":"%s","hide":"1"}`, s))
	}

	result := "[" + strings.Join(mappedStrings, ",") + "]"

	fmt.Println("Blacklist: " + result)

	replaceWordInFile(altTabPlist, "_BLACKLIST_", result)

	plutilCmd := fmt.Sprintf("plutil -convert binary1 %s", altTabPlist)
	i.run(plutilCmd)

	i.run("defaults read com.lwouis.alt-tab-macos.plist")
	i.run("open -a AltTab")

	optionsMap := make(map[string]bool)
	for _, value := range params.additionalOptions {
		optionsMap[strings.ToLower(value)] = true
	}

	fmt.Println("")

	if optionsMap["enable dock auto-hide (2s delay)"] {
		i.run("defaults write com.apple.dock autohide -bool true")
		i.run("defaults write com.apple.dock autohide-delay -float 2 && killall Dock")
	}
	if optionsMap[`change dock minimize animation to "scale"`] {
		i.run(`defaults write com.apple.dock "mineffect" -string "scale" && killall Dock`)
	}
	if optionsMap["enable home & end keys"] {
		fmt.Println("Enable Home & End keys...")
		copyFileFromEmbedFS("system/DefaultKeyBinding.dict", i.homeDir.LibraryDir()+"/KeyBindings/DefaultKeyBinding.dict")
	}
	if optionsMap["show hidden files in finder"] {
		i.run("defaults write com.apple.finder AppleShowAllFiles -bool true")
	}
	if optionsMap["show directories on top in finder"] {
		i.run("defaults write com.apple.finder _FXSortFoldersFirst -bool true")
	}
	if optionsMap["show full posix paths in finder window title"] {

		i.run("defaults write com.apple.finder _FXShowPosixPathInTitle -bool true")
	}

	fmt.Println("SUCCESS")

	return nil
}

func (i Installation) tryInstallDependencies() {

	var notInstalled []string
	var commands []string

	if !i.cmd.Exists("jq") {
		notInstalled = append(notInstalled, "jq")
		commands = append(commands, "brew install jq")
	}

	if !i.cmd.Exists("Karabiner-Elements.app") {
		notInstalled = append(notInstalled, "Karabiner-Elements")
		commands = append(commands, "brew install --cask karabiner-elements")
	}

	if !i.cmd.Exists("AltTab.app") {
		notInstalled = append(notInstalled, "AltTab")
		commands = append(commands, "brew install --cask alt-tab")
	}

	if !i.cmd.Exists("Rectangle.app") {
		notInstalled = append(notInstalled, "Rectangle")
		commands = append(commands, "brew install --cask rectangle")
	}

	if len(notInstalled) > 0 {
		installApp := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("The following dependencies will be installed: %s. Do you agree?", strings.Join(notInstalled, ", ")),
		}
		handleInterrupt(survey.AskOne(prompt, &installApp))

		if !installApp {
			fmt.Printf("Qutting...")
			i.cmd.Exit(0)
		}

		for _, c := range commands {
			i.run(c)
		}
	}
}

func (i Installation) run(cmd string) {
	i.cmd.Run(cmd)
}

func (i Installation) installIdeKeymap(ide IDE) error {

	var destDirs []string

	if ide.multipleDirs {
		a := i.homeDir.IdeKeymapPaths(ide)

		destDirs = a
	} else {
		destDirs = []string{
			i.homeDir.Path + "/" + ide.parentDir + "/" + ide.dir + "/" + ide.keymapsDir + "/" + ide.destKeymapsFile,
		}
	}

	if len(destDirs) == 0 {
		printColored(YELLOW, fmt.Sprintf("%s not found. Skipping...", ide.fullName))
		return nil
	}

	for _, d := range destDirs {
		err := copyFileFromEmbedFS(i.homeDir.SourceKeymap(ide), d)

		if err != nil {
			return err
		}
	}

	return nil
}

func findMatchingDirs(basePath, namePrefix, subDir, fileName string) ([]string, error) {

	var result []string

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {

		if path != basePath && strings.HasPrefix(info.Name(), namePrefix) {
			if err != nil {
				return err
			}

			if fileExists(filepath.Join(basePath, info.Name())) {
				destDir := filepath.Join(path, subDir)
				destFilePath := filepath.Join(destDir, fileName)
				result = append(result, destFilePath)
			}
		}

		return nil
	})

	return result, err
}

func (home HomeDir) IdeKeymapPaths(ide IDE) []string {
	return home.IdesKeymapPaths([]IDE{ide})
}

func (home HomeDir) IdesKeymapPaths(ide []IDE) []string {

	var result []string

	for _, e := range ide {

		dirs, _ := findMatchingDirs(home.Path+e.parentDir, e.dir, e.keymapsDir, e.destKeymapsFile)

		for _, e1 := range dirs {
			result = append(result, e1)
		}
	}

	return result
}

func applyRules(i Installation, file string) {
	copyFileFromEmbedFS("karabiner/"+file, i.homeDir.KarabinerComplexModificationsDir()+"/"+file)
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '(.profiles[] | select(.name == \"%s\") | .complex_modifications.rules) += $rules[].rules' %s --slurpfile rules %s/%s >tmp && mv tmp %s", i.profileName, i.profileName, i.homeDir.KarabinerConfigFile(), i.homeDir.KarabinerComplexModificationsDir(), file, i.homeDir.KarabinerConfigFile())
	i.run(jq)
}

func prepareForExternalMacKeyboard(i Installation) {
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '.profiles |= map(if .name == \"%s\" then walk(if type == \"object\" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' %s --indent 4 >tmp && mv tmp %s", i.profileName, i.profileName, i.homeDir.KarabinerConfigFile(), i.homeDir.KarabinerConfigFile())
	i.run(jq)
}

func validateParamValues(param string, values *[]string, validValues []string) error {

	if values != nil && len(*values) != 0 {

		vals := toLowerSlice(*values)
		valids := toLowerSlice(validValues)

		validMap := make(map[string]bool)
		for _, v := range valids {
			validMap[v] = true
		}

		var invalidValues []string
		for _, val := range vals {
			if !validMap[val] {
				invalidValues = append(invalidValues, val)
			}
		}

		if len(invalidValues) != 0 {
			joined := strings.Join(invalidValues, ", ")
			return errors.New("Invalid param '" + param + "' value/s '" + joined + "', valid values:\n" + strings.Join(validValues, "\n"))
		}
	}

	return nil
}

func toLowerSlice(slice []string) []string {
	for i, s := range slice {
		slice[i] = strings.ToLower(s)
	}
	return slice
}

func validateAll(params ...func() error) error {
	for _, paramFunc := range params {
		if err := paramFunc(); err != nil {
			return err
		}
	}
	return nil
}

func handleInterrupt(err error) {
	if errors.Is(err, terminal.InterruptErr) {
		fmt.Println("Quitting...")
		os.Exit(1)
	}
}
