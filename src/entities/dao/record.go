package dao

import (
	"time"
)

type CommonRecordDAO struct {
	Time       time.Time   `bson:"time" json:"time"`
	RemoteAddr string      `bson:"remote_addr" json:"remote_addr"`
	Data       interface{} `bson:"data" json:"data"`
}

type IdentifierTupleDAO struct {
	InstanceID string    `bson:"instance_id" json:"instance_id"`
	HostID     string    `bson:"host_id" json:"host_id"`
	IP         string    `bson:"ip" json:"ip"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

type IPGeoInfoDAO struct {
	Status      string    `bson:"status" json:"status"`
	IP          string    `bson:"ip" json:"query"`
	Region      string    `bson:"region" json:"region"`
	RegionName  string    `bson:"region_name" json:"regionName"`
	City        string    `bson:"city" json:"city"`
	Country     string    `bson:"country" json:"country"`
	CountryCode string    `bson:"country_code" json:"countryCode"`
	Timezone    string    `bson:"timezone" json:"timezone"`
	ISP         string    `bson:"isp" json:"isp"`
	Lat         float64   `bson:"lat" json:"lat"`
	Lon         float64   `bson:"lon" json:"lon"`
	Org         string    `bson:"org" json:"org"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}
