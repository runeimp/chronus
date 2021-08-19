package tzinfo

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	ErrorUnknownLocationStatus = "Unknown Location Status"
)

var (
	debug                       = false
	timeZoneDataByAbbreviation  map[string][]*TimeZoneLocation
	TimeZoneLocations           []*TimeZoneLocation
	unitedStatesAbbreviationMap map[string][]*TimeZoneLocation
)

type TimeZoneLocation struct {
	countryCodeAlpha2 string
	countryCodeAlpha3 string
	ianaName          string
	location          *time.Location
	nation            string
	offset            int
	status            string // timezone location status
	zone              string
}

// CountryCode returns the Alpha-2 (or Alpha-3 when necessary) country code if known
func (tzloc *TimeZoneLocation) CountryCode() string {
	// check if the Alpha-2 is set (not zero length) and not garbage (only two characters)
	if len(tzloc.countryCodeAlpha2) == 2 {
		return tzloc.countryCodeAlpha2
	}
	return tzloc.countryCodeAlpha3
}

// IANA returns the IANA time zone name
func (tzloc *TimeZoneLocation) IANA() string {
	return tzloc.ianaName
}

// Location returns a *time.Location for the time zone
func (tzloc *TimeZoneLocation) Location() *time.Location {
	if tzloc.location == nil {
		loc, err := time.LoadLocation(tzloc.ianaName)
		if err != nil {
			log.Printf("location error: %s\n", err.Error())
		} else {
			tzloc.location = loc
		}
	}

	return tzloc.location
}

// Nation returns the nation for the time zone if known
func (tzloc *TimeZoneLocation) Nation() string {
	return tzloc.nation
}

// Offset returns a formatted offset string. The default formatting is ±HHMM.
// For UTC/Zulu time is -0000. You can include modifiers to the zulu argument
// to adjust formatting. To add a colon for ±HH:MM add a colon to the
// zulu argument. To force a + instead of a - for UTC/Zulu time add one.
// To use use Z instead of a variation of -0000 just add it.
// Example: to specify ±HH:MM and use Z for UTC/Zulu time use the string ":Z".
// The order and capitalization of characters does not matter.
func (tzloc *TimeZoneLocation) Offset(zulu string) string {
	_, offset := time.Now().In(tzloc.location).Zone()
	return OffsetSecondsToString(offset, zulu)
}

// Status returns the current location status, one of alias, canonical, deprecated, or a zero length string if unset for the TZ Location
func (tzloc *TimeZoneLocation) Status() string {
	return tzloc.status
}

// StatusUpdate sets the current location status to one of alias, canonical, or deprecated
func (tzloc *TimeZoneLocation) StatusUpdate(s string) (err error) {
	status := strings.ToLower(s)
	switch status {
	case "", "alias", "canonical", "deprecated":
		tzloc.status = s
	default:
		err = errors.New(ErrorUnknownLocationStatus)
	}

	return err
}

// String returns a JSON encoded version of the data
func (tzloc *TimeZoneLocation) String() string {
	jsonStr := "{\n"

	jsonStr += `	"countryCodeAlpha2": "` + tzloc.countryCodeAlpha2 + `"` + ",\n"
	jsonStr += `	"countryCodeAlpha3": "` + tzloc.countryCodeAlpha3 + `"` + ",\n"
	jsonStr += `	"ianaName": "` + tzloc.ianaName + `"` + ",\n"
	jsonStr += `	"location": "` + tzloc.location.String() + `"` + ",\n"
	jsonStr += `	"nation": "` + tzloc.nation + `"` + ",\n"
	jsonStr += `	"offset": ` + strconv.Itoa(tzloc.offset) + ",\n"
	jsonStr += `	"status": "` + tzloc.status + `"` + ",\n"
	jsonStr += `	"zone": "` + tzloc.zone + `"` + ",\n"

	jsonStr = jsonStr[:len(jsonStr)-2] + "\n}\n"

	// byteSlice, err := json.Marshal(tzloc)
	// if err != nil {
	// 	return err.Error()
	// }
	// return string(byteSlice)
	return jsonStr
}

func (tzloc *TimeZoneLocation) Zone() (string, int) {
	return tzloc.zone, tzloc.offset
}

