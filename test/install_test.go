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
	"path/filepath"
	"testing"
)

func TestInstallWarpAlfredPC(t *testing.T) {

	params := param.Params{
		AppLauncher:    "alfred",
		Terminal:       "warp",
		KeyboardLayout: "pc",
		Ides:           []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, c, _ := runInstaller(t, params)

	actual := home.KarabinerConfigFile()
	expected := "expected/karabiner-expected-alfred-warp-pc.json"

	test_utils.AssertFilesEqual(t, actual, expected)
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
	})
}

func TestInstallNoneDefaultNone(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "default",
		KeyboardLayout: "none",
		Ides:           []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, c, _ := runInstaller(t, params)

	test_utils.AssertFilesEqual(t, home.KarabinerConfigFile(), "expected/karabiner-expected-none-default-none.json")
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
	})
}

func TestInstallItermSpotlightMac(t *testing.T) {

	params := param.Params{
		AppLauncher:    "spotlight",
		Terminal:       "iterm",
		KeyboardLayout: "mac",
		Ides:           []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, c, _ := runInstaller(t, params)

	test_utils.AssertFilesEqual(t, home.KarabinerConfigFile(), "expected/karabiner-expected-spotlight-iterm-mac.json")
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
	})
}

func TestInstallNoneLaunchpadPC(t *testing.T) {

	params := param.Params{
		AppLauncher:    "launchpad",
		Terminal:       "warp",
		KeyboardLayout: "pc",
		Ides:           []string{},
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, c, _ := runInstaller(t, params)

	test_utils.AssertFilesEqual(t, home.KarabinerConfigFile(), "expected/karabiner-expected-launchpad-none-pc.json")
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
	})
}

func TestInstallAllKeymaps(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "none",
		Ides:           param.IdeKeymapOptions(),
		Blacklist:      []string{},
		SystemSettings: []string{},
	}

	home, c, _ := runInstaller(t, params)

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
	})
}

func TestInstallSystemSettings(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "none",
		Ides:           []string{},
		Blacklist:      []string{},
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

	test_utils.AssertFilesEqual(t, "../assets/system/DefaultKeyBinding.dict", filepath.Join(home.LibraryDir(), "KeyBindings/DefaultKeyBinding.dict"))
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

func TestInstallBlacklist(t *testing.T) {

	params := param.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "none",
		Ides:           []string{},
		Blacklist:      []string{"Spotify", "FINDER", "com.apple.AppStore"},
		SystemSettings: []string{},
	}

	home, c, _ := runInstaller(t, params)

	fmt.Println(c)

	test_utils.AssertFilesEqual(t, filepath.Join(home.PreferencesDir(), "com.knollsoft.Rectangle.plist"), "expected/com.knollsoft.Rectangle.plist")
	test_utils.AssertFilesEqual(t, filepath.Join(home.PreferencesDir(), "com.lwouis.alt-tab-macos.plist"), "expected/com.lwouis.alt-tab-macos.plist")
}

func runInstaller(t *testing.T, params param.Params) (install.HomeDir, test_utils.MockCommander, error) {
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
