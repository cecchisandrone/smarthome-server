package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const intervalsSeparator = ","
const timeSeparator = "-"
const hourMinutesSeparator = "."

func ParseTimeIntervals(intervalsString string) (map[string][]time.Time, error) {

	timeIntervals := make(map[string][]time.Time)

	intervalTokens := strings.Split(intervalsString, intervalsSeparator)
	for _, intervalToken := range intervalTokens {
		timeTokens := strings.Split(intervalToken, timeSeparator)
		if len(timeTokens) != 2 || timeTokens[0] == "" || timeTokens[1] == "" {
			return nil, errors.New("Unable to parse string " + intervalsString)
		}

		// Here we have time tokens: [9, 12.30]
		currentTime := time.Now()

		// Start time
		startTimeTokens := strings.Split(timeTokens[0], hourMinutesSeparator)
		if len(startTimeTokens) == 2 && (startTimeTokens[0] == "" || startTimeTokens[1] == "") {
			return nil, errors.New("Unable to parse string " + timeTokens[0])
		}

		// Hour
		startHour, err := strconv.Atoi(startTimeTokens[0])
		if err != nil {
			return nil, errors.New("Unable to parse number for start hour " + startTimeTokens[0])
		}

		// Minutes
		startMinutes := 0
		if len(startTimeTokens) == 2 {
			startMinutes, err = strconv.Atoi(startTimeTokens[1])
			if err != nil {
				return nil, errors.New("Unable to parse number for start minutes " + startTimeTokens[1])
			}
		}

		startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), startHour, startMinutes, 0, 0, time.Local)

		// End time
		endTimeTokens := strings.Split(timeTokens[1], hourMinutesSeparator)
		if len(endTimeTokens) == 2 && (endTimeTokens[0] == "" || endTimeTokens[1] == "") {
			return nil, errors.New("Unable to parse string " + timeTokens[1])
		}

		// Hour
		endHour, err := strconv.Atoi(endTimeTokens[0])
		if err != nil {
			return nil, errors.New("Unable to parse number for end hour " + endTimeTokens[0])
		}

		// Minutes
		endMinutes := 0
		if len(endTimeTokens) == 2 {
			endMinutes, err = strconv.Atoi(endTimeTokens[1])
			if err != nil {
				return nil, errors.New("Unable to parse number for end minutes " + endTimeTokens[1])
			}
		}
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), endHour, endMinutes, 0, 0, time.Local)

		if startTime.After(endTime) || startTime.Equal(endTime) {
			return nil, errors.New("Start time should be after end time for interval " + intervalToken)
		}

		timeIntervals[intervalToken] = []time.Time{startTime, endTime}
	}
	return timeIntervals, nil
}
