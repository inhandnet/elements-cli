package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"github.com/Sirupsen/logrus"
	"strings"
)

var sess *mgo.Session

func Connect(url string) (err error) {
	sess, err = mgo.Dial(url)
	return
}

func S(action func(s *mgo.Session)) {
	s := sess.Copy()
	defer s.Close()
	action(s)
}

func Session() *mgo.Session {
	return sess
}

func UserOid(email string) (bson.ObjectId, error) {
	r := bson.M{}
	sess.DB("ABCDE_db").C("user_dbs").Find(bson.M{
		"username": email,
	}).One(&r)

	if oid, ok := r["oid"]; ok {
		return oid.(bson.ObjectId), nil
	}
	return "", errors.New("not found")
}

func FindDevice(serialNumber, email string) []bson.M {
	result := make([]bson.M, 100)
	query := bson.M{}
	if strings.HasPrefix(serialNumber, "/") && strings.HasSuffix(serialNumber, "/") {
		query["sn"] = bson.RegEx{Pattern: serialNumber[1: len(serialNumber)-1]}
	} else {
		query["sn"] = serialNumber
	}

	if email != "" {
		if oid, err := UserOid(email); err != nil {
			logrus.Fatalln("user not found", email)
		} else {
			query["oid"] = oid
		}
	}

	sess.DB("dn_pp").C("ovdp_device").Find(query).All(&result)
	return result
}

func FindSites(name, email string) []bson.M {
	result := make([]bson.M, 100)
	query := bson.M{}
	if strings.HasPrefix(name, "/") && strings.HasSuffix(name, "/") {
		query["name"] = bson.RegEx{Pattern: name[1: len(name)-1]}
	} else {
		query["name"] = name
	}

	oid, err := UserOid(email)
	if err != nil {
		logrus.Fatalln("user not found", email)
	}

	sess.DB(UserDb(oid)).C("site").Find(query).All(&result)
	return result
}

func UserDb(oid bson.ObjectId) string {
	return strings.ToUpper(oid.Hex()) + "_db"
}
