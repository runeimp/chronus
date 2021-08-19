package chronus

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/runeimp/chronus/tzinfo"
)

const (
	Version = "0.1.0"
)

const (
	// GitDateTime specifies the format used by Git
	GitDateTime = "Mon Jan 2 15:04:05 2006 -0700" // Tue Apr 13 17:58:42 2021 -0700

	// ISO8601 specifies the format for ISO 8601 date-time
	ISO8601 = "2006-01-02T15:04:05"

	// ISO8601Z specifies the format for ISO 8601 date-time with timezone
	ISO8601Z = "2006-01-02T15:04:05Z0700"

	// ISO8601alt specifies the format for ISO 8601 date-time with timezone
	ISO8601alt = "2006-01-02 15:04:05 Z0700 (MST)"

	// ISO8601file1 specifies the format for ISO 8601 date-time with timezone for files
	ISO8601file1 = "2006-01-02_1504Z0700"

	// ISO8601file1Seconds specifies the format for ISO 8601 date-time with timezone for files including seconds
	ISO8601file1Seconds = "2006-01-02_150405Z0700"

	// ISO8601file2 specifies the format for ISO 8601 date-time with timezone for files
	ISO8601file2 = "2006-01-02_1504_MST"

	// ISO8601file2Seconds specifies the format for ISO 8601 date-time with timezone for files including seconds
	ISO8601file2Seconds = "2006-01-02_150405_MST"

	// SQLDateTime represents the default ANSI SQL DATETIME (YEAR TO SECOND) standard
	SQLDateTime = "2006-01-02 15:04:05" // "YYYY-MM-DD HH:MM:SS"

	// SQLDateTimeWithTZ represents the an SQL DATETIME (YEAR TO SECOND) with timezone
	SQLDateTimeWithTZ = "2006-01-02 15:04:05 MST" // "YYYY-MM-DD HH:MM:SS TZ"

	// SQLDateYearToMonth represents the ANSI SQL DATETIME (YEAR TO MONTH) standard
	SQLDateYearToMonth = "2006-01" // "YYYY-MM

	// SQLDateYearToDay represents the ANSI SQL DATETIME (YEAR TO Day) standard
	SQLDateYearToDay = "2006-01-02" // "YYYY-MM-DD"

	// SQLDateYearToDayWithTZ represents the an SQL DATETIME (YEAR TO Day) with timezone
	SQLDateYearToDayWithTZ = "2006-01-02 MST" // "YYYY-MM-DD TZ"

	// SQLDateTimeYearToMinute represents the ANSI SQL DATETIME (YEAR TO MINUTE) standard
	SQLDateTimeYearToMinute = "2006-01-02 15:04" // "YYYY-MM-DD HH:MM"

	// SQLDateTimeYearToMinuteWithTZ represents the an SQL DATETIME (YEAR TO MINUTE) with timezone
	SQLDateTimeYearToMinuteWithTZ = "2006-01-02 15:04 MST" // "YYYY-MM-DD HH:MM TZ"

	// SQLDateTimeYearToMinuteWithOffset represents the an SQL DATETIME (YEAR TO MINUTE) with offset
	SQLDateTimeYearToMinuteWithOffset = "2006-01-02 15:04:05 -0700" // "YYYY-MM-DD HH:MM:SS Offset"

	// SQLDateTimeYearToMinuteWithOffset2 represents the an SQL DATETIME (YEAR TO MINUTE) with offset using colon separator
	SQLDateTimeYearToMinuteWithOffset2 = "2006-01-02 15:04:05 -07:00" // "YYYY-MM-DD HH:MM:SS Offset"

	// SQLDateTimeYearToSecond represents the ANSI SQL DATETIME (YEAR TO SECOND) standard
	SQLDateTimeYearToSecond = "2006-01-02 15:04:05" // "YYYY-MM-DD HH:MM:SS"

	// SQLDateTimeYearToSecondWithOffset represents the an SQL DATETIME (YEAR TO SECOND) with timezone offset
	SQLDateTimeYearToSecondWithOffset = "2006-01-02 15:04:05 -0700" // "YYYY-MM-DD HH:MM:SS TZ"

	// SQLDateTimeYearToSecondWithOffset2 represents the an SQL DATETIME (YEAR TO SECOND) with timezone offset using colon separator
	SQLDateTimeYearToSecondWithOffset2 = "2006-01-02 15:04:05 -07:00" // "YYYY-MM-DD HH:MM:SS Offset"

	// SQLDateTimeYearToSecondWithTZ represents the an SQL DATETIME (YEAR TO SECOND) with timezone
	SQLDateTimeYearToSecondWithTZ = "2006-01-02 15:04:05 MST" // "YYYY-MM-DD HH:MM:SS TZ"

	TimezoneUnitedStatesPDT = "-0700"
	TimezoneUnitedStatesPST = "-0800"

	// UKCommon is a horrible and common UK date-time format
	UKCommon = "2 Jan 2006 15:04"

	// UKCommonWithSeconds is a horrible and common US date-time format with seconds included
	UKCommonWithSeconds = "2 Jan 2006 15:04:05"

	// UKSlashDate UK common slash date format
	UKSlashDate = "02/01/2006" // "DD/MM/YYYY"

	// UnixTimeStamp is to denote the format is a reference to a UNIX timestamp
	UnixTimeStamp = "UNIX Timestamp"

	// UnixTimeStampFloat is to denote the format is a reference to a UNIX timestamp with subsecond value
	UnixTimeStampFloat = "UNIX Timestamp with Subsecond Value"

	// USCommonDate is a horrible and common US date format
	USCommonDate = "Jan 2 2006"

	// USCommonDateTime is a horrible and common UK date-time format
	USCommonDateTime = "Jan 2 2006 15:04"

	// USCommonDateTimeWithSeconds is a horrible and common UK date-time format with seconds included
	USCommonDateTimeWithSeconds = "Jan 2 2006 15:04:05"

	// USSlashDate US common slash date format
	USSlashDate = "01/02/2006" // "MM/DD/YYYY"

	// RFC5322A defines the standard format for RFC 5322 date-time
	RFC5322A = "Mon, 2 Jan 2006 15:04:05 -0700"

	// RFC5322B defines the alternate format for RFC 5322 date-time
	RFC5322B = "Mon, 2 Jan 2006 15:04:05 MST"

	// RFC5322C defines the standard format for RFC 5322 date-time with comment in parenthesis
	RFC5322C = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)" // Zone and comment

	regExGitDateTime     = `(\w+) (\w+) (\d+) (\d\d:\d\d:\d\d) (\d+) (-\d+)`                     // Tue Apr 13 17:58:42 2021 -0700
	regExAnsiGitRubyUnix = `(\w+) (\w+) ( ?\d+) (\d\d):(\d\d)(:\d\d)? ([A-Z0-9-]+) ([A-Z0-9-]+)` // Tue Apr 13 17:58:42 2021 -0700
	regExRFC3339         = `(\d+)-(\d\d)-(\d\d)T(\d\d):(\d\d):(\d\d)(\.\d+)?([+-]\d\d:?\d\d|Z)`  // 2006-01-02T15:04:05.999999999Z07:00

	regExUnixTimeStamp = `^[+-]?(\d+)\.?(\d+)?$` // Standard 1234567890 or Python Float with microtime 1234567890.123456
	// regExUKCommon               = `\d{1,2} \w+ \d{4} \d\d:\d\d(:\d\d)?`
	// regExUKCommonStrict         = `^\d{1,2} \w+ \d{4} \d\d:\d\d(:\d\d)?$`
	regExUSCommonDateTime       = `\w+ ( \d|\d{1,2}) \d{4} \d\d:\d\d(:\d\d)?`
	regExUSCommonDateTimeStrict = `^\w+ ( \d|\d{1,2}) \d{4} \d\d:\d\d(:\d\d)?$`
	// regExRFC5322A               = `^\w+, \d{1,2} \w+ \d{4} \d\d:\d\d(:\d\d)? \-\d+$`
	// regExRFC5322B               = `^\w+, \d{1,2} \w+ \d{4} \d\d:\d\d(:\d\d)? \w+$`
	// regExRFC5322C               = `^\w+, \d{1,2} \w+ \d{4} \d\d:\d\d(:\d\d)? \-\d+( \(\w+\))?$`
	regExRFCUnitedKingdom = `^(\w+, )?( \d|\d{1,2})( \w+)( \d{4})( \d\d:\d\d)(:\d\d)?( \-\d+| \w+)?( \(\w+\))?$`

	regExSQLDateTime = `(\d+-\d+-\d+) ?(\d+:\d+(:\d+)?)? ?(\w+|[+-]?\d+:?\d+)?`
)

