package install_test

import (
	"flag"
	"github.com/raxigan/pcfy-my-mac/cmd"
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/install"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"github.com/raxigan/pcfy-my-mac/test/test_utils"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallWithMacKeyboardLayout(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "mac",
		Ides:           []string{},
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
		Ides:           []string{},
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
		Ides:           []string{},
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
		Ides:           []string{},
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
		Ides:           []string{},
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
		Ides:           []string{},
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
		Ides:           []string{},
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
		Ides:           []string{},
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
		Ides:           []string{},
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
		Ides:           []string{},
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
		Ides:           []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, _, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-warp-terminal.json"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallMany(t *testing.T) {

	params := param.Params{
		AppLauncher:    "alfred",
		Terminal:       "warp",
		KeyboardLayout: "pc",
		Ides:           param.IdeKeymapOptions(),
		Blacklist:      []string{"Spotify", "FINDER", "com.apple.AppStore"},
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
	expected := "expected/karabiner-expected-alfred-warp-pc.json"

	test_utils.AssertFilesEqual(t, actual, expected)

	test_utils.AssertFilesEqual(t, filepath.Join(home.PreferencesDir(), "com.knollsoft.Rectangle.plist"), "expected/com.knollsoft.Rectangle.plist")
	test_utils.AssertFilesEqual(t, filepath.Join(home.PreferencesDir(), "com.lwouis.alt-tab-macos.plist"), "expected/com.lwouis.alt-tab-macos.plist")

	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.IntelliJ())), home.IdeKeymapPaths(param.IntelliJ())[0])
	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.IntelliJ())), home.IdeKeymapPaths(param.IntelliJ())[1])
	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.IntelliJCE())), home.IdeKeymapPaths(param.IntelliJCE())[0])
	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.GoLand())), home.IdeKeymapPaths(param.GoLand())[0])
	test_utils.AssertFilesEqual(t, filepath.Join("../assets", home.SourceKeymap(param.Fleet())), home.IdeKeymapPaths(param.Fleet())[0])

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
	})
}

func runInstaller(t *testing.T, params param.Params) (install.HomeDir, test_utils.MockCommander, error) {
	commander := *test_utils.NewMockCommander()
	homeDir := testHomeDir()
	test_utils.RemoveFiles(filepath.Join(homeDir.Path, ".config"))
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

func tearDown(i install.HomeDir) {
	test_utils.RemoveFiles(i.KarabinerConfigBackupFile(test_utils.FakeTimeProvider{}.Now()))
	test_utils.RemoveFiles(i.IdesKeymapPaths(param.IDEKeymaps)...)
	test_utils.RemoveFilesWithExt(i.LibraryDir(), "plist")
	test_utils.RemoveFilesWithExt(i.LibraryDir(), "dict")
	common.CopyFile(karabinerTestDefaultConfig(i), i.KarabinerConfigFile())
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // reset flags
}

func karabinerTestDefaultConfig(i install.HomeDir) string {
	return filepath.Join(i.KarabinerConfigDir(), "karabiner-default.json")
}
