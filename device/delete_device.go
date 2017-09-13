package device

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/inhandnet/elements-cli/mongo"
	"github.com/urfave/cli"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/resty.v0"
	"strings"
	"time"
)

var (
	token   = ""
	baseUrl = ""
)

func Prepare(c *cli.Context) (err error) {
	baseUrl = c.GlobalString("addr")
	password := c.GlobalString("password")
	username := c.GlobalString("username")

	token, err = adminToken(username, password)
	if err != nil {
		logrus.Fatalln("failed to init admin token:", err.Error())
	}

	time.AfterFunc(30*time.Minute, func() {
		Prepare(c)
	})
	return
}

func Delete(oid, id string, deleteSite bool) error {
	params := map[string]string{
		"access_token": token,
		"oid":          oid,
	}
	if deleteSite {
		params["delete_site"] = "1"
	}
	resp, _ := resty.R().SetQueryParams(params).SetResult(new(bson.M)).Delete(baseUrl + "/api/devices/" + id)

	result := resp.Result().(*bson.M)
	if _, ok := (*result)["error"]; ok {
		return errors.New((*result)["error"].(string))
	}
	return nil
}

func Find(serialNumber, email string) []bson.M {
	result := make([]bson.M, 100)
	query := bson.M{}
	if strings.HasPrefix(serialNumber, "/") && strings.HasSuffix(serialNumber, "/") {
		query["sn"] = bson.RegEx{Pattern: serialNumber[1: len(serialNumber)-1]}
	} else {
		query["sn"] = serialNumber
	}

	if email != "" {
		if oid, err := mongo.UserOid(email); err != nil {
			logrus.Fatalln("user not found", email)
		} else {
			query["oid"] = oid
		}
	}

	mongo.Session().DB("dn_pp").C("ovdp_device").Find(query).All(&result)
	return result
}

func adminToken(username, password string) (token string, err error) {
	resp, err := resty.R().SetFormData(map[string]string{
		"client_id":     "17953450251798098136",
		"client_secret": "08E9EC6793345759456CB8BAE52615F3",
		"grant_type":    "password",
		"username":      username,
		"password":      password,
		"password_type": "1",
	}).SetResult(new(bson.M)).Post(baseUrl + "/oauth2/access_token")
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
