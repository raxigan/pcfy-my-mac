package install

import (
	"errors"
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
	AppLauncher    string
	Terminal       string
	KeyboardLayout string
	Ides           []string
	SystemSettings []string
	Blacklist      []string
}

type Installation struct {
	Commander
	profileName      string
	installationTime time.Time
}

type FileParams struct {
	AppLauncher    *string `yaml:"app-launcher"`
	Terminal       *string
	KeyboardLayout *string `yaml:"keyboard-layout"`
	Ides           *[]string
	SystemSettings *[]string `yaml:"system-settings"`
	Blacklist      *[]string
	Extra          map[string]string `yaml:",inline"`
}

func CollectYamlParams(yml string) (FileParams, error) {

	fp := FileParams{}

	err := yaml.Unmarshal([]byte(yml), &fp)
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
			return validateParamValues("system-settings", fp.SystemSettings, SystemSettings)
		},
	)

	if validationErr != nil {
		return FileParams{}, validationErr
	}

	return FileParams{
		AppLauncher:    fp.AppLauncher,
		Terminal:       fp.Terminal,
		KeyboardLayout: fp.KeyboardLayout,
		Ides:           fp.Ides,
		SystemSettings: fp.SystemSettings,
		Blacklist:      fp.Blacklist,
	}, nil
}

func RunInstaller(homeDir HomeDir, commander Commander, tp TimeProvider, params Params) error {

	installation := Installation{
		Commander:        commander,
		profileName:      "PC mode GOLANG",
		installationTime: tp.Now(),
	}

	return installation.install(params, homeDir)
}

func CollectParams(fileParams FileParams) Params {

	questionsToAsk := questions

	params := Params{}

	if fileParams.AppLauncher != nil {
		params.AppLauncher = *fileParams.AppLauncher
		questionsToAsk = slices.DeleteFunc(questionsToAsk, func(e *survey.Question) bool { return e.Name == "appLauncher" })
	}

	if fileParams.Terminal != nil {
		params.Terminal = *fileParams.Terminal
		questionsToAsk = slices.DeleteFunc(questionsToAsk, func(e *survey.Question) bool { return e.Name == "terminal" })
	}

	if fileParams.KeyboardLayout != nil {
		params.KeyboardLayout = *fileParams.KeyboardLayout
		questionsToAsk = slices.DeleteFunc(questionsToAsk, func(e *survey.Question) bool { return e.Name == "keyboardLayout" })
	}

	if fileParams.Ides != nil {
		params.Ides = *fileParams.Ides
		questionsToAsk = slices.DeleteFunc(questionsToAsk, func(e *survey.Question) bool { return e.Name == "ides" })
	}

	if fileParams.Blacklist != nil {
		params.Blacklist = *fileParams.Blacklist
		questionsToAsk = slices.DeleteFunc(questionsToAsk, func(e *survey.Question) bool { return e.Name == "blacklist" })
	}

	if fileParams.SystemSettings != nil {
		params.SystemSettings = *fileParams.SystemSettings
		questionsToAsk = slices.DeleteFunc(questionsToAsk, func(e *survey.Question) bool { return e.Name == "systemSettings" })
	}

	handleInterrupt(survey.Ask(questionsToAsk, &params, survey.WithRemoveSelectAll(), survey.WithRemoveSelectNone()))

	return params
}

