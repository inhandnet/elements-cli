package main

import (
	"os"
	"github.com/inhandnet/elements-cli/fix"
	"github.com/inhandnet/elements-cli/device"
	"github.com/urfave/cli"
	"github.com/Sirupsen/logrus"
	"github.com/inhandnet/elements-cli/log"
)

func main() {
	NewApp().Run(os.Args)
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "elements util"
	app.Usage = "elements scripts utility."
	app.Version = "0.1.1"
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
		{
			Name:  "device",
			Usage: "device management",
			Subcommands: []cli.Command{
				{
					Name:   "delete",
					Usage:  "delete device",
					Flags:  []cli.Flag{
						cli.StringFlag{
							Name:  "url, h",
							Value: "http://10.5.16.105",
							Usage: "api base url",
						},
						cli.StringFlag{
							Name:  "password, p",
							Value: "admin",
							Usage: "admin password",
						},
						cli.StringFlag{
							Name: "serialNumber, s",
							Usage: "device serial number to delete",
						},
					},
					Action: func(c *cli.Context) {
						serialNumber := c.String("serialNumber")
						device.Prepare(c)

						d := device.Find(serialNumber)
						if d == nil {
							logrus.Fatalln("device not found")
						}

						id := d["_id"].(string)
						oid := d["oid"].(string)

						if err := device.Delete(oid, id); err != nil {
							logrus.Fatalln(err.Error())
						}
						log.PrintJSON(d)
					},
				},
				{
					Name:   "find",
					Usage:  "find device by serial number",
					Flags:  []cli.Flag{
						cli.StringFlag{
							Name:  "url, h",
							Value: "http://10.5.16.105",
							Usage: "api base url",
						},
						cli.StringFlag{
							Name:  "password, p",
							Value: "admin",
							Usage: "admin password",
						},
						cli.StringFlag{
							Name: "serialNumber, s",
							Usage: "device serial number to delete",
						},
					},
					Action: func(c *cli.Context) {
						serialNumber := c.String("serialNumber")
						device.Prepare(c)

						d := device.Find(serialNumber)
						if d == nil {
							logrus.Fatalln("device not found")
						}

						log.PrintJSON(d)
					},
				},
			},
		},
	}

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help",
		Usage: "show help",
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "print the version",
	}

	return app
}
