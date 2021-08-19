package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/runeimp/chronus"
	// "github.com/runeimp/chronus/tzinfo"
)

const (
	appName    = "Chronus"
	appVersion = "0.1.0"
	appLabel   = "Chronus v0.1.0"
)

const usage = `%s

Usage: %s [OPTIONS] [DATE_TIME]

OPTIONS:
`

var (
	countryCodePtr *string
	debugPtr       *bool
	helpPtr        *bool
	inputPtr       *bool
	iso8601Ptr     *bool
	labelPtr       *bool
	listPtr        *bool
	pythonPtr      *bool
	rfc3339Ptr     *bool
	sqlDateTimePtr *bool
	sqlPtr         *bool
	unixAllPtr     *bool
	unixFloatPtr   *bool
	versionPtr     *bool
	internetPtr    *bool
)

func main() {
	countryCodePtr = flag.String("country-code", "", "What country code should be used in calculations")
	debugPtr = flag.Bool("debug", false, "Display debugging info")
	helpPtr = flag.Bool("help", false, "Display this help info")
	inputPtr = flag.Bool("input", false, "Display the input referenced")
	iso8601Ptr = flag.Bool("iso8601", false, "Display time in ISO 8601 formats")
	labelPtr = flag.Bool("label", false, "Display label for single formats")
	listPtr = flag.Bool("list", false, "List all supported formats")
	pythonPtr = flag.Bool("python", false, "Display a Python timestamp")
	rfc3339Ptr = flag.Bool("rfc3339", false, "Display time in RFC 3339 formats")
	sqlPtr = flag.Bool("sql", false, "Display SQL Date Time Formats")
	sqlDateTimePtr = flag.Bool("sql-datetime", false, "Display a SQL DateTime")
	unixAllPtr = flag.Bool("unix-all", false, "Display time in UNIX formats")
	unixFloatPtr = flag.Bool("unix-float", false, "Display a UNIX floating point timestamp")
	versionPtr = flag.Bool("version", false, "Display version info")
	internetPtr = flag.Bool("web", false, "Display Internet Date Time Formats")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, appLabel, filepath.Base(os.Args[0]))

		flag.VisitAll(func(f *flag.Flag) {
			optionName := fmt.Sprintf("-%s", f.Name)
			if f.DefValue == "" {
				fmt.Fprintf(flag.CommandLine.Output(), "  %-13s  %s (no default)\n", optionName, f.Usage)
			} else {
				fmt.Fprintf(flag.CommandLine.Output(), "  %-13s  %s (default: %v)\n", optionName, f.Usage, f.DefValue)
			}
		})
		fmt.Println()
	}

	flag.Parse()

	if *helpPtr {
		usageAndExit(0)
	}

	if *listPtr {
		chronus.ListFormats()
		os.Exit(0)
	}

	if *versionPtr {
		fmt.Println(appLabel)
		os.Exit(0)
	}

	if *debugPtr {
		chronus.Debug()
	}

	if len(flag.Args()) == 0 {
		// usageAndExit(0)
		t := time.Now()
		input := t.Format(time.RFC3339Nano)
		outputFormatBlocks(input)
		// printFormatInt64WithLabel("UNIX Timestamp", t.Unix())
		// printFormatFloat64WithLabel("Python Timestamp", chronus.PythonTimestamp(t))
		// printFormatStringWithLabel(t, "SQL DateTime", chronus.SQLDateTime)
		// printFormatStringWithLabel(t, "ISO 8601 Alternate", chronus.ISO8601alt)
		// printFormatStringWithLabel(t, "ISO 8601 FileSafe 1 w/Seconds", chronus.ISO8601file1Seconds)
		// printFormatStringWithLabel(t, "ISO 8601 FileSafe 2 w/Seconds", chronus.ISO8601file2Seconds)
		// fmt.Println()
		os.Exit(0)
	}

	// example := "Thu, 6 May 2021 12:46:12 -0700 (PDT)"
	// example = "Thu, 6 May 2021 12:46:12 -0700"
	// example = "Thu, 6 May 2021 12:46:12 PDT"

	for _, input := range flag.Args() {
		outputFormatBlocks(input)
	}
}

