package main

import (
	"os"
	"github.com/codegangsta/cli"
	"github.com/inhandnet/elements-cli/fix"
)

func main() {
	NewApp().Run(os.Args)
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "elements util"
	app.Usage = "elements scripts utility."

	app.Commands = []cli.Command{
		{
			Name:  "fix",
			Usage: "fix mongodb documents",
			Subcommands: []cli.Command{
				{
					Name:   "migrate-online-stats",
					Usage:  "migrate device_oniline_stat to device.online.stats",
					Flags:  []cli.Flag{
						cli.StringFlag{
							Name:  "url",
							Value: "mongodb://admin:admin@localhost:27017/",
							Usage: "mongodb connect uri",
						},
						cli.BoolFlag{
							Name: "retain",
							Usage: "retain collection device_online_stat after migrate, default it will be dropped",
						},
					},
					Action: fix.MigrateDeviceOnlineEvents,
				},
			},
		},
	}

	cli.HelpFlag.Name = "help"
	return app
}
