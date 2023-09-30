package install_test

import (
	"flag"
	"github.com/raxigan/pcfy-my-mac/install"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallWithFirstOptionsFromPrompts(t *testing.T) {

	i, _, _ := runInstaller(t, nil)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-spotlight-default-pc.json"

	AssertFilesEqual(t, actual, expected)
}

func TestInstallFromYamlFile(t *testing.T) {

	os.Args = []string{"script_name", "--params=params.yml"}
	i, _, _ := runInstaller(t, nil)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-alfred-warp-pc.json"

	AssertFilesEqual(t, actual, expected)
}

func TestInstallWarpAlfredPC(t *testing.T) {

	yml := trim(`
		app-launcher: alfred
		terminal: warp
		keyboard-layout: pc
		ides: [ ]
		additional-options: [ ]
		blacklist: [ com.apple.Preview ]`,
	)

	i, c, _ := runInstaller(t, &yml)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-alfred-warp-pc.json"

	AssertFilesEqual(t, actual, expected)
	AssertSlicesEqual(t, c.CommandsLog, []string{
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"plutil -convert binary1 /homedir/Library/preferences/com.knollsoft.Rectangle.plist",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"plutil -convert binary1 /homedir/Library/preferences/com.lwouis.alt-tab-macos.plist",
		"defaults read com.lwouis.alt-tab-macos.plist",
		"open -a AltTab",
	})
}

func TestInstallNoneDefaultNone(t *testing.T) {

	yml := trim(`
		app-launcher: None
		terminal: Default
		keyboard-layout: None
		ides: [ ]
		additional-options: [ ]
		blacklist: [ com.apple.Preview ]`,
	)

	i, c, _ := runInstaller(t, &yml)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-none-default-none.json"

	AssertFilesEqual(t, actual, expected)
	AssertSlicesEqual(t, c.CommandsLog, []string{
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"plutil -convert binary1 /homedir/Library/preferences/com.knollsoft.Rectangle.plist",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"plutil -convert binary1 /homedir/Library/preferences/com.lwouis.alt-tab-macos.plist",
		"defaults read com.lwouis.alt-tab-macos.plist",
		"open -a AltTab",
	})
}

func TestInstallItermSpotlightMac(t *testing.T) {

	yml := trim(`
		app-launcher: spotlight
		terminal: iterm
		keyboard-layout: mac
		ides: [ ]
		additional-options: [ ]
		blacklist: [ ]`,
	)

	i, c, _ := runInstaller(t, &yml)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-spotlight-iterm-mac.json"

	AssertFilesEqual(t, actual, expected)
	AssertSlicesEqual(t, c.CommandsLog, []string{
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"plutil -convert binary1 /homedir/Library/preferences/com.knollsoft.Rectangle.plist",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"plutil -convert binary1 /homedir/Library/preferences/com.lwouis.alt-tab-macos.plist",
		"defaults read com.lwouis.alt-tab-macos.plist",
		"open -a AltTab",
	})
}

func TestInstallNoneLaunchpadPC(t *testing.T) {

	yml := trim(`
		app-launcher: launchpad
		terminal: warp
		keyboard-layout: pc
		ides: [ ]
		additional-options: [ ]
		blacklist: [ ]`,
	)

	i, c, _ := runInstaller(t, &yml)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-launchpad-none-pc.json"

	AssertFilesEqual(t, actual, expected)
	AssertSlicesEqual(t, c.CommandsLog, []string{
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"plutil -convert binary1 /homedir/Library/preferences/com.knollsoft.Rectangle.plist",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"plutil -convert binary1 /homedir/Library/preferences/com.lwouis.alt-tab-macos.plist",
		"defaults read com.lwouis.alt-tab-macos.plist",
		"open -a AltTab",
	})
}

func TestInstallAllKeymaps(t *testing.T) {

	yml := trim(`
		app-launcher: None
		terminal: None
		keyboard-layout: None
		ides: [ "all" ]
		additional-options: [ ]
		blacklist: [ ]`,
	)

	i, c, _ := runInstaller(t, &yml)

	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJ()), i.IdeKeymapPaths(install.IntelliJ())[0])
	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJ()), i.IdeKeymapPaths(install.IntelliJ())[1])
	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJCE()), i.IdeKeymapPaths(install.IntelliJCE())[0])
	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.GoLand()), i.IdeKeymapPaths(install.GoLand())[0])
	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.Fleet()), i.IdeKeymapPaths(install.Fleet())[0])
	AssertSlicesEqual(t, c.CommandsLog, []string{
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"plutil -convert binary1 /homedir/Library/preferences/com.knollsoft.Rectangle.plist",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"plutil -convert binary1 /homedir/Library/preferences/com.lwouis.alt-tab-macos.plist",
		"defaults read com.lwouis.alt-tab-macos.plist",
		"open -a AltTab",
	})
}

