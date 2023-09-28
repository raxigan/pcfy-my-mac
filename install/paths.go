package install

import (
	"os"
	"path/filepath"
	"time"
)

type HomeDir struct {
	Path string
}

func DefaultHomeDir() HomeDir {
	homeDirPath, _ := os.UserHomeDir()
	return HomeDir{Path: homeDirPath}
}

func (home HomeDir) KarabinerConfigDir() string {
	return home.Path + "/.config/karabiner"
}

func (home HomeDir) KarabinerConfigFile() string {
	return filepath.Join(home.KarabinerConfigDir(), "/karabiner.json")
}

func (home HomeDir) KarabinerConfigBackupFile(time time.Time) string {
	currentTime := time.Format("02-01-2006_15:04:05")
	return home.KarabinerConfigDir() + "/karabiner-" + currentTime + ".json"
}

func (home HomeDir) KarabinerComplexModificationsDir() string {
	return home.Path + "/.config/karabiner/assets/complex_modifications"
}

func (home HomeDir) ApplicationSupportDir() string {
	return home.Path + "/Library/Application Support"
}

func (home HomeDir) PreferencesDir() string {
	return home.Path + "/Library/preferences"
}

func (home HomeDir) LibraryDir() string {
	return home.Path + "/Library"
}

func (home HomeDir) SourceKeymap(ide IDE) string {
	return "keymaps/" + ide.srcKeymapsFile
}
