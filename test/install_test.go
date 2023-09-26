package install_test

import (
	"flag"
	"github.com/raxigan/pcfy-my-mac/install"
	"io"
	"os"
	"strings"
	"testing"
)

func TestInstallFromYamlFile(t *testing.T) {

	os.Args = []string{"script_name", "--params=params.yml"}
	i, _, _ := runInstaller(nil)
	defer reset(i)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-warp-alfred-pc.json"

	AssertFilesEqual(t, actual, expected)

	removeFiles(i.KarabinerConfigBackupFile())
}

func TestInstallWarpAlfredPC(t *testing.T) {

	yml := yaml(`
		app-launcher: alfred
		terminal: warp
		keyboard-type: pc
		ides: [ ]
		additional-options: [ ]
		blacklist: [ ]`,
	)

	i, c, _ := runInstaller(&yml)
	defer reset(i)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-warp-alfred-pc.json"

	AssertFilesEqual(t, actual, expected)
	AssertSlicesEqual(t, c.CommandsLog, []string{
		"clear",
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"defaults read com.lwouis.alt-tab-macos.plist",
		"open -a AltTab",
	})
}

func TestInstallItermSpotlightMac(t *testing.T) {

	yml := yaml(`
		app-launcher: spotlight
		terminal: iterm
		keyboard-type: mac
		ides: [ ]
		additional-options: [ ]
		blacklist: [ ]`,
	)

	i, c, _ := runInstaller(&yml)
	defer reset(i)

	actual := i.KarabinerConfigFile()
	expected := "expected/karabiner-expected-iterm-spotlight-mac.json"

	AssertFilesEqual(t, actual, expected)
	AssertSlicesEqual(t, c.CommandsLog, []string{
		"clear",
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"defaults read com.lwouis.alt-tab-macos.plist",
		"open -a AltTab",
	})
}

func TestInstallAllKeymaps(t *testing.T) {

	yml := yaml(`
		app-launcher: None
		terminal: None
		keyboard-type: None
		ides: [ "all" ]
		additional-options: [ ]
		blacklist: [ ]`,
	)

	i, c, _ := runInstaller(&yml)
	defer reset(i)

	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJ()), i.IdeKeymapPaths(install.IntelliJ())[0])
	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJ()), i.IdeKeymapPaths(install.IntelliJ())[1])
	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.IntelliJCE()), i.IdeKeymapPaths(install.IntelliJCE())[0])
	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.GoLand()), i.IdeKeymapPaths(install.GoLand())[0])
	AssertFilesEqual(t, "../configs/"+i.SourceKeymap(install.Fleet()), i.IdeKeymapPaths(install.Fleet())[0])
	AssertSlicesEqual(t, c.CommandsLog, []string{
		"clear",
		"killall Karabiner-Elements",
		"open -a Karabiner-Elements",
		"killall Rectangle",
		"defaults read com.knollsoft.Rectangle.plist",
		"open -a Rectangle",
		"killall AltTab",
		"defaults read com.lwouis.alt-tab-macos.plist",
		"open -a AltTab",
	})
}

func TestFailForUnknownParam(t *testing.T) {

	yml := yaml(`unknown: hello`)

	i, c, err := runInstaller(&yml)
	defer reset(i)

	AssertErrorContains(t, err, "Unknown parameter: unknown")
	AssertSlicesEqual(t, c.CommandsLog, []string{})
}

func TestFailForInvalidYaml(t *testing.T) {

	yml := yaml(`[] :app-launcher:`)

	i, c, err := runInstaller(&yml)
	defer reset(i)

	AssertErrorContains(t, err, "cannot unmarshal !!seq into install.FileParams")
	AssertSlicesEqual(t, c.CommandsLog, []string{})
}

func TestInstallYmlFileDoesNotExist(t *testing.T) {

	os.Args = []string{"script_name", "--params=nope.yml"}
	i, c, err := runInstaller(nil)
	defer reset(i)

	AssertErrorContains(t, err, "open nope.yml: no such file or directory")
	AssertSlicesEqual(t, c.CommandsLog, []string{})
}

func runInstaller(yml *string) (install.Installation, install.MockCommander, error) {
	wd, _ := os.Getwd()
	commander := install.MockCommander{}
	installer, err := install.RunInstaller(wd+"/homedir", &commander, yml)
	return installer, commander, err
}

func reset(i install.Installation) {
	removeFiles(i.KarabinerConfigBackupFile())
	copyFile(karabinerTestDefaultConfig(i), i.KarabinerConfigFile())
	removeFiles(i.IdesKeymapPaths(install.IDEKeymaps)...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func karabinerTestDefaultConfig(i install.Installation) string {
	return i.KarabinerConfigDir() + "/karabiner-default.json"
}

func removeFiles(paths ...string) {
	for _, path := range paths {
		os.Remove(path)
	}
}

func copyFile(src, dst string) {
	sourceFile, _ := os.Open(src)
	defer sourceFile.Close()

	destFile, _ := os.Create(dst)
	defer destFile.Close()

	io.Copy(destFile, sourceFile)
}

func yaml(yaml string) string {
	return strings.TrimSpace(strings.ReplaceAll(yaml, "\t", ""))
}
