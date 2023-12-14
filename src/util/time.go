package util

import "time"

func GetCSTTime() time.Time {
	return time.Now().UTC().Add(time.Hour * 8)
}

func GetCSTTimeLocation() *time.Location {
	return time.FixedZone("CST", 8*3600)
}