func (i Installation) install(params Params, home HomeDir) error {

	i.tryInstallDependencies()

	i.Run("killall Karabiner-Elements")

	// do karabiner.json backup
	original := home.KarabinerConfigFile()
	backupDest := home.KarabinerConfigBackupFile(i.installationTime)

	script.Exec("cp " + original + " " + backupDest).Wait()

	// delete existing profile
	deleteProfileJqCmd := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" 'del(.profiles[] | select(.name == \"%s\"))' %s >tmp && mv tmp %s", i.profileName, i.profileName, home.KarabinerConfigFile(), home.KarabinerConfigFile())
	i.Run(deleteProfileJqCmd)

	// add new karabiner profile
	copyFileFromEmbedFS("karabiner/karabiner-profile.json", "tmp")
	addProfileJqCmd := fmt.Sprintf("jq '.profiles += $profile' %s --slurpfile profile tmp --indent 4 >INPUT.tmp && mv INPUT.tmp %s && rm tmp", home.KarabinerConfigFile(), home.KarabinerConfigFile())
	i.Run(addProfileJqCmd)

	// rename the profile
	renameJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name == \"_PROFILE_NAME_\" then .name = \"%s\" else . end)' %s > tmp && mv tmp %s", i.profileName, home.KarabinerConfigFile(), home.KarabinerConfigFile())
	i.Run(renameJqCmd)

	// unselect other profiles
	unselectJqCmd := fmt.Sprintf("jq '.profiles |= map(if .name != \"%s\" then .selected = false else . end)' %s > tmp && mv tmp %s", i.profileName, home.KarabinerConfigFile(), home.KarabinerConfigFile())
	i.Run(unselectJqCmd)

	i.applyRules(home, "main.json")
	i.applyRules(home, "finder.json")

	switch strings.ToLower(params.AppLauncher) {
	case "spotlight":
		i.applyRules(home, "spotlight.json")
	case "launchpad":
		i.applyRules(home, "launchpad.json")
	case "alfred":
		{
			if i.Exists("Alfred 4.app") || i.Exists("Alfred 5.app") {

				i.applyRules(home, "alfred.json")

				dirs, err := findMatchingDirs(home.ApplicationSupportDir()+"/Alfred/Alfred.alfredpreferences/preferences/local", "", "hotkey", "prefs.plist")

				if err != nil {
					return err
				}

				for _, e := range dirs {
					copyFileFromEmbedFS("alfred/prefs.plist", e)
				}
			} else {
				printColored(Yellow, fmt.Sprintf("Alfred app not found. Skipping..."))
			}
		}
	}

	switch strings.ToLower(params.KeyboardLayout) {
	case "mac":
		prepareForExternalMacKeyboard(home, i)
	}

	switch strings.ToLower(params.Terminal) {
	case "default":
		i.applyRules(home, "apple-terminal.json")
	case "iterm":
		if i.Exists("iTerm.app") {
			i.applyRules(home, "iterm.json")
		} else {
			printColored(Yellow, fmt.Sprintf("iTerm app not found. Skipping..."))
		}
	case "warp":
		{
			if i.Exists("Warp.app") {
				i.applyRules(home, "warp.json")
			} else {
				printColored(Yellow, fmt.Sprintf("Warp app not found. Skipping..."))
			}
		}
	}

	// reformat using 2 spaces indentation
	i.Run(fmt.Sprintf("jq '.' %s > tmp && mv tmp %s", home.KarabinerConfigFile(), home.KarabinerConfigFile()))

	i.Run("open -a Karabiner-Elements")

	for _, ide := range params.Ides {
		name, _ := IdeKeymapByFullName(ide)
		i.installIdeKeymap(home, name)
	}

	i.Run("killall Rectangle")

	rectanglePlist := filepath.Join(home.PreferencesDir(), "com.knollsoft.Rectangle.plist")
	copyFileFromEmbedFS("rectangle/Settings.xml", rectanglePlist)

	plutilCmdRectangle := fmt.Sprintf("plutil -convert binary1 %s", rectanglePlist)
	i.Run(plutilCmdRectangle)
	i.Run("defaults read com.knollsoft.Rectangle.plist")
	i.Run("open -a Rectangle")

	i.Run("killall AltTab")

	altTabPlist := filepath.Join(home.PreferencesDir(), "/com.lwouis.alt-tab-macos.plist")
	copyFileFromEmbedFS("alt-tab/Settings.xml", altTabPlist)

	// set up blacklist

	var mappedStrings []string
	for _, s := range params.Blacklist {
		mappedStrings = append(mappedStrings, fmt.Sprintf(`{"ignore":"0","bundleIdentifier":"%s","hide":"1"}`, s))
	}

	result := "[" + strings.Join(mappedStrings, ",") + "]"

	replaceWordInFile(altTabPlist, "_BLACKLIST_", result)

	plutilCmd := fmt.Sprintf("plutil -convert binary1 %s", altTabPlist)
	i.Run(plutilCmd)

	i.Run("defaults read com.lwouis.alt-tab-macos.plist")
	i.Run("open -a AltTab")

	optionsMap := make(map[string]bool)
	for _, value := range params.SystemSettings {
		optionsMap[strings.ToLower(value)] = true
	}

	if optionsMap["enable dock auto-hide (2s delay)"] {
		i.Run("defaults write com.apple.dock autohide -bool true")
		i.Run("defaults write com.apple.dock autohide-delay -float 2 && killall Dock")
	}

	if optionsMap[`change dock minimize animation to "scale"`] {
		i.Run(`defaults write com.apple.dock "mineffect" -string "scale" && killall Dock`)
	}

	if optionsMap["enable home & end keys"] {
		copyFileFromEmbedFS("system/DefaultKeyBinding.dict", home.LibraryDir()+"/KeyBindings/DefaultKeyBinding.dict")
	}

	if optionsMap["show hidden files in finder"] {
		i.Run("defaults write com.apple.finder AppleShowAllFiles -bool true")
	}

	if optionsMap["show directories on top in finder"] {
		i.Run("defaults write com.apple.finder _FXSortFoldersFirst -bool true")
	}
	if optionsMap["show full posix paths in finder window title"] {
		i.Run("defaults write com.apple.finder _FXShowPosixPathInTitle -bool true")
	}

	printColored(Green, "PC'fied!")

	return nil
}

