package godate

import (
	"math"
	"strconv"
	"strings"
	"time"
)

type GoDate struct {
	Time     time.Time
	TimeZone *time.Location
}

//IsBefore checks if the GoDate is before the passed GoDate
func (d *GoDate) IsBefore(compare *GoDate) bool {
	return d.Time.Before(compare.Time)
}

//IsAfter checks if the GoDate is before the passed GoDate
func (d *GoDate) IsAfter(compare *GoDate) bool {
	return d.Time.After(compare.Time)
}

//Sub subtracts the 'count' from the GoDate using the unit passed
func (d GoDate) Sub(count int, unit int) *GoDate {
	return d.Add(-count, unit)
}

//Add adds the 'count' from the GoDate using the unit passed
func (d GoDate) Add(count int, unit int) *GoDate {
	//milliSecondOffset := time.Millisecond * time.Duration(count/int(math.Abs(float64(count))))
	switch unit {
	case MINUTES:
		duration := time.Minute
		d.Time = d.Time.Add(duration * time.Duration(count))
	case HOURS:
		duration := time.Hour
		d.Time = d.Time.Add(duration * time.Duration(count))
	case DAYS:
		d.Time = d.Time.AddDate(0, 0, count)
	case WEEKS:
		d.Time = d.Time.AddDate(0, 0, 7*count)
	case MONTHS:
		d.Time = d.Time.AddDate(0, count, 0)
	case YEARS:
		d.Time = d.Time.AddDate(count, 0, 0)
	}
	return &d
}

//Get the difference as a duration type
func (d *GoDate) DifferenceAsDuration(compare *GoDate) time.Duration {
	return d.Time.Sub(compare.Time)
}

//Difference Returns the difference between the Godate and another in the specified unit
//If the difference is negative then the 'compare' date occurs after the date
//Else it occurs before the the date
func (d GoDate) Difference(compare *GoDate, unit int) int {
	difference := d.DifferenceAsFloat(compare, unit)
	return int(difference)
}

//Get the difference as a float
func (d GoDate) DifferenceAsFloat(compare *GoDate, unit int) float64 {
	duration := d.DifferenceAsDuration(compare)
	switch unit {
	case MINUTES:
		return duration.Minutes()
	case HOURS:
		return duration.Hours()
	case DAYS:
		return float64(duration / DAY)
	case WEEKS:
		return float64(duration / WEEK)
	case MONTHS:
		return float64(duration / MONTH)
	default:
		return float64(duration.Hours() / 24)
	}
}

//Gets the difference between the relative to the date value in the form of
//1 month before
//1 month after
func (d GoDate) DifferenceForHumans(compare *GoDate) string {
	differenceString, differenceInt := d.AbsDifferenceForHumans(compare)
	if differenceInt > 0 {
		return differenceString + " before"
	} else {
		return differenceString + " after"
	}
}

//Gets the difference between the relative to current time value in the form of
//1 month ago
//1 month from now
func (d GoDate) DifferenceFromNowForHumans() string {
	now := Now(d.TimeZone)
	differenceString, differenceInt := now.AbsDifferenceForHumans(&d)
	if differenceInt > 0 {
		return differenceString + " ago"
	} else {
		return differenceString + " from now"
	}
}

//Get the abs difference relative to compare time in the form
//1 month
//2 days
func (d GoDate) AbsDifferenceForHumans(compare *GoDate) (string, int) {
	sentence := make([]string, 2, 2)
	duration := time.Duration(math.Abs(float64(d.DifferenceAsDuration(compare))))
	unit := 0
	if duration >= YEAR {
		unit = YEARS
	} else if duration < YEAR && duration >= MONTH {
		unit = MONTHS
	} else if duration < MONTH && duration >= WEEK {
		unit = WEEKS
	} else if duration < WEEK && duration >= DAY {
		unit = DAYS
	} else if duration < DAY && duration >= time.Hour {
		unit = HOURS
	} else if duration < time.Hour && duration >= time.Minute {
		unit = MINUTES
	} else {
		unit = SECONDS
	}
	difference := d.Difference(compare, unit)
	sentence[0] = strconv.Itoa(int(math.Abs(float64(difference))))
	if difference == 1 || difference == -1 {
		sentence[1] = strings.TrimSuffix(UnitStrings[unit], "s")
	} else {
		sentence[1] = UnitStrings[unit]
	}
	return strings.Join(sentence, " "), difference
}

func (d *GoDate) StartOfHour() *GoDate {
	y, m, day := d.Time.Date()
	return &GoDate{time.Date(y, m, day, d.Time.Hour(), 0, 0, 0, d.TimeZone), d.TimeZone}
}

func (d *GoDate) StartOfDay() *GoDate {
	y, m, day := d.Time.Date()
	return &GoDate{time.Date(y, m, day, 0, 0, 0, 0, d.TimeZone), d.TimeZone}
}

func (d *GoDate) StartOfWeek() *GoDate {
	day := d.StartOfDay().Time.Weekday()
	if day != FirstDayOfWeek {
		return d.Sub(int(day-FirstDayOfWeek), DAYS).StartOfDay()
	} else{
		return d.StartOfDay()
	}
}

func (d *GoDate) StartOfMonth() *GoDate {
	y, m, _ := d.Time.Date()
	return &GoDate{time.Date(y, m, 1, 0, 0, 0, 0, d.TimeZone), d.TimeZone}
}

func (d *GoDate) StartOfQuarter() *GoDate {
	startMonth := d.StartOfMonth()
	off := (startMonth.Time.Month() - 1) % 3
	return startMonth.Sub(int(off), MONTHS)
}

func (d *GoDate) StartOfYear() *GoDate {
	y, _, _ := d.Time.Date()
	return &GoDate{time.Date(y, 1, 1, 0, 0, 0, 0, d.TimeZone), d.TimeZone}
}

func (d *GoDate) EndOfHour() *GoDate {
	nextHour := d.StartOfHour().Add(1, HOURS)
	return &GoDate{nextHour.Time.Add(-time.Millisecond), d.TimeZone}
}

func (d *GoDate) EndOfDay() *GoDate {
	nextDay := d.StartOfDay().Add(1, DAYS)
	return &GoDate{nextDay.Time.Add(-time.Millisecond), d.TimeZone}
}

func (d *GoDate) EndOfWeek() *GoDate {
	nextWeek := d.StartOfWeek().Add(1, WEEKS)
	return &GoDate{nextWeek.Time.Add(-time.Millisecond), d.TimeZone}
}

func (d *GoDate) EndOfMonth() *GoDate {
	nextWeek := d.StartOfMonth().Add(1, MONTHS)
	return &GoDate{nextWeek.Time.Add(-time.Millisecond), d.TimeZone}
}

func (d *GoDate) EndOfQuarter() *GoDate {
	nextWeek := d.StartOfQuarter().Add(3, MONTHS)
	return &GoDate{nextWeek.Time.Add(-time.Millisecond), d.TimeZone}
}

func (d *GoDate) EndOfYear() *GoDate {
	nextWeek := d.StartOfYear().Add(1, MONTHS)
	return &GoDate{nextWeek.Time.Add(-time.Millisecond), d.TimeZone}
}

//Check if this is the weekend
func (d *GoDate) IsWeekend() bool {
	day := d.Time.Weekday()
	return day == time.Saturday || day == time.Sunday
}

func (d *GoDate) Format(format string) string{
	return d.Time.Format(format)
}

func (d GoDate) String() string{
	return d.Format("Mon Jan 2 15:04:05 -0700 MST 2006")
}
