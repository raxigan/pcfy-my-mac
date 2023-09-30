package install_test

import (
	"flag"
	"github.com/raxigan/pcfy-my-mac/install"
	"github.com/raxigan/pcfy-my-mac/test/test_utils"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallWarpAlfredPC(t *testing.T) {

	params := install.Params{
		AppLauncher:       "alfred",
		Terminal:          "warp",
		KeyboardLayout:    "pc",
		Ides:              []install.IDE{},
		Blacklist:         []string{},
		AdditionalOptions: []string{},
	}

	i, c, _ := runInstaller(t, params)

	actual := i.KarabinerConfigFile()
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

	params := install.Params{
		AppLauncher:       "none",
		Terminal:          "default",
		KeyboardLayout:    "none",
		Ides:              []install.IDE{},
		Blacklist:         []string{},
		AdditionalOptions: []string{},
	}

	i, c, _ := runInstaller(t, params)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-none-default-none.json"

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

func TestInstallItermSpotlightMac(t *testing.T) {

	params := install.Params{
		AppLauncher:       "spotlight",
		Terminal:          "iterm",
		KeyboardLayout:    "mac",
		Ides:              []install.IDE{},
		Blacklist:         []string{},
		AdditionalOptions: []string{},
	}

	i, c, _ := runInstaller(t, params)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-spotlight-iterm-mac.json"

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

func TestInstallNoneLaunchpadPC(t *testing.T) {

	params := install.Params{
		AppLauncher:       "launchpad",
		Terminal:          "warp",
		KeyboardLayout:    "pc",
		Ides:              []install.IDE{},
		Blacklist:         []string{},
		AdditionalOptions: []string{},
	}

	i, c, _ := runInstaller(t, params)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-launchpad-none-pc.json"

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

func TestInstallAllKeymaps(t *testing.T) {

	params := install.Params{
		AppLauncher:       "none",
		Terminal:          "none",
		KeyboardLayout:    "none",
		Ides:              []install.IDE{install.IntelliJ(), install.IntelliJCE(), install.PyCharm(), install.GoLand(), install.Fleet()},
		Blacklist:         []string{},
		AdditionalOptions: []string{},
	}

	i, c, _ := runInstaller(t, params)

	test_utils.AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJ()), i.IdeKeymapPaths(install.IntelliJ())[0])
	test_utils.AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJ()), i.IdeKeymapPaths(install.IntelliJ())[1])
	test_utils.AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJCE()), i.IdeKeymapPaths(install.IntelliJCE())[0])
	test_utils.AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.GoLand()), i.IdeKeymapPaths(install.GoLand())[0])
	test_utils.AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.Fleet()), i.IdeKeymapPaths(install.Fleet())[0])
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

func TestInstallEnableHomeAndEndKeys(t *testing.T) {

	params := install.Params{
		AppLauncher:       "none",
		Terminal:          "none",
		KeyboardLayout:    "none",
		Ides:              []install.IDE{},
		Blacklist:         []string{},
		AdditionalOptions: []string{"Enable Home & End keys"},
	}

	homeDir, _, _ := runInstaller(t, params)

	actual := "../configs/system/DefaultKeyBinding.dict"
	expected := homeDir.LibraryDir() + "/KeyBindings/DefaultKeyBinding.dict"

	test_utils.AssertFilesEqual(t, actual, expected)
}

func TestInstallAdditionalOptions(t *testing.T) {

	params := install.Params{
		AppLauncher:    "none",
		Terminal:       "none",
		KeyboardLayout: "none",
		Ides:           []install.IDE{},
		Blacklist:      []string{},
		AdditionalOptions: []string{
			"Enable Dock auto-hide (2s delay)",
			`Change Dock minimize animation to "scale"`,
			"Enable Home & End keys",
			"Show hidden files in Finder",
			"Show directories on top in Finder",
			"Show full POSIX paths in Finder window title",
		},
	}

	_, c, _ := runInstaller(t, params)

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

func runInstaller(t *testing.T, params install.Params) (install.HomeDir, test_utils.MockCommander, error) {
	commander := test_utils.MockCommander{}
	homeDir := testHomeDir()
	err := install.RunInstaller(homeDir, &commander, test_utils.FakeTimeProvider{}, params)
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
	test_utils.RemoveFiles(i.IdesKeymapPaths(install.IDEKeymaps)...)
	test_utils.RemoveFilesWithExt(i.LibraryDir(), "plist")
	test_utils.RemoveFilesWithExt(i.LibraryDir(), "dict")
	test_utils.CopyFile(karabinerTestDefaultConfig(i), i.KarabinerConfigFile())
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // reset flags
}

func karabinerTestDefaultConfig(i install.HomeDir) string {
	return i.KarabinerConfigDir() + "/karabiner-default.json"
}
