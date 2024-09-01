package param

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"gopkg.in/yaml.v3"
	"slices"
	"strings"
)

type Params struct {
	AppLauncher    string
	Terminal       string
	KeyboardLayout string
	Keymaps        []string
	SystemSettings []string
	Blacklist      []string
}

type FileParams struct {
	AppLauncher    *string `yaml:"app-launcher"`
	Terminal       *string
	KeyboardLayout *string `yaml:"keyboard-layout"`
	Keymaps        *[]string
	SystemSettings *[]string `yaml:"system-settings"`
	Blacklist      *[]string
	Extra          map[string]string `yaml:",inline"`
}

func CollectParams(paramsFile string) (Params, error) {
	fileParams := FileParams{}

	if paramsFile != "" {
		yamlStr, err := common.TextFromFile(paramsFile)

		if err != nil {
			return Params{}, err
		}

		fileParams, err = CollectYamlParams(yamlStr)

		if err != nil {
			return Params{}, err
		}
	}

	return CollectSurveyParams(fileParams), nil
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

	validationErr := ValidateAll(
		func() error {
			if fp.AppLauncher != nil {
				return ValidateParamValues("app-launcher", &[]string{*fp.AppLauncher}, []string{Spotlight, Launchpad, Alfred, None})
			}

			return nil
		},
		func() error {
			if fp.Terminal != nil {
				return ValidateParamValues("terminal", &[]string{*fp.Terminal}, []string{Default, ITerm, Warp, Wave, None})
			}

			return nil
		},
		func() error {
			if fp.KeyboardLayout != nil {
				return ValidateParamValues("keyboard-layout", &[]string{*fp.KeyboardLayout}, []string{PC, Mac, None})
			}

			return nil
		},
		func() error {
			return ValidateParamValues("ides", fp.Keymaps, append(IdeKeymapOptions(), []string{"all"}...))
		},
		func() error {
			return ValidateParamValues("system-settings", fp.SystemSettings, SystemSettings)
		},
	)

	if validationErr != nil {
		return FileParams{}, validationErr
	}

	return FileParams{
		AppLauncher:    fp.AppLauncher,
		Terminal:       fp.Terminal,
		KeyboardLayout: fp.KeyboardLayout,
		Keymaps:        fp.Keymaps,
		SystemSettings: fp.SystemSettings,
		Blacklist:      fp.Blacklist,
	}, nil
}

func CollectSurveyParams(fileParams FileParams) Params {

	questionsToAsk := questions

	fp := Params{}

	qNameToIfShouldNotBeAsked := map[string]bool{
		"appLauncher":    fileParams.AppLauncher != nil,
		"terminal":       fileParams.Terminal != nil,
		"keyboardLayout": fileParams.KeyboardLayout != nil,
		"keymaps":        fileParams.Keymaps != nil,
		"systemSettings": fileParams.SystemSettings != nil,
	}

	for k, v := range qNameToIfShouldNotBeAsked {
		if v {
			questionsToAsk = slices.DeleteFunc(questionsToAsk, func(e *survey.Question) bool { return e.Name == k })
		}
	}

	if !qNameToIfShouldNotBeAsked["keymaps"] {

		ides := findIdes()

		if len(ides) > 0 {
			questionsToAsk = append(questionsToAsk,
				&survey.Question{
					Name: "keymaps",
					Prompt: &survey.MultiSelect{
						Message: "Select tools to install keymap for:",
						Options: ides,
						Help:    "IDEs/tools to apply the PC keymaps to",
					},
				},
			)
		}
	}

	common.HandleInterrupt(survey.Ask(questionsToAsk, &fp, survey.WithRemoveSelectAll(), survey.WithRemoveSelectNone(), survey.WithKeepFilter(false)))

	return Params{
		AppLauncher:    common.GetOrDefaultString(fp.AppLauncher, fileParams.AppLauncher),
		Terminal:       common.GetOrDefaultString(fp.Terminal, fileParams.Terminal),
		KeyboardLayout: common.GetOrDefaultString(fp.KeyboardLayout, fileParams.KeyboardLayout),
		Keymaps:        common.GetOrDefaultSlice(fp.Keymaps, fileParams.Keymaps),
		Blacklist:      common.GetOrDefaultSlice(fp.Blacklist, fileParams.Blacklist),
		SystemSettings: common.GetOrDefaultSlice(fp.SystemSettings, fileParams.SystemSettings),
	}
}

func ToSimpleParamName(name string) string {
	loweredAndSnaked := strings.TrimSpace(strings.ReplaceAll(strings.ToLower(name), " ", "-"))
	noBrackets := strings.ReplaceAll(strings.ReplaceAll(loweredAndSnaked, "(", ""), ")", "")
	return strings.ReplaceAll(noBrackets, "\"", "")
}

func findIdes() []string {
	var ides []string

	for _, e := range IDEKeymaps {
		if common.Exists(e.FullName + ".app") {
			ides = append(ides, strings.TrimSuffix(e.FullName, ".app"))
		}
	}

	return ides
}