// Constants from Go time
const (
	// Date and Time Formats
	ANSIC        = time.ANSIC       // "Mon Jan _2 15:04:05 2006"
	UnixDateTime = time.UnixDate    // "Mon Jan _2 15:04:05 MST 2006"
	RubyDateTime = time.RubyDate    // "Mon Jan 02 15:04:05 -0700 2006"
	RFC822       = time.RFC822      // "02 Jan 06 15:04 MST"
	RFC822Z      = time.RFC822Z     // "02 Jan 06 15:04 -0700"
	RFC850       = time.RFC850      // "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123      = time.RFC1123     // "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z     = time.RFC1123Z    // "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339      = time.RFC3339     // "2006-01-02T15:04:05Z07:00"
	RFC3339Nano  = time.RFC3339Nano // "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen      = time.Kitchen     // "3:04PM"

	// Handy time stamps
	Stamp      = time.Stamp      // "Jan _2 15:04:05"
	StampMilli = time.StampMilli // "Jan _2 15:04:05.000"
	StampMicro = time.StampMicro // "Jan _2 15:04:05.000000"
	StampNano  = time.StampNano  // "Jan _2 15:04:05.000000000"
)

var (
	CountryCode         = os.Getenv("CHRONUS_COUNTRY_CODE")
	debug               = false
	reIsOffset          = regexp.MustCompile(`[+-]?\d{4}?`)
	reIsOffsetWithColon = regexp.MustCompile(`[+-]?\d{1,2}:\d{2}`)
)

