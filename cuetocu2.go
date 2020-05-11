package cuetocu2

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ErrorMultiBin attempt to be generated a cu2 for a cue file with multiplebin
var ErrorMultiBin = errors.New("Cu2 can not be created for multibin cue files, merge your bin files first")

type track struct {
	ID        int
	TrackType string
	Indexes   []index
}

type index struct {
	ID    int
	Stamp string
}

// Generate creates a cu2 file in the given destination using the given cue file as base
func Generate(cuePath string, destination string) error {
	bin, err := getCueBinPath(cuePath)
	binPath := filepath.Join(filepath.Dir(cuePath), bin)

	if err != nil {
		return err
	}

	fi, err := os.Stat(binPath)
	if err != nil {
		return err
	}

	cueMap, err := cueToCueMap(cuePath)
	if err != nil {
		return err
	}

	blockSize, err := getBlockSize(cueMap[0].TrackType)
	size := sectorsToStamp(fi.Size() / int64(blockSize))

	cu2, err := cueMapToCu2(cueMap, size)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_ = os.MkdirAll(destination, os.ModePerm)

	cueName := filepath.Base(cuePath)[0:len(filepath.Base(cuePath))-len(filepath.Ext(filepath.Base(cuePath)))] + ".cu2"
	cu2Path := filepath.Join(destination, cueName)
	f, err := os.Create(cu2Path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(cu2)

	return nil
}

func cueMapToCu2(cueMap []track, size string) (string, error) {
	cu2 := fmt.Sprintf("ntracks %d\n", len(cueMap))
	cu2 += fmt.Sprintf("size      %s\n", size)
	cu2 += "data1     00:02:00\n"

	for i, track := range cueMap {
		if i == 0 {
			continue
		}

		if len(track.Indexes) == 1 {
			sectors, err := stampToSectors(track.Indexes[0].Stamp)
			if err != nil {
				return "", err
			}

			stamp := sectorsToStamp(int64(sectors + (2 * sectorsPerSecond)))
			cu2 += fmt.Sprintf("pregap%02d  %s\n", track.ID, stamp)
			cu2 += fmt.Sprintf("track%02d   %s\n", track.ID, stamp)

			continue
		}

		cu2 += fmt.Sprintf("pregap%02d  %s\n", track.ID, track.Indexes[1].Stamp)
		sectors, err := stampToSectors(track.Indexes[1].Stamp)
		if err != nil {
			return "", err
		}

		cu2 += fmt.Sprintf("track%02d   %s\n", track.ID, sectorsToStamp(int64(sectors+(2*sectorsPerSecond))))
	}

	sectors, err := stampToSectors(size)
	if err != nil {
		return "", err
	}

	cu2 += fmt.Sprintf("\ntrk end   %s", sectorsToStamp(int64(sectors+(2*sectorsPerSecond))))

	return cu2, nil
}

func getCueBinPath(cuePath string) (string, error) {
	f, err := os.Open(cuePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		switch fields[0] {
		case "FILE":
			return strings.Replace(strings.Join(fields[1:len(fields)-1], " "), "\"", "", -1), nil
		}
	}

	return "", errors.New("No file reference was found in the given cue sheet")
}

func cueToCueMap(cuePath string) ([]track, error) {
	f, err := os.Open(cuePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tracks []track
	scanner := bufio.NewScanner(f)
	files := 0

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		switch fields[0] {
		case "FILE":
			files++
			if files > 1 {
				return nil, ErrorMultiBin
			}
		case "TRACK":
			var track track
			track.ID, err = strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			track.TrackType = fields[2]
			tracks = append(tracks, track)
		case "INDEX":
			var index index
			index.ID, err = strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}

			index.Stamp = fields[2]
			if err != nil {
				return nil, err
			}

			lastTrack := &tracks[len(tracks)-1]
			lastTrack.Indexes = append(
				lastTrack.Indexes,
				index,
			)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}
