package cmd

import (
	"github.com/raxigan/pcfy-my-mac/cmd/common"
	"github.com/raxigan/pcfy-my-mac/cmd/install"
	"github.com/raxigan/pcfy-my-mac/cmd/param"
	"github.com/raxigan/pcfy-my-mac/cmd/task"
)

func RunInstaller(homeDir install.HomeDir, commander install.Commander, tp install.TimeProvider, params param.Params) error {

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
		task.ApplyMainKarabinerRules(),
		task.ApplyAppLauncherRules(),
		task.ApplyKeyboardLayoutRules(),
		task.ApplyTerminalRules(),
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

	for _, task := range tasks {
		i.Commander.TryPrint(common.Colored(common.Blue, "TASK"), task.Name)

		err := task.Exec(i)

		if err != nil {
			return err
		}
	}

	common.PrintlnColored(common.Cyan, "PC'fied!")

	return nil
}
