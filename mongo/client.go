package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

var sess *mgo.Session

func Connect(url string) (err error) {
	sess, err = mgo.Dial(url)
	return
}

func S(action func(s *mgo.Session))  {
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