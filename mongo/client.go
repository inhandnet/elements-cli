package mongo

import (
	"gopkg.in/mgo.v2"
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