func Debug() {
	debug = true
}

// DebugPrintf only prints output if chronus.debug is true
func DebugPrintf(f string, args ...interface{}) {
	if debug {
		log.Printf(f, args...)
	}
}

func init() {
	// t := time.Now()
	// unixTimeStamp := t.Unix()
	// unixTimeStampNano := t.UnixNano()
	// goTime := time.Unix(unixTimeStamp, 0)
	// var year, week int = t.ISOWeek() // ISO 8601 year and week
	// CountryCode = os.Getenv("CHRONUS_COUNTRY_CODE")
}

// GetFormat determines the correct format for the provided date-time-zone string
func GetFormat(dtz string) (format string, tzloc *tzinfo.TimeZoneLocation) {
	DebugPrintf("chronus.GetFormat() | dtz: %q\n", dtz)

	// DebugPrintf("chronus.GetFormat() | %s\n", "UNIX TimeStamp")
	// Check if it's a UNIX TimeStamp
	format = GetUnixTimeStampFormat(dtz)
	if len(format) > 0 {
		return format, tzloc
	}

	// DebugPrintf("chronus.GetFormat() | %s\n", "RFC 3339")
	// Check if it's RFC 3339
	format = GetRFC3339Format(dtz)
	if len(format) > 0 {
		DebugPrintf("chronus.GetFormat() | format: %q\n", format)
		return format, tzloc
	}

	// DebugPrintf("chronus.GetFormat() | %s\n", "SQL DateTime")
	// Check if it's an SQL DateTime style
	format, tzloc = GetSQLFormat(dtz)

	if len(format) > 0 {
		// log.Printf("chronus.GetFormat() | tzloc: %v\n", tzloc)
		// zone, offset := tzloc.Zone()
		// log.Printf("chronus.GetFormat() | zone: %q | offset: %d\n", zone, offset)
		// zulu := "Z"
		// offsetString := tzinfo.OffsetSecondsToString(offset, zulu)
		// log.Printf("chronus.GetFormat() | offsetString: %q\n", offsetString)
		return format, tzloc
	}

	// DebugPrintf("chronus.GetFormat() | %s\n", "ANSI C")
	// Check if it's an ANSI C, Git, Ruby, Unix, etc. format
	re := regexp.MustCompile(regExAnsiGitRubyUnix)
	matched := re.MatchString(dtz)

	if matched {
		re = regexp.MustCompile(regExGitDateTime)
		matched = re.MatchString(dtz)
		if matched {
			return GitDateTime, tzloc
		}
	}

	// DebugPrintf("chronus.GetFormat() | %s\n", "US Common")
	// Check if it's at least partially a US Common format
	re = regexp.MustCompile(regExUSCommonDateTime)
	matched = re.MatchString(dtz)

	if matched {
		// Check if it is the US Common format
		re = regexp.MustCompile(regExUSCommonDateTimeStrict)
		matched = re.MatchString(dtz)

		if matched == false {
			//
		}
	} else {
		i := 0
		re = regexp.MustCompile(regExRFCUnitedKingdom)
		matched = re.MatchString(dtz)
		matches := re.FindStringSubmatch(dtz)
		formatMatches := re.FindStringSubmatch(RFC5322C)

		format = ""
		for c, s := range matches {
			if len(strings.TrimSpace(s)) > 0 {
				if c > 0 {
					if formatMatches[c] == " -0700" && regexp.MustCompile("^ [A-Z]+$").MatchString(s) {
						format += " MST"
					} else {
						format += formatMatches[c]
					}
				}
				i++
			}
		}
		// fmt.Printf("chronus.GetFormat() | RFC 5322 UK matched: %-5t | matches: %d %q\n", matched, i, matches)
		// fmt.Printf("chronus.GetFormat() | RFC 5322 UK matched: %-5t | matches: %d %q\n", matched, i, formatMatches)
	}

	DebugPrintf("chronus.GetFormat() | matched: %t\n", matched)
	DebugPrintf("chronus.GetFormat() | format: %q\n", format)

	return format, tzloc
}

