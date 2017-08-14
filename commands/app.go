package commands

import (
	"github.com/urfave/cli"
	"github.com/Droi-SDK/droi-checker/version"
)

func Run(args []string) {
	app := cli.NewApp()
	app.Name = "droi-checker"
	app.Version = version.Version
	app.Usage = "Command line to check DroiBaaS integrate in android/ios"
	app.EnableBashCompletion = true
	app.Action = checkAction
	app.Run(args)
	/*app.Commands = []cli.Command{
		{
			Name:   "android",
			Usage:  "查看当前登录用户以及应用信息",
			Action: androidChecker(),
		},
		{
			Name:   "info",
			Usage:  "查看当前登录用户以及应用信息",
			Action: iosChecker(),
		},
	}*/

}