// stdError sends a formatted string to stderr
func stdError(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, f, args...)
}

func outputFormatBlocks(input string) {
	var (
		err error
		t   time.Time
	)
	if len(*countryCodePtr) > 0 {
		chronus.CountryCode = *countryCodePtr
	}
	t, err = chronus.Parse(input)
	zName, zOffset := t.Zone()
	chronus.DebugPrintf("main.outputFormatBlocks() | t.Zone().name %q | .offset %d\n", zName, zOffset)
	if err != nil {
		fmt.Printf("Time Parse Error: %s\n", err.Error())
	}
	if *inputPtr {
		fmt.Printf("                        Input: %q\n", input)
		// fmt.Printf("                       Format: %q\n", format)
	}

	if *internetPtr {
		printFormatStringWithLabel(t, "RFC 3339 DateTime", chronus.RFC3339)
	}

	if *iso8601Ptr {
		// fmt.Printf("                     ISO 8601: %s\n", t.Format(chronus.ISO8601))
		// fmt.Printf("           ISO 8601 Alternate: %s\n", t.Format(chronus.ISO8601alt))
		// fmt.Printf("           ISO 8601 FileSafe1: %s\n", t.Format(chronus.ISO8601file1))
		// fmt.Printf("           ISO 8601 FileSafe2: %s\n", t.Format(chronus.ISO8601file2))
		// fmt.Printf("ISO 8601 FileSafe 1 w/Seconds: %s\n", t.Format(chronus.ISO8601file1Seconds))
		// fmt.Printf("ISO 8601 FileSafe 2 w/Seconds: %s\n", t.Format(chronus.ISO8601file2Seconds))
		// fmt.Printf("                     RFC 3339: %s\n", t.Format(chronus.RFC3339))

		printFormatStringWithLabel(t, "ISO 8601", chronus.ISO8601)
		printFormatStringWithLabel(t, "ISO 8601 Alternate", chronus.ISO8601alt)
		printFormatStringWithLabel(t, "ISO 8601 FileSafe 1", chronus.ISO8601file1)
		printFormatStringWithLabel(t, "ISO 8601 FileSafe 2", chronus.ISO8601file2)
		printFormatStringWithLabel(t, "ISO 8601 FileSafe 1 w/Seconds", chronus.ISO8601file1Seconds)
		printFormatStringWithLabel(t, "ISO 8601 FileSafe 2 w/Seconds", chronus.ISO8601file2Seconds)
		printFormatStringWithLabel(t, "RFC 3339", chronus.RFC3339)
		// printFormatStringWithLabel(t, "____", ____)
		// printFormatStringWithLabel(t, "____", ____)
		// printFormatStringWithLabel(t, "____", ____)
		fmt.Println()
	}

	if *pythonPtr {
		if *labelPtr {
			printFormatFloat64WithLabel("Python Timestamp", chronus.PythonTimestamp(t))
		} else {
			fmt.Printf("%f\n", chronus.PythonTimestamp(t))
		}
	}

	if *rfc3339Ptr {
		if *labelPtr {
			printFormatStringWithLabel(t, "RFC 3339 DateTime", chronus.RFC3339)
		} else {
			fmt.Printf("%s\n", t.Format(chronus.RFC3339))
		}
	}

	if *sqlPtr {
		fmt.Printf("                 SQL DateTime: %s\n", t.Format(chronus.SQLDateTime))
		fmt.Printf(" SQL DateTime Year to Seconds: %s\n", t.Format(chronus.SQLDateTimeYearToSecond))
		fmt.Printf("  SQL DateTime Year to Minute: %s\n", t.Format(chronus.SQLDateTimeYearToMinute))
		fmt.Printf("         SQL Date Year to Day: %s\n", t.Format(chronus.SQLDateYearToDay))
		fmt.Printf("       SQL Date Year to Month: %s\n", t.Format(chronus.SQLDateYearToMonth))
	}

	if *sqlDateTimePtr {
		if *labelPtr {
			fmt.Printf("                 SQL DateTime: %s\n", t.Format(chronus.SQLDateTime))
		} else {
			fmt.Println(t.Format(chronus.SQLDateTime))
		}
	}

	if *unixAllPtr {
		printFormatInt64WithLabel("UNIX Timestamp", chronus.UnixTimestamp(t))
		printFormatFloat64WithLabel("UNIX Floating Point Timestamp", chronus.UnixFloat(t))
		printFormatInt64WithLabel("UNIX Timestamp in Nanoseconds", chronus.UnixNano(t))
		printFormatStringWithLabel(t, "ANSI C DateTime", chronus.ANSIC)
		printFormatStringWithLabel(t, "Git DateTime", chronus.GitDateTime)
		printFormatStringWithLabel(t, "Ruby DateTime", chronus.RubyDateTime)
		printFormatStringWithLabel(t, "UNIX DateTime", chronus.UnixDateTime)
		printFormatStringWithLabel(t, "RFC850 DateTime", chronus.RFC850)
		printFormatStringWithLabel(t, "RFC1123 DateTime", chronus.RFC1123)
		printFormatStringWithLabel(t, "RFC1123Z DateTime", chronus.RFC1123Z)
		printFormatStringWithLabel(t, "RFC3339 DateTime", chronus.RFC3339)
		// printFormatStringWithLabel(t, "____", ____)
		// printFormatStringWithLabel(t, "____", ____)
		fmt.Println()
	}

	if *unixFloatPtr {
		fmt.Printf("UNIX Floating Point Timestamp: %f\n", chronus.UnixFloat(t))
	}

	switch {
	case *iso8601Ptr:
	case *pythonPtr:
	case *rfc3339Ptr:
	case *sqlPtr:
	case *sqlDateTimePtr:
	case *unixAllPtr:
	case *unixFloatPtr:
	case *internetPtr:
	default:
		printFormatInt64WithLabel("UNIX Timestamp", t.Unix())
		printFormatFloat64WithLabel("Python Timestamp", chronus.PythonTimestamp(t))
		printFormatStringWithLabel(t, "SQL DateTime", chronus.SQLDateTime)
		printFormatStringWithLabel(t, "UK Slash Date (DD/MM/YYYY)", chronus.UKSlashDate)
		printFormatStringWithLabel(t, "US Slash Date (MM/DD/YYYY)", chronus.USSlashDate)
		printFormatStringWithLabel(t, "RFC 3339 DateTime", chronus.RFC3339)
		printFormatStringWithLabel(t, "ISO 8601", chronus.ISO8601)
		printFormatStringWithLabel(t, "ISO 8601 Alternate", chronus.ISO8601alt)
		printFormatStringWithLabel(t, "ISO 8601 FileSafe 1 w/Seconds", chronus.ISO8601file1Seconds)
		printFormatStringWithLabel(t, "ISO 8601 FileSafe 2 w/Seconds", chronus.ISO8601file2Seconds)
		fmt.Println()
	}
}

func printFormatFloat64WithLabel(label string, f float64) {
	fmt.Printf("%29s: %f\n", label, f)
}

func printFormatInt64WithLabel(label string, d int64) {
	fmt.Printf("%29s: %d\n", label, d)
}

func printFormatStringWithLabel(t time.Time, label, format string) {
	fmt.Printf("%29s: %s\n", label, t.Format(format))
	// fmt.Printf("%29s: %s | format: %q\n", label, t.Format(format), format)
}

func usageAndExit(exitCode int) {
	flag.Usage()
	os.Exit(exitCode)
}