// GetRFC3339Format determines the correct format for an RFC 3339 based string
func GetRFC3339Format(dtz string) (format string) {
	// DebugPrintf("chronus.GetRFC3339Format() | dtz: %q\n", dtz)
	re := regexp.MustCompile(regExRFC3339)
	if re.MatchString(dtz) {
		matches := re.FindStringSubmatch(dtz)
		// DebugPrintf("chronus.GetRFC3339Format() | matches: %q\n", matches)
		// DebugPrintf("chronus.GetRFC3339Format() | submatches: %d\n", len(matches)-1)
		// DebugPrintf("chronus.GetRFC3339Format() | matches[1]: %q (year)\n", matches[1])
		// DebugPrintf("chronus.GetRFC3339Format() | matches[2]: %q (month)\n", matches[2])
		// DebugPrintf("chronus.GetRFC3339Format() | matches[3]: %q (day)\n", matches[3])
		// DebugPrintf("chronus.GetRFC3339Format() | matches[4]: %q (hour)\n", matches[4])
		// DebugPrintf("chronus.GetRFC3339Format() | matches[5]: %q (minute)\n", matches[5])
		// DebugPrintf("chronus.GetRFC3339Format() | matches[6]: %q (second)\n", matches[6])
		// DebugPrintf("chronus.GetRFC3339Format() | matches[7]: %q (micro)\n", matches[7])
		// DebugPrintf("chronus.GetRFC3339Format() | matches[8]: %q (timezone)\n", matches[8])
		if len(matches[7]) > 0 {
			// format = RFC3339Nano
			format = "2006-01-02T15:04:05.999999999" + GetTimeZoneFormat(matches[8])
		}
		// format = RFC3339
		format = "2006-01-02T15:04:05" + GetTimeZoneFormat(matches[8])
	}
	return format
}

