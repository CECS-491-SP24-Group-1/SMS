package timex

import (
	"strconv"
	"time"
	_ "unsafe" //For the fmtInt and fmtFloat calls
)

/*
Extension of Golang's `time.Duration` that adds new duration
types (days, weeks, months, and years) and proper support for
them in the string functions. This type aims to be mostly
compatible with Golang's std version, so it can be dropped in
place in most cases. As a side effect of this, the underlying
type is an `int64` and has a maximum value of about 290 years.
For cases where this type may need to interact with other types
in the `time` package, it can simply be cast to a plain `Duration`
type via the following snippet:

	time.Duration(<your duration variable here>)

This is a safe operation since `time.Duration` and `DurationX`
are both `int64` types under the hood.
*/
type DurationX int64

/*
Common duration types. Keep in mind that anything past hour are
based on the most common durations for these values (eg: months
being 30 days and years being 365 days).

To count the number of units in a Duration, divide:

	second := time.Second
	fmt.Print(int64(second/time.Millisecond)) // prints 1000

To convert an integer number of units to a Duration, multiply:

	seconds := 10
	fmt.Print(time.Duration(seconds)*time.Second) // prints 10s
*/
const (
	Nanosecond  DurationX = 1
	Microsecond           = 1000 * Nanosecond
	Millisecond           = 1000 * Microsecond
	Second                = 1000 * Millisecond
	Minute                = 60 * Second
	Hour                  = 60 * Minute
	Day                   = 24 * Hour
	Week                  = 7 * Day
	Month                 = 30 * Day
	Year                  = 365 * Day
)

// Returns the name of a duration constant.
func (d DurationX) NameFor() string {
	unitMap := map[DurationX]string{
		Nanosecond:  "nanosecond",
		Microsecond: "microsecond",
		Millisecond: "millisecond",
		Second:      "second",
		Minute:      "minute",
		Hour:        "hour",
		Day:         "day",
		Week:        "week",
		Month:       "month",
		Year:        "year",
	}
	tus, ok := unitMap[d]
	if !ok {
		return ""
	}
	return tus
}

// Returns the abbreviation of a duration constant.
func (d DurationX) AbbrFor() string {
	return abbrFor(int64(d))
}

/*
String returns a string representing the duration in the form "72h3m0.5s".
Leading zero units are omitted. As a special case, durations less than one
second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
that the leading digit is non-zero. The zero duration formats as 0s. Includes
all described units for this type except weeks.
*/
func (d DurationX) String() string {
	return stringHelper(int64(d), false)
}

// Identical to `DurationX.String()` but adds a space between the duration components.
func (d DurationX) StringSp() string {
	return stringHelper(int64(d), true)
}

// Nanoseconds returns the duration as an integer nanosecond count.
func (d DurationX) Nanoseconds() int64 { return int64(d) }

// Microseconds returns the duration as an integer microsecond count.
func (d DurationX) Microseconds() int64 { return int64(d) / 1e3 }

// Milliseconds returns the duration as an integer millisecond count.
func (d DurationX) Milliseconds() int64 { return int64(d) / 1e6 }

// Seconds returns the duration as a floating point number of seconds.
func (d DurationX) Seconds() float64 { return time.Duration(d).Seconds() }

// Minutes returns the duration as a floating point number of minutes.
func (d DurationX) Minutes() float64 { return time.Duration(d).Minutes() }

// Hours returns the duration as a floating point number of hours.
func (d DurationX) Hours() float64 { return time.Duration(d).Hours() }

// Days returns the duration as a floating point number of days.
func (d DurationX) Days() float64 {
	day := d / Day
	nsec := d % Day
	return float64(day) + float64(nsec)/(24*60*60*1e9) //days * hours * minutes * seconds
}

// Weeks returns the duration as a floating point number of weeks.
func (d DurationX) Weeks() float64 {
	week := d / Week
	nsec := d % Week
	return float64(week) + float64(nsec)/(7*24*60*60*1e9) //weeks * days * hours * minutes * seconds
}

// Months returns the duration as a floating point number of months.
func (d DurationX) Months() float64 {
	month := d / Month
	nsec := d % Month
	return float64(month) + float64(nsec)/(30*24*60*60*1e9) //months * days * hours * minutes * seconds
}

// Years returns the duration as a floating point number of years.
func (d DurationX) Years() float64 {
	year := d / Year
	nsec := d % Year
	return float64(year) + float64(nsec)/(365*24*60*60*1e9) //years * days * hours * minutes * seconds
}

/*
Truncate returns the result of rounding d toward zero to a multiple of m.
If m <= 0, Truncate returns d unchanged.
*/
func (d DurationX) Truncate(m DurationX) DurationX {
	return DurationX(time.Duration(d).Truncate(time.Duration(m)))
}

/*
Round returns the result of rounding d to the nearest multiple of m.
The rounding behavior for halfway values is to round away from zero.
If the result exceeds the maximum (or minimum)
value that can be stored in a Duration,
Round returns the maximum (or minimum) duration.
If m <= 0, Round returns d unchanged.
*/
func (d DurationX) Round(m DurationX) DurationX {
	return DurationX(time.Duration(d).Round(time.Duration(m)))
}

/*
Abs returns the absolute value of d.
As a special case, math.MinInt64 is converted to math.MaxInt64.
*/
func (d DurationX) Abs() DurationX {
	return DurationX(time.Duration(d).Abs())
}

// Returns the `time.Duration` equivalent object for this object.
func (d DurationX) ToDur() time.Duration {
	return time.Duration(d)
}

// Stubbed functions from `time.go`
//
//go:linkname time_fmtInt time.fmtInt
func time_fmtInt(buf []byte, v uint64) int

//go:linkname time_fmtFrac time.fmtFrac
func time_fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64)

// Private utilities
// Returns the abbreviation of a duration constant.
func abbrFor(d int64) string {
	abbrMap := map[DurationX]string{
		Nanosecond:  "ns",
		Microsecond: "\u03BCs", //Small mu
		Millisecond: "ms",
		Second:      "s",
		Minute:      "m",
		Hour:        "h",
		Day:         "d",
		Week:        "w",
		Month:       "mo",
		Year:        "y",
	}
	abbrs, ok := abbrMap[DurationX(d)]
	if !ok {
		return ""
	}
	return abbrs
}

// Helper for `DurationX.String()`.
func stringHelper(duration int64, sp bool) string {
	//Allocate enough space for the output string; more efficient than manipulating a string
	buf := make([]byte, 0, 120)

	//Check for negative duration
	if duration < 0 {
		buf = append(buf, '-')
		duration = -duration
	}

	//Setup a unit map
	units := []int64{
		int64(Year),
		int64(Month),
		int64(Day),
		int64(Hour),
		int64(Minute),
		int64(Second),
		int64(Millisecond),
		int64(Microsecond),
		int64(Nanosecond),
	}

	//Tack on the appropriate units using the `abbrFor` function
	for _, unit := range units {
		if duration >= unit {
			val := duration / unit
			duration %= unit
			buf = append(buf, strconv.FormatInt(val, 10)...)
			buf = append(buf, abbrFor(unit)...)
			if sp {
				buf = append(buf, ' ')
			}
		}
	}

	//Convert the buffer to a string and return it
	return string(buf)
}
