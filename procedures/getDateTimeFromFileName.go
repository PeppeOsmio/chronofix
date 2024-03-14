package procedures

import (
	"errors"
	"metadata_restorer/structs"
	"regexp"
	"strconv"
	"time"
)

func GetDateTimeFromFileName(fileName string, parsers []*structs.MediaDateTimeParser) (matchedIndex int, timeToSet time.Time, err error) {
	year := 0
	month := 0
	day := 0
	hours := 0
	minutes := 0
	seconds := 0
	matchedIndex = -1

	for i, parser := range parsers {
		regex := regexp.MustCompile(parser.Regex)
		if !regex.MatchString(fileName) {
			continue
		}
		matches := regex.FindStringSubmatch(fileName)
		if len(matches) == 0 {
			continue
		}

		matchedIndex = i
		if parser.YearIndex > 0 {
			parsedYear, err := strconv.ParseInt(matches[parser.YearIndex], 10, 64)
			if err != nil {
				return matchedIndex, timeToSet, errors.New("Error parsing year: " + err.Error())
			}
			year = int(parsedYear)
		}

		if parser.MonthIndex > 0 {
			parsedMonth, err := strconv.ParseInt(matches[parser.MonthIndex], 10, 64)
			if err != nil {
				return matchedIndex, timeToSet, errors.New("Error parsing month: " + err.Error())
			}
			month = int(parsedMonth)
		}
		if parser.DayIndex > 0 {
			parsedDay, err := strconv.ParseInt(matches[parser.DayIndex], 10, 64)
			if err != nil {
				return matchedIndex, timeToSet, errors.New("Error parsing day: " + err.Error())
			}
			day = int(parsedDay)
		}
		if parser.HoursIndex > 0 {
			parsedHours, err := strconv.ParseInt(matches[parser.HoursIndex], 10, 64)
			if err != nil {
				return matchedIndex, timeToSet, errors.New("Error parsing hours: " + err.Error())
			}
			hours = int(parsedHours)
		}
		if parser.MinutesIndex > 0 {
			parsedMinutes, err := strconv.ParseInt(matches[parser.MinutesIndex], 10, 64)
			if err != nil {
				return matchedIndex, timeToSet, errors.New("Error parsing minutes: " + err.Error())
			}
			minutes = int(parsedMinutes)
		}
		if parser.SecondsIndex > 0 {
			parsedSeconds, err := strconv.ParseInt(matches[parser.SecondsIndex], 10, 64)
			if err != nil {
				return matchedIndex, timeToSet, errors.New("Error parsing seconds: " + err.Error())
			}
			seconds = int(parsedSeconds)
		}
		break
	}
	if matchedIndex == -1 {
		return -1, timeToSet, nil
	}
	return matchedIndex, time.Date(year, time.Month(month), day, hours, minutes, seconds, 0, time.Local), nil
}
