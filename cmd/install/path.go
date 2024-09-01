package install

import (
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"os"
	"path/filepath"
	"strings"
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
	return filepath.Join(home.Path, ".config/karabiner")
}

func (home HomeDir) KarabinerConfigFile() string {
	return filepath.Join(home.KarabinerConfigDir(), "karabiner.json")
}

func (home HomeDir) KarabinerConfigBackupFile(time time.Time) string {
	currentTime := time.Format("02-01-2006_15:04:05")
	return filepath.Join(home.KarabinerConfigDir(), "/karabiner-"+currentTime+".json")
}

func (home HomeDir) KarabinerComplexModificationsDir() string {
	return filepath.Join(home.Path, ".config/karabiner/assets/complex_modifications")
}

func (home HomeDir) ApplicationSupportDir() string {
	return filepath.Join(home.Path, "Library/Application Support")
}

func (home HomeDir) PreferencesDir() string {
	return filepath.Join(home.Path, "Library/Preferences")
}

func (home HomeDir) LibraryDir() string {
	return filepath.Join(home.Path, "Library")
}

func (home HomeDir) LaunchAgents() string {
	return filepath.Join(home.LibraryDir(), "LaunchAgents")
}

func (home HomeDir) SourceKeymap(ide param.IDE) string {
	return filepath.Join("keymaps/", ide.SrcKeymapsFile)
}

func (home HomeDir) IdeKeymapPaths(ide param.IDE) []string {
	return home.IdesKeymapPaths([]param.IDE{ide})
}

func (home HomeDir) IdesKeymapPaths(ide []param.IDE) []string {

	var result []string

	for _, e := range ide {

		keymapDest := initializeWithDefault(e.DestKeymapsFile, formatString(e.FullName)+".xml")
		dirs, _ := common.FindMatchingPaths(filepath.Join(home.Path, e.KeymapsDir), keymapDest)

		for _, e1 := range dirs {
			result = append(result, e1)
		}
	}

	return result
}

func initializeWithDefault(original, defaultValue string) string {
	if original == "" {
		return defaultValue
	}
	return original
}

func formatString(input string) string {
	lowercased := strings.ToLower(input)
	formatted := strings.ReplaceAll(lowercased, " ", "-")
	return formatted
}
