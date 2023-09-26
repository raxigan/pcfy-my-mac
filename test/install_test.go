package install_test

import (
	"github.com/raxigan/pcfy-my-mac/install"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestInstallFromYamlFile(t *testing.T) {

	os.Args = []string{"script_name", "--params=params.yml"}
	i, _ := runInstaller(nil)

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

	i, c := runInstaller(&yml)
	defer resetHomeDir(i)

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

	i, c := runInstaller(&yml)
	defer resetHomeDir(i)

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

	i, c := runInstaller(&yml)
	defer resetHomeDir(i)

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

func runInstaller(yml *string) (install.Installation, install.MockCommander) {
	wd, _ := os.Getwd()
	commander := install.MockCommander{}
	return install.RunInstaller(wd+"/homedir", &commander, yml), commander
}

func resetHomeDir(i install.Installation) {
	removeFiles(i.KarabinerConfigBackupFile())
	copyFile(karabinerTestDefaultConfig(i), i.KarabinerConfigFile())
	removeFiles(i.IdesKeymapPaths(install.IDEKeymaps)...)
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
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func yaml(yaml string) string {
	return strings.TrimSpace(strings.ReplaceAll(yaml, "\t", ""))
}
