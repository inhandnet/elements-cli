package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/inhandnet/elements-cli/fix"
	"github.com/inhandnet/elements-cli/log"
	"github.com/urfave/cli"
	"os"
	"github.com/inhandnet/elements-cli/mongo"
	"github.com/inhandnet/elements-cli/util"
	"github.com/jeffail/tunny"
	"sync"
	"gopkg.in/mgo.v2/bson"
	"github.com/inhandnet/elements-cli/client"
)

func main() {
	NewApp().Run(os.Args)
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "elements util"
	app.Usage = "elements scripts utility."
	app.Version = "0.1.3"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "mongo-uri, m",
			Value: "mongodb://admin:admin@localhost:27017/",
			Usage: "mongodb connect uri",
		},
		cli.StringFlag{
			Name:  "addr, a",
			Value: "http://localhost",
			Usage: "api base url",
		},
		cli.StringFlag{
			Name:  "username, u",
			Value: "admin",
			Usage: "platform account username",
		},
		cli.StringFlag{
			Name:  "password, p",
			Value: "admin",
			Usage: "platform account password",
		},
	}
	app.Before = func(ctx *cli.Context) (err error) {
		if uri := ctx.GlobalString("mongo-uri"); uri != "" {
			err = mongo.Connect(uri)
		}
		return
	}
	app.Commands = []cli.Command{
		{
			Name:  "fix",
			Usage: "fix mongodb documents",
			Subcommands: []cli.Command{
				{
					Name:  "migrate-online-stats",
					Usage: "migrate device_oniline_stat to device.online.stats",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "retain",
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
					Before: client.Prepare,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "serialNumber, s",
							Usage: "device serial number to delete",
						},
						cli.IntFlag{
							Name:  "pool",
							Usage: "worker pool size",
							Value: 1,
						},
						cli.BoolFlag{
							Name:  "delete-site",
							Usage: "also delete site",
						},
						cli.StringFlag{
							Name:  "email",
							Usage: "only delete device of user",
						},
					},
					Action: func(c *cli.Context) {
						serialNumber := c.String("serialNumber")

						list := mongo.FindDevice(serialNumber, c.String("email"))
						if len(list) == 0 {
							logrus.Fatalln("device not found")
						}

						log.PrintJSON(list)

						pool, _ := tunny.CreatePool(c.Int("pool"), func(i interface{}) interface{} {
							d := i.(bson.M)
							id := util.ObjectIdStr(d["deviceId"])
							oid := util.ObjectIdStr(d["oid"])
							logrus.Infoln("Deleting", d["sn"], "in", oid)
							if err := client.DeleteDevice(oid, id, c.Bool("delete-site")); err != nil {
								logrus.Fatalln(err.Error())
							}

							return nil
						}).Open()
						wg := sync.WaitGroup{}

						wg.Add(len(list))
						for _, d := range list {
							pool.SendWorkAsync(d, func(i interface{}, e error) {
								wg.Done()
							})
						}
						wg.Wait()
					},
				},
				{
					Name:  "find",
					Usage: "find device by serial number",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "serialNumber, s",
							Usage: "device serial number to delete",
						},
						cli.StringFlag{
							Name:  "email",
							Usage: "only delete device of user",
						},
					},
					Action: func(c *cli.Context) {
						serialNumber := c.String("serialNumber")

						d := mongo.FindDevice(serialNumber, c.String("email"))
						if d == nil {
							logrus.Fatalln("device not found")
						}

						log.PrintJSON(d)
					},
				},
			},
		},
		{
			Name:  "site",
			Usage: "site management",
			Subcommands: []cli.Command{
				{
					Name:   "delete",
					Usage:  "search and delete sites",
					Before: client.Prepare,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "site name, use /abc/ to search regex",
						},
						cli.IntFlag{
							Name:  "pool",
							Usage: "worker pool size",
							Value: 1,
						},
						cli.StringFlag{
							Name:  "email",
							Usage: "only delete device of user",
						},
					},
					Action: func(c *cli.Context) {
						list := mongo.FindSites(c.String("name"), c.String("email"))
						if len(list) == 0 {
							logrus.Fatalln("sites not found")
						}

						log.PrintJSON(list)

						pool, _ := tunny.CreatePool(c.Int("pool"), func(i interface{}) interface{} {
							d := i.(bson.M)
							id := util.ObjectIdStr(d["_id"])
							oid := util.ObjectIdStr(d["oid"])
							logrus.Infoln("Deleting", d["name"], "in", oid)
							if err := client.DeleteSite(oid, id); err != nil {
								logrus.Fatalln(err.Error())
							}

							return nil
						}).Open()
						wg := sync.WaitGroup{}

						wg.Add(len(list))
						for _, d := range list {
							pool.SendWorkAsync(d, func(i interface{}, e error) {
								wg.Done()
							})
						}
						wg.Wait()
					},
				},
				{
					Name:  "find",
					Usage: "find sites by name",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "serialNumber, s",
							Usage: "device serial number to delete",
						},
						cli.StringFlag{
							Name:  "email",
							Usage: "only delete device of user",
						},
					},
					Action: func(c *cli.Context) {
						d := mongo.FindSites(c.String("name"), c.String("email"))
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
