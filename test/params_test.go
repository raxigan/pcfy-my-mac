package install_test

import (
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"github.com/raxigan/pcfy-my-mac/test/test_utils"
	"testing"
)

func TestFailForUnknownParam(t *testing.T) {

	yml := test_utils.Trim(`unknown: hello`)
	_, err := param.CollectYamlParams(yml)

	test_utils.AssertErrorContains(t, err, "Unknown parameter: unknown")
}

func TestFailForInvalidYaml(t *testing.T) {

	yml := test_utils.Trim(`[] :app-launcher:`)
	_, err := param.CollectYamlParams(yml)

	test_utils.AssertErrorContains(t, err, "cannot unmarshal !!seq into param.FileParams")
}

func TestInstallYmlFileDoesNotExist(t *testing.T) {

	_, err := common.TextFromFile("i-do-not-exist.yml")

	test_utils.AssertErrorContains(t, err, "open i-do-not-exist.yml: no such file or directory")
}

func TestInstallInvalidAppLauncher(t *testing.T) {

	yml := test_utils.Trim(`app-launcher: unknown`)
	_, err := param.CollectYamlParams(yml)

	test_utils.AssertErrorContains(t, err, `Invalid param 'app-launcher' value/s 'unknown', valid values:
		spotlight
		launchpad
		alfred
		none`)
}

func TestInstallInvalidTerminal(t *testing.T) {

	yml := test_utils.Trim(`terminal: unknown`)
	_, err := param.CollectYamlParams(yml)

	test_utils.AssertErrorContains(t, err, `Invalid param 'terminal' value/s 'unknown', valid values:
		default
		iterm
		warp
		wave
		none`)
}

func TestInstallInvalidKeyboardLayout(t *testing.T) {

	yml := test_utils.Trim(`keyboard-layout: unknown`)
	_, err := param.CollectYamlParams(yml)

	test_utils.AssertErrorContains(t, err, `Invalid param 'keyboard-layout' value/s 'unknown', valid values:
		pc
		mac`)
}

func TestReadParamsFromYmlFile(t *testing.T) {

	params, _ := param.CollectParams("assets/params.yml")

	test_utils.AssertEquals(t, params.AppLauncher, "alfred")
	test_utils.AssertEquals(t, params.Terminal, "warp")
	test_utils.AssertEquals(t, params.KeyboardLayout, "pc")
	test_utils.AssertSlicesEqual(t, params.Keymaps, []string{"fleet"})
	test_utils.AssertSlicesEqual(t, params.Blacklist, []string{
		"Spotify",
		"FINDER",
		"com.apple.AppStore"},
	)
	test_utils.AssertSlicesEqual(t, params.SystemSettings, []string{
		"enable-dock-auto-hide-2s-delay",
		"change-dock-minimize-animation-to-scale"},
	)
}

func TestReadParamsFromNonexistentFile(t *testing.T) {

	_, err := param.CollectParams("i-do-not-exist.yml")

	test_utils.AssertErrorContains(t, err, "open i-do-not-exist.yml: no such file or directory")
}
