package device

import (
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/httplib"
	"errors"
	"github.com/urfave/cli"
	"gopkg.in/resty.v0"
	"gopkg.in/mgo.v2/bson"
)

var (
	token   = ""
	baseUrl = ""
)

func Prepare(c *cli.Context) {
	baseUrl = c.String("url")
	password := c.String("password")

	var err error
	token, err = adminToken(password)
	if err != nil {
		logrus.Fatalln("request token failed, " + err.Error())
	}
}

func Delete(oid, id string) error {
	resp, _ := resty.R().SetQueryParams(map[string]string{
		"access_token": token,
		"oid": oid,
	}).SetResult(new(bson.M)).Delete(baseUrl + "/api/devices/" + id)

	result := resp.Result().(*bson.M)
	if _, ok := (*result)["error"]; ok {
		return errors.New((*result)["error"].(string))
	}
	return nil
}

func Find(serialNumber string) map[string]interface{} {
	resp := map[string]interface{}{}
	req := httplib.Get(baseUrl + "/api/devices")
	req.Param("serial_number", serialNumber)
	req.Param("access_token", token)
	req.Param("verbose", "50")
	req.ToJSON(&resp)

	if _, ok := resp["error"]; ok {
		return nil
	}

	if result, ok := resp["result"].([]interface{}); ok {
		if len(result) > 0 {
			return result[0].(map[string]interface{})
		}
	}
	return nil
}

func adminToken(password string) (token string, err error) {
	resp := map[string]interface{}{}
	req := httplib.Post(baseUrl + "/oauth2/access_token")

	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Param("client_id", "17953450251798098136")
	req.Param("client_secret", "08E9EC6793345759456CB8BAE52615F3")
	req.Param("grant_type", "password")
	req.Param("username", "admin")
	req.Param("password", password)
	req.Param("password_type", "1")
	req.ToJSON(&resp)

	if err != nil {
		return
	}
	if _, ok := resp["access_token"]; ok {
		token = resp["access_token"].(string)
	} else {
		err = errors.New(resp["error"].(string))
	}
	return
}
