package common

import "github.com/rs/xid"

func GenUniqueID() string {
	return xid.New().String()
}
