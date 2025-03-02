package cmd

import (
	"fmt"
	"github.com/raxigan/pcfy-my-mac/cmd/install"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"github.com/raxigan/pcfy-my-mac/cmd/task"
)

func Launch(homeDir install.HomeDir, commander install.Commander, tp install.TimeProvider, params param.Params) error {

	installation := install.Installation{
		Commander:        commander,
		HomeDir:          homeDir,
		Params:           params,
		ProfileName:      "PCfy",
		InstallationTime: tp.Now(),
	}

	return Install(installation)
}

func Install(i install.Installation) error {

	installDeps := task.DownloadDependencies()
	err := installDeps.Execute(i)
	if err != nil {
		return err
	}

	tasks := []task.Task{
		task.CloseKarabiner(),
		task.BackupKarabinerConfig(),
		task.DeleteExistingKarabinerProfile(),
		task.CreateKarabinerProfile(),
		task.NameKarabinerProfile(),
		task.UnselectOtherKarabinerProfiles(),
		task.ApplyTerminalRules(),
		task.ApplyMainKarabinerRules(),
		task.ApplyAppLauncherRules(),
		task.ApplyKeyboardLayoutRules(),
		task.ReformatKarabinerConfigFile(),
		task.OpenKarabiner(),
		task.CopyIdeKeymaps(),
		task.CloseRectangle(),
		task.CopyRectanglePreferences(),
		task.OpenRectangle(),
		task.CloseAltTab(),
		task.InstallAltTabPreferences(),
		task.OpenAltTab(),
		task.ApplySystemSettings(),
		task.CopyHidutilRemappingFile(),
		task.ExecuteHidutil(),
	}

	for _, t := range tasks {
		i.Commander.Progress()
		i.Commander.TryLog(install.TaskMsg, t.Name)

		err := t.Execute(i)

		if err != nil {
			return err
		}
	}

	i.TryLog(install.TaskMsg, "Installed successfully")

	fmt.Println("PC'fied")
	i.Commander.Run("clear")
	fmt.Println(`
Almost ready!

1. Restart the tools (if any) you installed the keymaps for, and then select
   the new keymap "PCfy" in settings.
2. Grant appropriate system permissions to the following tools when prompted:
 • Karabiner-Elements
 • Alt-Tab
 • Rectangle`)

	return nil
}
