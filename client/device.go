package client

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/resty.v0"
)

func DeleteDevice(oid, id string, deleteSite bool) error {
	params := map[string]string{
		"access_token": Token,
		"oid":          oid,
	}
	if deleteSite {
		params["delete_site"] = "1"
	}
	resp, _ := resty.R().SetQueryParams(params).SetResult(new(bson.M)).Delete(BaseUrl + "/api/devices/" + id)

	result := resp.Result().(*bson.M)
	if _, ok := (*result)["error"]; ok {
		return errors.New((*result)["error"].(string))
	}
	return nil
}


func DeleteSite(oid, id string) error {
	params := map[string]string{
		"access_token": Token,
		"oid":          oid,
	}
	resp, _ := resty.R().SetQueryParams(params).SetResult(new(bson.M)).Delete(BaseUrl + "/api/sites/" + id)

	result := resp.Result().(*bson.M)
	if _, ok := (*result)["error"]; ok {
		return errors.New((*result)["error"].(string))
	}
	return nil
}
