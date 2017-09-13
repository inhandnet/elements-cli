package client

import (
	"github.com/Sirupsen/logrus"
	"time"
	"github.com/urfave/cli"
	"gopkg.in/resty.v0"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

var (
	Token   = ""
	BaseUrl = ""
)


func Prepare(c *cli.Context) (err error) {
	BaseUrl = c.GlobalString("addr")
	password := c.GlobalString("password")
	username := c.GlobalString("username")

	f := func() {
		Token, err = authenticate(username, password)
		if err != nil {
			logrus.Fatalln("failed to init admin Token:", err.Error())
		}
	}
	time.AfterFunc(30*time.Minute, f)
	f()
	return
}


func authenticate(username, password string) (token string, err error) {
	resp, err := resty.R().SetFormData(map[string]string{
		"client_id":     "17953450251798098136",
		"client_secret": "08E9EC6793345759456CB8BAE52615F3",
		"grant_type":    "password",
		"username":      username,
		"password":      password,
		"password_type": "1",
	}).SetResult(new(bson.M)).Post(BaseUrl + "/oauth2/access_token")
	if err != nil {
		logrus.Fatalln(err.Error())
		return
	}

	result := resp.Result().(*bson.M)

	if _, ok := (*result)["access_token"]; ok {
		token = (*result)["access_token"].(string)
	} else {
		err = errors.New((*result)["error"].(string))
	}
	return
}
