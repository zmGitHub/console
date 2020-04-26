package common

import "github.com/deckarep/golang-set"

var MunicipalityList = []interface{}{
	"北京",
	"天津",
	"上海",
	"重庆",
}

var MunicipalitySet = mapset.NewSetFromSlice(MunicipalityList)

func IsMunicipality(city string) bool {
	return MunicipalitySet.Contains(city)
}