// GetSQLFormat determines the correct format for the provided SQL based date-time-zone string
func GetSQLFormat(dtz string) (format string, tzloc *tzinfo.TimeZoneLocation) {
	DebugPrintf("chronus.GetSQLFormat() | dtz: %q\n", dtz)
	re := regexp.MustCompile(regExSQLDateTime)
	if re.MatchString(dtz) {
		matches := re.FindStringSubmatch(dtz)
		DebugPrintf("chronus.GetSQLFormat() | matches: %q\n", matches)
		DebugPrintf("chronus.GetSQLFormat() | submatches: %d\n", len(matches)-1)
		DebugPrintf("chronus.GetSQLFormat() | matches[1]: %q (date)\n", matches[1])
		DebugPrintf("chronus.GetSQLFormat() | matches[2]: %q (time)\n", matches[2])
		DebugPrintf("chronus.GetSQLFormat() | matches[3]: %q (seconds)\n", matches[3])
		DebugPrintf("chronus.GetSQLFormat() | matches[4]: %q (timezone)\n", matches[4])

		// date := matches[1]
		timeStr := matches[2]
		seconds := matches[3]
		timezone := matches[4]
		tzIsOffset := TimezoneIsOffset(timezone)

		format = SQLDateTime

		if len(timeStr) > 0 {
			format = SQLDateTimeYearToMinute
			if len(seconds) > 0 {
				format = SQLDateTimeYearToSecond
			}
		}

		if len(timezone) > 0 {
			if len(timeStr) > 0 {
				if len(seconds) > 0 {
					format = SQLDateTimeYearToSecondWithTZ
					if tzIsOffset == 1 {
						format = SQLDateTimeYearToSecondWithOffset
					} else if tzIsOffset == 2 {
						format = SQLDateTimeYearToSecondWithOffset2
					}
				} else {
					format = SQLDateTimeYearToMinuteWithTZ
					if tzIsOffset == 1 {
						format = SQLDateTimeYearToMinuteWithOffset
					} else if tzIsOffset == 2 {
						format = SQLDateTimeYearToMinuteWithOffset2
					}
				}
			}
		}

		if len(timezone) > 0 {
			if tzIsOffset == 0 {
				// CountryCode := os.Getenv("CHRONUS_COUNTRY_CODE")
				// zone, offset := time.Now().UTC().Zone()
				tzloc = tzinfo.GetCurrentTimeZoneLocation()
				zone, offset := tzloc.Zone()
				DebugPrintf("chronus.GetSQLFormat() | zone: %q | offset: %d\n", zone, offset)
				zulu := "Z"
				offsetString := tzinfo.OffsetSecondsToString(offset, zulu)
				DebugPrintf("chronus.GetSQLFormat() | offsetString: %q\n", offsetString)
				if len(CountryCode) > 0 {
					CountryCode = strings.ToUpper(CountryCode)
					var err error
					switch CountryCode {
					case "US", "USA":
						tzloc, err = tzinfo.GetUSTimeZoneLocationByTZAbbreviation(timezone)
						if err != nil {
							DebugPrintf("chronus.GetSQLFormat() | matches: %q\n", err.Error())
						}
					default:
						// NOTE: need to add more TZ Location by TZ Abbreviation methods ~RuneImp
					}
				}
			}
		}
	}
	DebugPrintf("chronus.GetSQLFormat() | format: %q\n", format)

	return format, tzloc
}

func GetTimeZoneFormat(tz string) (format string) {
	match, _ := regexp.MatchString(`[A-Z][A-Z]+`, tz)
	if match {
		return "MST"
	}

	match, _ = regexp.MatchString(`[+-]\d+:\d+`, tz)
	if match {
		return "Z07:00"
	}

	return "Z0700"
}

// GetUnixTimeStampFormat determines the correct format for the provided UNIX timestamp string
func GetUnixTimeStampFormat(dtz string) (format string) {
	re := regexp.MustCompile(regExUnixTimeStamp)
	if re.MatchString(dtz) {
		matches := re.FindStringSubmatch(dtz)
		DebugPrintf("chronus.GetUnixTimeStampFormat() | UNIX timestamp matched: true | matches: %q\n", matches)
		if len(matches) == 3 && len(matches[2]) > 0 {
			format = UnixTimeStampFloat
		}
		format = UnixTimeStamp
	}

	return format
}