func init() {
	timeZoneDataByAbbreviation = make(map[string][]*TimeZoneLocation)
	unitedStatesAbbreviationMap = make(map[string][]*TimeZoneLocation)
	setTimeZoneLocation("AKDT", "-0900", "US/Alaska", "America")
	setTimeZoneLocation("AKST", "-1000", "US/Alaska", "America")
	setTimeZoneLocation("CDT", "-0500", "US/Central", "America")
	setTimeZoneLocation("CST", "-0600", "US/Central", "America")
	setTimeZoneLocation("EDT", "-0400", "US/Eastern", "America")
	setTimeZoneLocation("EST", "-0500", "US/Eastern", "America")
	setTimeZoneLocation("HDT", "-1000", "US/Aleutian", "America")
	setTimeZoneLocation("HST", "-1100", "US/Aleutian", "America")
	setTimeZoneLocation("HST", "-1100", "US/Hawaii", "America")
	setTimeZoneLocation("MDT", "-0600", "US/Mountain", "America")
	setTimeZoneLocation("MST", "-0700", "US/Mountain", "America")
	setTimeZoneLocation("PDT", "-0700", "US/Pacific", "America")
	setTimeZoneLocation("PST", "-0800", "US/Pacific", "America")
	setTimeZoneLocation("SST", "-1200", "US/Samoa", "America")
}

// DebugPrintf only prints output if tzinfo.debug is true
func DebugPrintf(f string, args ...interface{}) {
	if debug {
		log.Printf(f, args...)
	}
}

func GetCurrentTimeZoneLocation() *TimeZoneLocation {
	t := time.Now()
	zone, offset := t.Zone()

	// zulu := "Z"
	// offsetString := OffsetSecondsToString(offset, zulu)

	ianaLoc := "" // Somewhere/OverTheRainbow
	nation := ""  // Oz
	loc := time.FixedZone(zone, offset)

	tzloc := &TimeZoneLocation{
		zone:     zone,
		ianaName: ianaLoc,
		location: loc,
		nation:   nation,
		offset:   offset,
	}

	return tzloc
}

func GetNationByTZAbbreviation(abbr string) (string, error) {
	if data, ok := timeZoneDataByAbbreviation[abbr]; ok {
		return data[0].nation, nil
	}

	err := fmt.Errorf("zone %q not found in timezone data")
	return "", err
}

func GetTimeZoneLocationsByTZAbbreviation(abbr string) (loc []*TimeZoneLocation, err error) {
	if tzlocs, ok := timeZoneDataByAbbreviation[abbr]; ok {
		return tzlocs, err
	}

	err = fmt.Errorf("zone %q not found in timezone data")
	return nil, err
}

func GetUSTimeZoneLocationByTZAbbreviation(abbr string) (loc *TimeZoneLocation, err error) {
	if tzlocs, ok := unitedStatesAbbreviationMap[abbr]; ok {
		for i, tzloc := range tzlocs {
			DebugPrintf("chronus.GetUnitedStatesLocationByAbbreviation() | i: %d | tzloc: %q\n", i, tzloc)
			return tzloc, err
			break
		}
	}

	err = fmt.Errorf("zone %q not found in United States timezone data", abbr)
	return nil, err
}

func GetUSTimeZoneLocationsByTZAbbreviation(abbr string) (loc []*TimeZoneLocation, err error) {
	if tzlocs, ok := unitedStatesAbbreviationMap[abbr]; ok {
		return tzlocs, err
	}

	err = fmt.Errorf("zone %q not found in United States timezone data")
	return nil, err
}

func OffsetStringToSeconds(offsetStr string) (offset int) {
	// log.Printf("tzinfo.OffsetStringToSeconds() | offsetStr: %q\n", offsetStr)
	offsetStr = strings.TrimSpace(offsetStr)

	sign := '+'
	h10 := 'X'
	h01 := 'X'
	m10 := 'X'
	m01 := 'X'

	for _, r := range offsetStr {
		// log.Printf("tzinfo.OffsetStringToSeconds() | i: %d | r: %c\n", i, r)
		switch r {
		case '+', '-':
			// sign
			sign = r
		case ':':
			// colon ignored
		default:
			// digit
			if h10 == 'X' {
				h10 = r
			} else if h01 == 'X' {
				h01 = r
			} else if m10 == 'X' {
				m10 = r
			} else if m01 == 'X' {
				m01 = r
			} else {
				log.Fatalln("more than four digits in offset string")
			}
		}
	}

	hoursStr := fmt.Sprintf("%c%c", h10, h01)
	hoursInt, err := strconv.Atoi(hoursStr)
	if err != nil {
		// log.Printf("tzinfo.OffsetStringToSeconds() | err: %q\n", )
		log.Fatalln(err)
	}
	minutesStr := fmt.Sprintf("%c%c", m10, m01)
	minutesInt, err := strconv.Atoi(minutesStr)
	if err != nil {
		// log.Printf("tzinfo.OffsetStringToSeconds() | err: %q\n", err.Error())
		log.Fatalln(err)
	}

	if hoursInt > 0 {
		offset = hoursInt * 60 * 60
	}
	if minutesInt > 0 {
		offset += minutesInt * 60
	}
	if sign == '-' {
		offset = -offset
	}
	// log.Printf("tzinfo.OffsetStringToSeconds() | offset: %d\n", offset)

	offsetStr = OffsetSecondsToString(offset, "Z")
	// log.Printf("tzinfo.OffsetStringToSeconds() | offsetStr: %s\n\n", offsetStr)

	return offset
}

