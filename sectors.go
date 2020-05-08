package cuetocu2

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const sectorsPerSecond int = 75
const secondsPerMinute int = 60

func stampToSectors(stamp string) (int, error) {
	s := strings.Split(stamp, ":")
	if len(s) != 3 {
		return 0, errors.New("The given stamp is invalid")
	}

	minutes, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(s[1])
	if err != nil {
		return 0, err
	}

	fields, err := strconv.Atoi(s[2])
	if err != nil {
		return 0, err
	}

	sectors := minutes*secondsPerMinute*sectorsPerSecond + seconds*sectorsPerSecond + fields

	return sectors, nil
}

func sectorsToStamp(sectors int64) string {
	minutes := sectors / int64((sectorsPerSecond * secondsPerMinute))
	fields := sectors % int64((sectorsPerSecond * secondsPerMinute))
	seconds := fields / int64(sectorsPerSecond)
	fields = sectors % int64(sectorsPerSecond)

	return fmt.Sprintf("%02d:%02d:%02d", minutes, seconds, fields)
}