// ListFormats prints a list of all supported time formats
func ListFormats() {
	fmt.Printf("   ISO8601alt: %q\n", ISO8601alt)
	fmt.Printf("     ISO8601Z: %q\n", ISO8601Z)
	fmt.Printf("      ISO8601: %q\n", ISO8601)
	fmt.Printf("  SQLDateTime: %q\n", SQLDateTime)
	fmt.Printf("UnixTimeStamp: %q\n", "1136239445")

}

// Parse attempts to convert a given string into a Go time.Time
func Parse(dtz string) (t time.Time, err error) {
	format, tzloc := GetFormat(dtz)
	DebugPrintf("chronus.Parse() | dtz: %q\n", dtz)
	DebugPrintf("chronus.Parse() | format: %q\n", format)
	DebugPrintf("chronus.Parse() | tzloc: %s\n", tzloc.String())

	// loc []*tzinfo.TimeZoneLocation

	// // if len(tzlocs) > 0 {
	// // 	tzloc := tzlocs[0]
	// // 	DebugPrintf("chronus.GetFormat() | format: %q\n", format)
	// // }

	// for i, tzloc := range tzlocs {
	// 	DebugPrintf("chronus.GetFormat() | i = %d | format: %q\n", i, tzloc)
	// }

	switch format {
	case UnixTimeStamp:
		i, err := strconv.ParseInt(dtz, 10, 64)
		if err == nil {
			t = time.Unix(i, 0)
		}
	case UnixTimeStampFloat:
		f, err := strconv.ParseFloat(dtz, 64)
		if err == nil {
			i := int64(f * 1000000000)
			t = time.Unix(0, i)
		}
	default:
		if tzloc != nil {
			t, err = time.ParseInLocation(format, dtz, tzloc.Location())
		} else {
			t, err = time.Parse(format, dtz)
		}

		if err != nil {
			fmt.Printf("Time Parse Error: %s\n", err.Error())
			DebugPrintf("input format: %q\n", format)
		}
	}

	return t, err
}

// UnixFloat converts Go time.Time into a floating point UNIX timestamp (ala Python) since the UNIX Epoch
func PythonTimestamp(t time.Time) float64 {
	return UnixFloat(t)
}

// TimezoneIsOffset returns 1 if the string is an offset or 2 if it is an offset with a colon
func TimezoneIsOffset(tz string) int {
	if reIsOffset.MatchString(tz) {
		return 1
	}
	if reIsOffsetWithColon.MatchString(tz) {
		return 2
	}
	return 0
}

// UnixFloat converts Go time.Time into a floating point UNIX timestamp (ala Python) since the UNIX Epoch
func UnixFloat(t time.Time) float64 {
	ts := fmt.Sprintf("%d", t.Unix())
	tsn := fmt.Sprintf("%d", t.UnixNano())
	fString := fmt.Sprintf("%s.%s", ts, tsn[len(ts):])
	fFloat, _ := strconv.ParseFloat(fString, 64)
	return fFloat
	// return float64(t.UnixNano()) // 1_000_000_000.0
}

// UnixMilli converts Go time.Time into milliseconds since the UNIX Epoch
func UnixMilli(t time.Time) int64 {
	nanos := t.UnixNano()
	millis := nanos / 1000000
	return millis
}

// UnixMilliToTime converts milliseconds since the UNIX Epoch into Go time.Time
func UnixMilliToTime(s int64) time.Time {
	millis := s * 1000000
	return time.Unix(0, millis)
}

// UnixNano converts Go time.Time into nanoseconds since the UNIX Epoch
func UnixNano(t time.Time) int64 {
	return t.UnixNano()
}

// UnixNanoToTime converts nanoseconds since the UNIX Epoch into Go time.Time
func UnixNanoToTime(s int64) time.Time {
	return time.Unix(0, s)
}

// UnixTimestamp converts Go time.Time into seconds since the UNIX Epoch
func UnixTimestamp(t time.Time) int64 {
	return t.Unix()
}

// UnixTimestampToTime converts seconds since the UNIX Epoch into Go time.Time
func UnixTimestampToTime(s int64) time.Time {
	return time.Unix(s, 0)
}
