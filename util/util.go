package util

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func ObjectIdStr(v interface{}) (str string) {
	switch v.(type) {
	case string:
		str = v.(string)
	case bson.ObjectId:
		str = v.(bson.ObjectId).Hex()
	case fmt.Stringer:
		str = v.(fmt.Stringer).String()
	}
	return
}