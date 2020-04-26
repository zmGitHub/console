package db

import (
	"github.com/ipipdotnet/ipdb-go"

	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var IPCity *ipdb.City

type IPGeo struct {
	city *ipdb.City
}

func (geo *IPGeo) GetLocation(ip string) (country, province, city, isp string, err error) {
	info, err := geo.city.FindInfo(ip, "CN")
	if err != nil {
		return "", "", "", "", err
	}

	country = info.CountryName
	province = info.RegionName
	city = info.CityName
	isp = info.IspDomain
	return
}

func LoadIPGeoDB(config *conf.IPGeoConfig) (ipGeo *IPGeo, err error) {
	IPCity, err := ipdb.NewCity(config.IPDBPath)
	if err != nil {
		return nil, err
	}

	return &IPGeo{city: IPCity}, nil
}