func (i Installation) tryInstallDependencies() {

	var notInstalled []string
	var commands []string

	if !i.Exists("jq") {
		notInstalled = append(notInstalled, "jq")
		commands = append(commands, "brew install jq")
	}

	if !i.Exists("Karabiner-Elements.app") {
		notInstalled = append(notInstalled, "Karabiner-Elements")
		commands = append(commands, "brew install --cask karabiner-elements")
	}

	if !i.Exists("AltTab.app") {
		notInstalled = append(notInstalled, "AltTab")
		commands = append(commands, "brew install --cask alt-tab")
	}

	if !i.Exists("Rectangle.app") {
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
			i.Exit(0)
		}

		for _, c := range commands {
			i.Run(c)
		}
	}
}

func (i Installation) installIdeKeymap(home HomeDir, ide IDE) error {

	var destDirs []string

	if ide.multipleDirs {
		destDirs = home.IdeKeymapPaths(ide)
	} else {
		destDirs = []string{filepath.Join(home.Path, ide.parentDir, ide.dir, ide.keymapsDir, ide.destKeymapsFile)}
	}

	if len(destDirs) == 0 {
		printColored(Yellow, fmt.Sprintf("%s not found. Skipping...", ide.fullName))
		return nil
	}

	for _, d := range destDirs {
		err := copyFileFromEmbedFS(home.SourceKeymap(ide), d)

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

func (i Installation) applyRules(home HomeDir, file string) {
	copyFileFromEmbedFS("karabiner/"+file, home.KarabinerComplexModificationsDir()+"/"+file)
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '(.profiles[] | select(.name == \"%s\") | .complex_modifications.rules) += $rules[].rules' %s --slurpfile rules %s/%s >tmp && mv tmp %s", i.profileName, i.profileName, home.KarabinerConfigFile(), home.KarabinerComplexModificationsDir(), file, home.KarabinerConfigFile())
	i.Run(jq)
}

func prepareForExternalMacKeyboard(home HomeDir, i Installation) {
	jq := fmt.Sprintf("jq --arg PROFILE_NAME \"%s\" '.profiles |= map(if .name == \"%s\" then walk(if type == \"object\" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' %s --indent 4 >tmp && mv tmp %s", i.profileName, i.profileName, home.KarabinerConfigFile(), home.KarabinerConfigFile())
	i.Run(jq)
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
