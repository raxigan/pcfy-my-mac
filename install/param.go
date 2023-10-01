package install

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type Params struct {
	AppLauncher    string
	Terminal       string
	KeyboardLayout string
	Ides           []string
	SystemSettings []string
	Blacklist      []string
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
