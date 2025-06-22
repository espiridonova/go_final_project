package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout = "20060102"
	daysLimit  = 400
	sunday     = 7
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", err
	}

	var resDate time.Time
	repeatParams := strings.Split(repeat, " ")

	switch repeatParams[0] {
	case "d":
		resDate, err = addDateDays(date, now, repeatParams)
		if err != nil {
			return "", err
		}
	case "y":
		resDate = addDateYear(date, now)
	case "w":
		resDate, err = addDateWeekDays(date, now, repeatParams)
		if err != nil {
			return "", err
		}
	case "m":
		resDate, err = addDateMonthDays(date, now, repeatParams)
		if err != nil {
			return "", err
		}
	case "":
		return "", nil
	default:
		return "", errors.New("invalid repeat")
	}

	return resDate.Format(dateLayout), nil
}

func prepareDays(param string) (int, error) {
	days, err := strconv.Atoi(param)
	if err != nil || days <= 0 || days > daysLimit {
		return 0, errors.New("invalid days in repeat")
	}
	return days, nil
}

func prepareWeekDays(params string) (map[int]struct{}, error) {
	weekDaysParams := strings.Split(params, ",")
	weekDays := make(map[int]struct{}, len(weekDaysParams))
	for _, wd := range weekDaysParams {
		weekDay, err := strconv.Atoi(wd)
		if err != nil {
			return nil, errors.New("invalid week day in repeat")
		}
		if weekDay < 1 || weekDay > sunday {
			return nil, errors.New("invalid week day")
		}
		if weekDay == sunday {
			weekDays[0] = struct{}{}
		}
		weekDays[weekDay] = struct{}{}
	}
	return weekDays, nil
}

func convertParams(params string) ([]int, error) {
	p := strings.Split(params, ",")
	r := make([]int, 0, len(p))

	for _, param := range p {
		v, err := strconv.Atoi(param)
		if err != nil {
			return nil, errors.New("invalid param")
		}
		r = append(r, v)
	}

	return r, nil
}

func prepareMonthDays(days, months []int) (map[int][]int, error) {
	res := make(map[int][]int)

	if len(months) == 0 {
		// add all month
		for i := 1; i < 13; i++ {
			months = append(months, i)
		}
	}

	for _, month := range months {
		if month < 1 || month > 12 {
			return nil, errors.New("invalid month")
		}

		for _, day := range days {
			if day < -2 || day > 31 || day == 0 {
				return nil, errors.New("invalid day")
			}
			res[month] = append(res[month], day)
		}
	}

	return res, nil
}

func addDateDays(date, now time.Time, repeatParams []string) (time.Time, error) {
	if len(repeatParams) < 2 {
		return time.Time{}, errors.New("invalid repeat")
	}
	days, err := prepareDays(repeatParams[1])
	if err != nil {
		return time.Time{}, err
	}

	if days == 1 && cmpDates(date, now) {
		return date, nil
	}

	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			return date, nil
		}
	}
}

func addDateYear(date, now time.Time) time.Time {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			return date
		}
	}
}

func addDateWeekDays(date, now time.Time, repeatParams []string) (time.Time, error) {
	if len(repeatParams) < 2 {
		return time.Time{}, errors.New("invalid repeat")
	}

	weekDays, err := prepareWeekDays(repeatParams[1])
	if err != nil {
		return time.Time{}, err
	}

	for {
		date = date.AddDate(0, 0, 1)
		weekDay := int(date.Weekday())
		_, ok := weekDays[weekDay]
		if !ok {
			continue
		}
		if afterNow(date, now) {
			return date, nil
		}
	}
}

func addDateMonthDays(date, now time.Time, repeatParams []string) (time.Time, error) {
	var days []int
	var months []int
	var err error

	if len(repeatParams) < 2 {
		return time.Time{}, errors.New("invalid month repeat")
	}

	if len(repeatParams) >= 2 {
		days, err = convertParams(repeatParams[1])
		if err != nil {
			return time.Time{}, err
		}
	}

	if len(days) == 0 {
		return time.Time{}, errors.New("invalid month days")
	}

	if len(repeatParams) == 3 {
		months, err = convertParams(repeatParams[2])
		if err != nil {
			return time.Time{}, err
		}
	}

	monthDays, err := prepareMonthDays(days, months)
	if err != nil {
		return time.Time{}, err
	}

	if len(monthDays) == 0 {
		return time.Time{}, errors.New("failed add date by month days; invalid param")
	}

	for {
		date = date.AddDate(0, 0, 1)
		lastDay := getLastDayOfMonth(date)

		month := int(date.Month())
		if len(monthDays[month]) == 0 {
			continue
		}

		posDays := make(map[int]struct{})
		negDays := make(map[int]struct{})
		for _, d := range monthDays[month] {
			if d > 0 {
				posDays[d] = struct{}{}
			} else {
				negDays[d] = struct{}{}
			}
		}

		day := date.Day()

		_, posOK := posDays[day]
		_, negOK := negDays[day-(lastDay+1)]

		if (posOK || negOK) && afterNow(date, now) {
			return date, nil
		}
	}
}

func getLastDayOfMonth(date time.Time) int {
	year, month, _ := date.Date()
	loc := date.Location()

	temp := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	_, _, lastDayOfMonth := temp.AddDate(0, 1, -1).Date()

	return lastDayOfMonth
}

func afterNow(date, now time.Time) bool {
	date1 := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	date2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return date1.Unix() > date2.Unix()
}

func cmpDates(date, now time.Time) bool {
	date1 := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	date2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return date1.Unix() == date2.Unix()
}
