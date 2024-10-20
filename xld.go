package checkers

import (
	"bufio"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/agnivade/levenshtein"
)

const xld_date_format = "20060102"
const minimum_xld_version = "20100701"

var cutoff, _ = time.Parse(xld_date_format, minimum_xld_version)

func XLD(log *bufio.Scanner) error {
	if err := validVersion(log); err != nil {
		return err
	}

	fmt.Println("Valid XLD version")

	album, err := albumInfo(log)
	if err != nil {
		return err
	}

	fmt.Printf("Checking XLD log for %s - %s\n", album.Artist, album.Title)

	err = verifySettings(log)
	if err != nil {
		return err
	}

	fmt.Println("Verified rip settings")

	err = checkAccurateRip(log)
	if err != nil {
		return err
	}

	fmt.Println("Verified AccurateRip data")

	return nil
}

func validVersion(log *bufio.Scanner) error {
	var datestr string

	_, err := fmt.Sscanf(log.Text(), "X Lossless Decoder version %s", &datestr)
	if err != nil {
		return err
	}

	date, err := time.Parse(xld_date_format, datestr)
	if err != nil {
		return err
	}

	if date.Before(cutoff) {
		return errors.New("xld version is before cutoff")
	}

	return nil
}

func albumInfo(log *bufio.Scanner) (*Album, error) {
	for log.Scan() {
		if strings.Contains(log.Text(), " / ") {
			break
		}
	}

	album := Album{}
	album.Artist, album.Title, _ = strings.Cut(log.Text(), " / ")

	if album.Artist == "" || album.Title == "" {
		return nil, errors.New("invalid album information")
	}

	return &album, nil
}

func verifySettings(log *bufio.Scanner) error {
	for log.Scan() {
		if strings.HasPrefix(log.Text(), "Used") {
			break
		}
	}

	_, driveName, _ := strings.Cut(log.Text(), ": ")

	if strings.Contains(driveName, "null") {
		return errors.New("null drive information")
	}

	if slices.Contains(virtual_drives, driveName) {
		return errors.New("virtual drive detected")
	}

	parsed := ParseDrive([]any{driveName, 0.0, "", 0.0})
	lowestDistance := 999
	var drive Drive

	for _, d := range Drives {
		distance := levenshtein.ComputeDistance(d.Identifier, parsed.Identifier)

		if distance < lowestDistance {
			lowestDistance = distance
			drive = d
		}

		if distance == 0 {
			break
		}
	}

	if drive.Score < 100 {
		return errors.New("drive score is too low")
	}

	log.Scan()

	if !strings.Contains(log.Text(), "Pressed CD") {
		return errors.New("media type must be pressed cd")
	}

	log.Scan()
	log.Scan()

	if !strings.Contains(log.Text(), "XLD Secure Ripper") {
		return errors.New("ripper mode must be secure")
	}

	log.Scan()

	if !strings.Contains(log.Text(), "OK") && !strings.Contains(log.Text(), "YES") {
		return errors.New("disable audio cache must be ok or yes")
	}

	log.Scan()

	if !strings.Contains(log.Text(), "NO") {
		return errors.New("must not make use of c2 pointers")
	}

	log.Scan()

	_, offset_s, _ := strings.Cut(log.Text(), ": ")

	offset, err := strconv.Atoi(offset_s)
	if err != nil {
		return errors.New("invalid offset value")
	}

	if offset == 0 {
		return errors.New("read offset is almost never zero")
	}

	if offset != drive.Offset {
		fmt.Println(drive.Identifier)
		fmt.Println(drive.Offset)
		return errors.New("read offset does not match drive data")
	}

	log.Scan()

	_, maxRetries, _ := strings.Cut(log.Text(), ": ")
	if maxRetries == "" {
		return errors.New("no max retry information found")
	}

	var count int
	_, err = fmt.Sscanf(maxRetries, "%d", &count)
	if err != nil || count < 10 {
		return errors.New("max retries must be at least 10")
	}

	return nil
}

func checkAccurateRip(log *bufio.Scanner) error {
	for log.Scan() {
		if strings.HasPrefix(strings.TrimSpace(log.Text()), "-") {
			break
		}
	}

	if strings.Contains(log.Text(), "not") {
		return errors.New("at least one track was not accurately ripped")
	}

	return nil
}

func checkAllTracks(log *bufio.Scanner) error {
	for log.Scan() {
		if strings.Contains(log.Text(), "Statistics") {
			break
		}
	}

	log.Scan()

	_, readErrors, _ := strings.Cut(log.Text(), ": ")

	count, err := strconv.Atoi(readErrors)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("rip had at least one read error")
	}

	for log.Scan() {
		if strings.Contains(log.Text(), "Damaged") {
			break
		}
	}

	_, damaged, _ := strings.Cut(log.Text(), ": ")

	count, err = strconv.Atoi(damaged)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("rip had at least one damaged sector")
	}

	return nil
}
