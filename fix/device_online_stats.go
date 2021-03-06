package fix

import (
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/urfave/cli"
	"github.com/inhandnet/elements-cli/mongo"
)

func MigrateDeviceOnlineEvents(c *cli.Context) {
	iter := mongo.Session().DB("ABCDE_db").C("organizations").Find(bson.M{}).Iter()

	org := new(Org)
	for iter.Next(org) {
		db := org.DbName()
		logrus.Infoln("Running in", db)
		migrate(db, c.Bool("retain"))
	}

}

type Org struct {
	Id   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name,omitempty"`
}

func (o *Org) DbName() string {
	var hex string

	if o.Id.Hex() == "0000000000000000000abcde" {
		hex = "abcde"
	} else {
		hex = o.Id.Hex()
	}
	return strings.ToUpper(hex) + "_db"
}

func migrate(db string, retain bool) {
	conn := mongo.Session().Copy()
	defer conn.Close()
	coll := conn.DB(db).C("device.online.events")

	coll.DropIndex("time", "deviceId")
	coll.EnsureIndex(mgo.Index{
		Background: true,
		Key:        []string{"time"},
	})
	coll.EnsureIndex(mgo.Index{
		Background: true,
		Key:        []string{"deviceId"},
	})

	oldc := conn.DB(db).C("device_online_stat")
	total, _ := oldc.Find(bson.M{}).Count()
	iter := oldc.Find(bson.M{}).Iter()

	doc := new(DeviceOnlineStat)
	bulk := coll.Bulk()
	bulk.Unordered()
	i := 0
	for iter.Next(doc) {
		bulk.Upsert(bson.M{
			"deviceId": doc.DeviceId,
			"time":     time.Unix(doc.Login, 0),
			"type":     1,
		}, bson.M{
			"$set": bson.M{"exception": false},
		})

		if doc.Logout != 0 {
			bulk.Upsert(bson.M{
				"deviceId": doc.DeviceId,
				"time":     time.Unix(doc.Logout, 0),
				"type":     0,
			}, bson.M{
				"$set": bson.M{"exception": doc.Exception == 1},
			})
		}

		i++
		if i%400 == 0 {
			if _, err := bulk.Run(); err != nil {
				logrus.Fatalln(err)
			}
			bulk = coll.Bulk()
			bulk.Unordered()
		}

		if i%10000 == 0 || i == total {
			logrus.Infof("Progress: %v/%v in %v finished.", i, total, db)
		}

	}
	if _, err := bulk.Run(); err != nil {
		logrus.Fatalln(err)
	}
	if !retain {
		oldc.DropCollection()
	}
}

type DeviceOnlineStat struct {
	DeviceId  bson.ObjectId `bson:"deviceId,omitempty"`
	Login     int64         `bson:"login,omitempty"`
	Logout    int64         `bson:"logout,omitempty"`
	Exception int           `bson:"exception,omitempty"`
}