// OffsetSecondsToString returns a formatted offset string. The default formatting is ±HHMM. For UTC/Zulu time is -0000.
// You can include modifiers to the zulu argument to adjust formatting. To add a colon for ±HH:MM add one to the zulu argument.
// To force a + instead of a - for UTC/Zulu time add one. To use use Z instead of a variation of -0000 just add it.
// Example: to specify ±HH:MM and use Z for UTC/Zulu time use the string ":Z", or ":z". The order of characters does not matter.
func OffsetSecondsToString(offset int, zulu string) (offsetStr string) {
	colonStr := ""
	hours := 3600
	hourSeconds := 0
	minutes := 60
	minuteSeconds := 0
	positiveOffset := offset

	if offset == 0 {
		if len(zulu) > 0 {
			if strings.ContainsAny(zulu, "Zz") {
				return "Z"
			}

			result := "-00"
			if strings.Contains(zulu, "+") {
				result = "+00"
			}
			if strings.Contains(zulu, ":") {
				result += ":00"
			} else {
				result += "00"
			}
			return result
		}

		return "-0000"
	}

	if offset < 0 {
		positiveOffset = -offset
	}

	hourSeconds = positiveOffset / hours
	minuteSeconds = (positiveOffset - (hourSeconds * hours)) / minutes
	offsetStr = fmt.Sprintf("%02d%s%02d", hourSeconds, colonStr, minuteSeconds)

	DebugPrintf("tzinfo.OffsetSecondsToString() | offset: %d | positiveOffset: %d\n", offset, positiveOffset)
	DebugPrintf("tzinfo.OffsetSecondsToString() | hourSeconds: %02d\n", hourSeconds)
	DebugPrintf("tzinfo.OffsetSecondsToString() | minuteSeconds: %02d\n", minuteSeconds)

	if offset > 0 {
		offsetStr = "+" + offsetStr
	} else {
		offsetStr = "-" + offsetStr
	}

	return offsetStr
}

func setTimeZoneLocation(abbr, offset, ianaLoc, nation string) {
	var tzloc *TimeZoneLocation

	loc, err := time.LoadLocation(ianaLoc)
	if err != nil {
		log.Printf("location error: %s\n", err.Error())
	}

	nationUpper := strings.ToUpper(nation)
	switch nationUpper {
	case "AMERICA", "US", "USA":
		status := "canonical"
		switch loc.String() {
		case "America/Atka",
			"America/Fort_Wayne",
			"America/Indianapolis",
			"America/Knox_IN",
			"America/Louisville",
			"America/Shiprock",
			"Navajo",
			"US/Alaska",
			"US/Aleutian",
			"US/Arizona",
			"US/Central",
			"US/East-Indiana",
			"US/Eastern",
			"US/Hawaii",
			"US/Indiana-Starke",
			"US/Michigan",
			"US/Mountain",
			"US/Pacific":
			status = "deprecated"
		}
		tzloc = &TimeZoneLocation{
			countryCodeAlpha2: "US",
			countryCodeAlpha3: "USA",
			zone:              abbr,
			ianaName:          ianaLoc,
			location:          loc,
			nation:            nation,
			offset:            OffsetStringToSeconds(offset),
			status:            status,
		}
		unitedStatesAbbreviationMap[abbr] = append(unitedStatesAbbreviationMap[abbr], tzloc)
		timeZoneDataByAbbreviation[abbr] = append(unitedStatesAbbreviationMap[abbr], tzloc)
	default:
		tzloc = &TimeZoneLocation{
			zone:     abbr,
			ianaName: ianaLoc,
			location: loc,
			nation:   nation,
			offset:   OffsetStringToSeconds(offset),
		}
	}
	TimeZoneLocations = append(TimeZoneLocations, tzloc)
}
