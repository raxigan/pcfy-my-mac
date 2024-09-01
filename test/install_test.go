package install_test

import (
	"flag"
	"fmt"
	"github.com/raxigan/pcfy-my-mac/cmd"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/install"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"github.com/raxigan/pcfy-my-mac/test/test_utils"
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

	home, c, _ := runInstaller(t, params)

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

	test_utils.AssertSlicesEqual(t, c.CommandsLog, []string{
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"plutil -convert binary1 /homedir/Library/Preferences/com.knollsoft.Rectangle.plist",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"plutil -convert binary1 /homedir/Library/Preferences/com.lwouis.alt-tab-macos.plist",
		"defaults read com.lwouis.alt-tab-macos.plist",
		"open -a AltTab",
		"defaults write com.apple.dock autohide -bool true",
		"defaults write com.apple.dock autohide-delay -float 2 && killall Dock",
		`defaults write com.apple.dock "mineffect" -string "scale" && killall Dock`,
		"defaults write com.apple.finder AppleShowAllFiles -bool true",
		"defaults write com.apple.finder _FXSortFoldersFirst -bool true",
		"defaults write com.apple.finder _FXShowPosixPathInTitle -bool true",
		"hidutil property --set '{\"UserKeyMapping\":[ { \"HIDKeyboardModifierMappingSrc\": 0x7000000E0, \"HIDKeyboardModifierMappingDst\": 0x7000000E3 }, { \"HIDKeyboardModifierMappingSrc\": 0x7000000E3, \"HIDKeyboardModifierMappingDst\": 0x7000000E0 }, { \"HIDKeyboardModifierMappingSrc\": 0x7000000E4, \"HIDKeyboardModifierMappingDst\": 0x7000000E7 }, { \"HIDKeyboardModifierMappingSrc\": 0x7000000E7, \"HIDKeyboardModifierMappingDst\": 0x7000000E4 } ]}'",
		"clear",
	})
}

func runInstaller(t *testing.T, params param.Params) (install.HomeDir, test_utils.MockCommander, error) {
	common.ExecCommand = fakeExecCommand
	defer func() { common.ExecCommand = exec.Command }()
	commander := *test_utils.NewMockCommander()
	homeDir := testHomeDir()
	err := cmd.Launch(homeDir, &commander, test_utils.FakeTimeProvider{}, params)
	t.Cleanup(func() { tearDown(homeDir) })
	return homeDir, commander, err
}

func testHomeDir() install.HomeDir {
	wd, _ := os.Getwd()
	return install.HomeDir{
		Path: filepath.Join(wd, "homedir"),
	}
}

func tearDown(homeDir install.HomeDir) {
	test_utils.RemoveFiles(filepath.Join(homeDir.Path, ".config"))
	test_utils.RemoveFiles(homeDir.KarabinerConfigBackupFile(test_utils.FakeTimeProvider{}.Now()))
	test_utils.RemoveFiles(homeDir.IdesKeymapPaths(param.IDEKeymaps)...)
	test_utils.RemoveFilesWithExt(homeDir.LibraryDir(), "plist")
	test_utils.RemoveFilesWithExt(homeDir.LibraryDir(), "dict")
	common.CopyFile(karabinerTestDefaultConfig(homeDir), homeDir.KarabinerConfigFile())
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // reset flags
}

func karabinerTestDefaultConfig(i install.HomeDir) string {
	return filepath.Join(i.KarabinerConfigDir(), "karabiner-default.json")
}

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, "siema")
	os.Exit(0)
}
