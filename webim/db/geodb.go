package db

import "github.com/oschwald/geoip2-golang"

var GeoDB *geoip2.Reader

func LoadGeoDB(dbPath string) (err error) {
	GeoDB, err = geoip2.Open(dbPath)
	return
}
