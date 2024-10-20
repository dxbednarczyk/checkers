package checkers

import (
	"regexp"
	"strings"

	"github.com/agnivade/levenshtein"
)

var virtualDrives = []string{
	"Generic DVD-ROM",
	"Generic DVD-ROM SCSI CdRom Device",
}

var (
	tsst       = regexp.MustCompile("TSSTcorp(BD|CD|DVD)")
	hlds       = regexp.MustCompile("HL-DT-ST(BD|CD|DVD)")
	special    = regexp.MustCompile(`^[ _-]+`)
	spacedash  = regexp.MustCompile(`\s+-\s`)
	multispace = regexp.MustCompile(`\s+`)
	revision   = regexp.MustCompile(`\(revision [a-zA-Z0-9\.\,\-]*\)`)
	adapter    = regexp.MustCompile(` Adapter.*$`)
)

type Album struct {
	Artist string
	Title  string
}

type Drive struct {
	Identifier string
	Offset     int
	Score      int
}

func ParseDrive(drive []any) Drive {
	name := drive[0].(string)

	name = strings.ReplaceAll(name, "JLMS", "Lite-ON")
	name = tsst.ReplaceAllString(name, "TSSTCorp")
	name = hlds.ReplaceAllString(name, "HL-DT-ST")
	name = strings.ReplaceAll(name, "HL-DT-ST", "LG Electronics")
	name = strings.ReplaceAll(name, "Matshita", "Panasonic")
	name = strings.ReplaceAll(name, "MATSHITA", "Panasonic")
	name = special.ReplaceAllString(name, "")
	name = spacedash.ReplaceAllString(name, " ")
	name = multispace.ReplaceAllString(name, " ")
	name = revision.ReplaceAllString(name, "")
	name = adapter.ReplaceAllString(name, "")
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	return Drive{
		Identifier: name,
		Offset:     int(drive[1].(float64)),
		Score:      int(drive[3].(float64)),
	}
}

func GetClosestDrive(driveName string) Drive {
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

	return drive
}
