package misc 

import (
	"github.com/sdming/gosnow"
)

var snow *gosnow.SnowFlake

func StartUUId() error {
	s, err := gosnow.Default()
	if err!=nil {
		return err
	}
	snow = s
	return nil
}

func UUId() uint64 {
	id, err := snow.Next()
	if err!=nil {
		return 0
	}
	return id
}