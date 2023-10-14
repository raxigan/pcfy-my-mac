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

	tasks := []task.Task{
		task.DownloadDependencies(),
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
	}

	for _, t := range tasks {
		i.Commander.TryLog(install.TaskMsg, t.Name)

		err := t.Execute(i)

		if err != nil {
			return err
		}
	}

	fmt.Println("PC'fied!")

	return nil
}
