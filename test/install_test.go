package install_test

import (
	"flag"
	"fmt"
	"github.com/raxigan/pcfy-my-mac/cmd"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/install"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"github.com/raxigan/pcfy-my-mac/test/test_utils"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInstallWithMacKeyboardLayout(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "mac",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-mac-keyboard-layout.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallWithNoneKeyboardLayout(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "none",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-mac-keyboard-layout.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallWithPcKeyboardLayout(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "pc",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-pc-keyboard-layout.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallWithUnknownKeyboardLayout(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "unknown",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, err := runInstaller(t, params)

	test_utils.AssertErrorContains(t, err, "Unknown keyboard layout: unknown")

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-pc-keyboard-layout.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallWithSpotlightAppLauncher(t *testing.T) {

	params := param.Params{
		AppLauncher:    "spotlight",
		Terminal:       "none",
		KeyboardLayout: "none",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-spotlight-app-launcher.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallWithLaunchpadAppLauncher(t *testing.T) {

	params := param.Params{
		AppLauncher:    "launchpad",
		Terminal:       "none",
		KeyboardLayout: "none",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-launchpad-app-launcher.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallWithAlfredAppLauncher(t *testing.T) {

	params := param.Params{
		AppLauncher:    "alfred",
		Terminal:       "none",
		KeyboardLayout: "none",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-alfred-app-launcher.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallWithUnknownAppLauncher(t *testing.T) {

	params := param.Params{
		AppLauncher:    "unknown",
		Terminal:       "none",
		KeyboardLayout: "none",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	_, _, err := runInstaller(t, params)

	test_utils.AssertErrorContains(t, err, "Unknown app launcher: unknown")
}

func TestInstallWithDefaultTerminal(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "default",
		KeyboardLayout: "none",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-default-terminal.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallWithItermTerminal(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "iterm",
		KeyboardLayout: "none",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-iterm-terminal.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallWithWarpTerminal(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "warp",
		KeyboardLayout: "none",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-warp-terminal.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallAndDoNotCreateNewKarabinerConfigIfItAlreadyExists(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "pc",
		Keymaps:        []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	homeDir := testHomeDir()
	os.MkdirAll(homeDir.KarabinerConfigDir(), 0755)
	common.CopyFile("assets/custom.json", homeDir.KarabinerConfigFile())

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-pc-keyboard-layout-with-custom.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallMany(t *testing.T) {

	params := param.Params{
		AppLauncher:    "alfred",
		Terminal:       "warp",
		KeyboardLayout: "pc",
		Keymaps:        param.IdeKeymapOptions(),
		Blacklist:      []string{"com.spotify.client", "com.apple.finder", "com.apple.AppStore"},
		SystemSettings: []string{
			"Enable Dock auto-hide (2s delay)",
			"change-dock-minimize-animation-to-scale",
			"Enable Home and End keys",
			"Show hidden files in Finder",
			"Show directories on top in Finder",
			"Show full POSIX paths in Finder window title",
		},
	}

	home, output, _ := runInstaller(t, params)
	fmt.Println(output == "")

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-alfred-warp-pc.json"

	test_utils.AssertFilesEqual(t, actual, expected)
	test_utils.AssertFilesEqual(t, filepath.Join(home.PreferencesDir(), "com.knollsoft.Rectangle.plist"), "../assets/rectangle/com.knollsoft.Rectangle.plist")
	test_utils.AssertFilesEqual(t, filepath.Join(home.PreferencesDir(), "com.lwouis.alt-tab-macos.plist"), "expected/com.lwouis.alt-tab-macos.plist")

	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.IntelliJ())), home.IdeKeymapPaths(param.IntelliJ())[0])
	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.IntelliJ())), home.IdeKeymapPaths(param.IntelliJ())[1])
	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.IntelliJCE())), home.IdeKeymapPaths(param.IntelliJCE())[0])
	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.GoLand())), home.IdeKeymapPaths(param.GoLand())[0])
	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.Fleet())), home.IdeKeymapPaths(param.Fleet())[0])

	test_utils.AssertFilesEqual(t, "../assets/system/com.github.pcfy-my-mac.plist", filepath.Join(home.LaunchAgents(), "com.github.pcfy-my-mac.plist"))

	expected = test_utils.ReadFile("expected/output.txt")
	assert.Equal(t, expected, output)
}

func runInstaller(t *testing.T, params param.Params) (install.HomeDir, string, error) {
	common.ExecCommand = fakeExecCommand
	os.Setenv("GO_WANT_HELPER_PROCESS", "1")
	os.Setenv("HOME", testHomeDir().Path)
	defer func() { common.ExecCommand = exec.Command }()
	commander := install.NewDefaultCommander(true)
	homeDir := testHomeDir()
	var err error = nil
	output, err := captureOutput(func() error {
		err := cmd.Launch(homeDir, commander, test_utils.FakeTimeProvider{}, params)
		return err
	})
	t.Cleanup(func() { tearDown(homeDir) })
	return homeDir, output, err
}

func testHomeDir() install.HomeDir {
	wd, _ := os.Getwd()
	return install.HomeDir{
		Path: filepath.Join(wd, "homedir"),
	}
}

//goland:noinspection ALL because it's test code
func tearDown(homeDir install.HomeDir) {
	test_utils.RemoveFiles(filepath.Join(homeDir.Path, ".config"))
	test_utils.RemoveFiles(homeDir.KarabinerConfigBackupFile(test_utils.FakeTimeProvider{}.Now()))
	test_utils.RemoveFiles(homeDir.IdesKeymapPaths(param.IDEKeymaps)...)
	test_utils.RemoveFilesWithExt(homeDir.LibraryDir(), "plist")
	test_utils.RemoveFilesWithExt(homeDir.LibraryDir(), "dict")
	test_utils.RemoveDirs(homeDir.ApplicationSupportDir(), "keymaps")
	test_utils.RemoveDirs(homeDir.ApplicationSupportDir(), "hotkey")
	common.CopyFile(karabinerTestDefaultConfig(homeDir), homeDir.KarabinerConfigFile())
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // reset flags
}

func karabinerTestDefaultConfig(i install.HomeDir) string {
	return filepath.Join(i.KarabinerConfigDir(), "karabiner-default.json")
}

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	execCommand := exec.Command(os.Args[0], cs...)
	return execCommand
}

func captureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	s := string(out)
	fmt.Println(s)
	return s, err
}
