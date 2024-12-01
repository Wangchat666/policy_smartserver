package utils

import "time"

func StringToTime(input string) (time.Time, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func StringToTimeByDot(input string) (time.Time, error) {
	layout := "2006.01.02"
	t, err := time.Parse(layout, input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
