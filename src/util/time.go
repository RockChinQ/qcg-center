package util

import "time"

func GetCSTTime() time.Time {
	utcTime := time.Now().UTC()

	cstTime := time.Date(utcTime.Year(), utcTime.Month(), utcTime.Day(), utcTime.Hour(), utcTime.Minute(), utcTime.Second(), utcTime.Nanosecond(), GetCSTTimeLocation())

	return cstTime
}

func GetCSTTimeLocation() *time.Location {
	return time.FixedZone("CST", 8*3600)
}

func PrintTime(t time.Time) {
	println(t.Format("2006-01-02 15:04:05"))
	// 打印时区
	println(t.Location().String())
}
