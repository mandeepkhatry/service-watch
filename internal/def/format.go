package def

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
)

var StringFormat = map[string]string{
	"ipv4":      randomdata.IpV4Address(),
	"ipv6":      randomdata.IpV6Address(),
	"date-time": "2018-11-13T20:20:39+00:00",
	"time":      "20:20:39+00:00",
	"date":      "2018-11-13",
	"email":     "abde@xyz.com",
	"uuid":      uuid.New().String(),
}
