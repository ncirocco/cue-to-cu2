package cuetocu2

import "fmt"

var blockSizes = map[string]int{
	"CDG":        2448,
	"MODE2/2336": 2336,
	"CDI/2336":   2336,
	"AUDIO":      2352,
	"MODE1/2352": 2352,
	"MODE2/2352": 2352,
	"CDI/2352":   2352,
	"MODE1/2048": 2048,
}

func getBlockSize(trackType string) (int, error) {
	if blockSize, ok := blockSizes[trackType]; ok {
		return blockSize, nil
	}
	return 0, fmt.Errorf("No blockSize defined for %s", trackType)
}