func TestFailForUnknownParam(t *testing.T) {

	yml := trim(`unknown: hello`)

	_, c, err := runInstaller(t, &yml)

	AssertErrorContains(t, err, "Unknown parameter: unknown")
	AssertSlicesEqual(t, c.CommandsLog, []string{})
}

func TestFailForInvalidYaml(t *testing.T) {

	yml := trim(`[] :app-launcher:`)

	_, c, err := runInstaller(t, &yml)

	AssertErrorContains(t, err, "cannot unmarshal !!seq into install.FileParams")
	AssertSlicesEqual(t, c.CommandsLog, []string{})
}

func TestInstallYmlFileDoesNotExist(t *testing.T) {

	os.Args = []string{"script_name", "--params=nope.yml"}
	_, c, err := runInstaller(t, nil)

	AssertErrorContains(t, err, "open nope.yml: no such file or directory")
	AssertSlicesEqual(t, c.CommandsLog, []string{})
}

func TestInstallEnableHomeAndEndKeys(t *testing.T) {

	yml := trim(`
		app-launcher: None
		terminal: None
		keyboard-layout: None
		ides: [ ]
		additional-options: [ "Enable Home & End keys" ]
		blacklist: [ ]`,
	)

	homeDir, _, _ := runInstaller(t, &yml)

	actual := "../configs/system/DefaultKeyBinding.dict"
	expected := homeDir.LibraryDir() + "/KeyBindings/DefaultKeyBinding.dict"

	AssertFilesEqual(t, actual, expected)
}

func TestInstallInvalidAppLauncher(t *testing.T) {

	yml := trim(`
		app-launcher: Unknown
		terminal: None
		keyboard-layout: None
		ides: [ ]
		additional-options: [ ]
		blacklist: [ ]`,
	)

	_, _, err := runInstaller(t, &yml)

	AssertErrorContains(t, err, `Invalid param 'app-launcher' value/s 'unknown', valid values:
		spotlight
		launchpad
		alfred
		none`)
}

func TestInstallInvalidTerminal(t *testing.T) {

	yml := trim(`
		app-launcher: None
		terminal: unknown
		keyboard-layout: None
		ides: [ ]
		additional-options: [ ]
		blacklist: [ ]`,
	)

	_, _, err := runInstaller(t, &yml)

	AssertErrorContains(t, err, `Invalid param 'terminal' value/s 'unknown', valid values:
		default
		iterm
		warp
		none`)
}

func TestInstallInvalidKeyboardLayout(t *testing.T) {

	yml := trim(`
		app-launcher: None
		terminal: None
		keyboard-layout: unknown
		ides: [ ]
		additional-options: [ ]
		blacklist: [ ]`,
	)

	_, _, err := runInstaller(t, &yml)

	AssertErrorContains(t, err, `Invalid param 'keyboard-layout' value/s 'unknown', valid values:
		pc
		mac`)
}

func TestInstallAdditionalOptions(t *testing.T) {

	yml := trim(`
		app-launcher: alfred
		terminal: warp
		keyboard-layout: pc
		ides: [ ]
		additional-options:
		- Enable Dock auto-hide (2s delay)
		- Change Dock minimize animation to "scale"
		- Enable Home & End keys
		- Show hidden files in Finder
		- Show directories on top in Finder
		- Show full POSIX paths in Finder window title
		blacklist: [ com.apple.Preview ]`,
	)

	i, c, _ := runInstaller(t, &yml)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-alfred-warp-pc.json"

	AssertFilesEqual(t, actual, expected)
	AssertSlicesEqual(t, c.CommandsLog, []string{
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"plutil -convert binary1 /homedir/Library/preferences/com.knollsoft.Rectangle.plist",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"plutil -convert binary1 /homedir/Library/preferences/com.lwouis.alt-tab-macos.plist",
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

func runInstaller(t *testing.T, yml *string) (install.HomeDir, MockCommander, error) {
	commander := MockCommander{}
	homeDir := testHomeDir()
	err := install.RunInstaller(homeDir, &commander, FakeTimeProvider{}, yml)
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
	removeFiles(i.KarabinerConfigBackupFile(FakeTimeProvider{}.Now()))
	removeFiles(i.IdesKeymapPaths(install.IDEKeymaps)...)
	removeFilesWithExt(i.LibraryDir(), "plist")
	removeFilesWithExt(i.LibraryDir(), "dict")
	copyFile(karabinerTestDefaultConfig(i), i.KarabinerConfigFile())
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // reset flags
}

func karabinerTestDefaultConfig(i install.HomeDir) string {
	return i.KarabinerConfigDir() + "/karabiner-default.json"
}
