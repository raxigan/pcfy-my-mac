package install

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/raxigan/pcfy-my-mac/configs"
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

type Installation struct {
	Commander
	Params
	HomeDir
	profileName      string
	installationTime time.Time
}

func RunInstaller(homeDir HomeDir, commander Commander, tp TimeProvider, params Params) error {

	installation := Installation{
		Commander:        commander,
		HomeDir:          homeDir,
		Params:           params,
		profileName:      "PCfy",
		installationTime: tp.Now(),
	}

	return installation.install()
}

func CollectParams(fileParams FileParams) Params {

	questionsToAsk := questions

	fp := Params{}

	m := map[string]bool{
		"appLauncher":    fileParams.AppLauncher != nil,
		"terminal":       fileParams.Terminal != nil,
		"keyboardLayout": fileParams.KeyboardLayout != nil,
		"ides":           fileParams.Ides != nil,
		"blacklist":      fileParams.Blacklist != nil,
		"systemSettings": fileParams.SystemSettings != nil,
	}

	for k, v := range m {
		if v {
			questionsToAsk = slices.DeleteFunc(questionsToAsk, func(e *survey.Question) bool { return e.Name == k })
		}
	}

	handleInterrupt(survey.Ask(questionsToAsk, &fp, survey.WithRemoveSelectAll(), survey.WithRemoveSelectNone()))

	return Params{
		AppLauncher:    getOrDefaultString(fp.AppLauncher, fileParams.AppLauncher),
		Terminal:       getOrDefaultString(fp.Terminal, fileParams.Terminal),
		KeyboardLayout: getOrDefaultString(fp.KeyboardLayout, fileParams.KeyboardLayout),
		Ides:           getOrDefaultSlice(fp.Ides, fileParams.Ides),
		Blacklist:      getOrDefaultSlice(fp.Blacklist, fileParams.Blacklist),
		SystemSettings: getOrDefaultSlice(fp.SystemSettings, fileParams.SystemSettings),
	}
}

func (i Installation) install() error {

	tasks := []Task{
		DownloadDependencies(),
		CloseKarabiner(),
		BackupKarabinerConfig(),
		DeleteExistingKarabinerProfile(),
		CreateKarabinerProfile(),
		NameKarabinerProfile(),
		UnselectOtherKarabinerProfiles(),
		ApplyMainKarabinerRules(),
		ApplyAppLauncherRules(),
		ApplyKeyboardLayoutRules(),
		ApplyTerminalRules(),
		ReformatKarabinerConfigFile(),
		OpenKarabiner(),
		CopyIdeKeymaps(),
		CloseRectangle(),
		CopyRectanglePreferences(),
		OpenRectangle(),
		CloseAltTab(),
		InstallAltTabPreferences(),
		OpenAltTab(),
		ApplySystemSettings(),
	}

	for _, task := range tasks {
		i.Commander.TryPrint(Colored(Blue, "TASK"), task.name)

		err := task.exec(i)

		if err != nil {
			return err
		}
	}

	printlnColored(Green, "PC'fied!")

	return nil
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
