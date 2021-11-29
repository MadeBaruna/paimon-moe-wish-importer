package main

import (
	"errors"
	"net/url"
)

func GetUrl(val string, server string) (*url.URL, error) {
	u, err := url.Parse(val)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	u.Path = "event/gacha_info/api/getGachaLog"
	u.Host = "hk4e-api-os.mihoyo.com"
	if server == "china" {
		u.Host = "hk4e-api.mihoyo.com"
	}
	u.Fragment = ""

	return u, nil
}
