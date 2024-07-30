package parser

import (
	"errors"
	"github.com/shev-dm/TODO-project/internal/models"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Repeat struct {
	Rules  string
	days   []int
	months []int
}

func split(str string) ([]int, error) {
	var numbers []int
	part := strings.Split(str, ",")
	for _, day := range part {
		dayInt, err := strconv.Atoi(day)
		if err != nil {
			return numbers, err
		}
		numbers = append(numbers, dayInt)
	}
	return numbers, nil
}

func parser(repeat string) (Repeat, error) {
	var rules Repeat
	parts := strings.Split(repeat, " ")
	switch len(parts) {
	case 1:
		rules.Rules = parts[0]
		return rules, nil
	case 2:
		rules.Rules = parts[0]
		days, err := split(parts[1])
		if err != nil {
			return rules, err
		}
		rules.days = days
		return rules, nil

	case 3:
		rules.Rules = parts[0]
		days, err := split(parts[1])
		if err != nil {
			return rules, err
		}
		rules.days = days

		months, err := split(parts[2])
		if err != nil {
			return rules, err
		}
		rules.months = months

		return rules, nil
	default:
		return rules, errors.New("ошибка в формате правила")
	}
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("не указано правило повтора задачи")
	}
	rules, err := parser(repeat)
	if err != nil {
		return "", err
	}

	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	oneDay := 24 * time.Hour
	now = now.Truncate(oneDay)
	dateTime = dateTime.Truncate(oneDay)

	switch rules.Rules {
	case "d":
		if rules.months != nil || len(rules.days) != 1 {
			return "", errors.New("неверный формат данных")
		}
		day := rules.days[0]
		if day > 400 {
			return "", errors.New("количество дней более 400")
		}
		dateTime = dateTime.AddDate(0, 0, day)
		for {
			if !now.Before(dateTime) {
				dateTime = dateTime.AddDate(0, 0, day)
			} else {
				break
			}
		}
		return dateTime.Format("20060102"), nil
	case "y":
		if rules.months != nil || rules.days != nil {
			return "", errors.New("неверный формат данных")
		}
		dateTime = dateTime.AddDate(1, 0, 0)
		for {
			if !now.Before(dateTime) {
				dateTime = dateTime.AddDate(1, 0, 0)
			} else {
				break
			}
		}
		return dateTime.Format("20060102"), nil
	case "w":
		if rules.months != nil || rules.days == nil || len(rules.days) > 7 {
			return "", errors.New("неверный формат данных")
		}
		for _, day := range rules.days {
			if day > 7 || day < 0 {
				return "", errors.New("неверный формат данных")
			}
		}
		var allNextDaysRulesW []int
		for _, day := range rules.days {
			today := int(now.Weekday())
			if today == 0 {
				today = 7
			}
			diff := (day - today + 7) % 7
			// diff==0 значит ближайший день это сегодня, нужно добавить 7, чтобы получить ближайшее число не сегодня
			if diff == 0 {
				allNextDaysRulesW = append(allNextDaysRulesW, 7)
				continue
			}
			// diff от 1 до 6, значит столько дней до ближайшего дня недели
			allNextDaysRulesW = append(allNextDaysRulesW, diff)
		}
		nextDayInt := slices.Min(allNextDaysRulesW)
		dateTime = now.AddDate(0, 0, nextDayInt)
		return dateTime.Format("20060102"), nil
	case "m":
		if rules.months == nil {
			var startYear int
			var startMonth time.Month
			if now.Before(dateTime) {
				startYear = dateTime.Year()
				startMonth = dateTime.Month()
			} else {
				startYear = now.Year()
				startMonth = now.Month()
			}
			var allNextDaysRulesM []time.Time
			for _, day := range rules.days {
				if day < -2 || day > 31 {
					return "", errors.New("неверный формат данных")
				}
				if day > 0 {
					currentMonth := time.Date(startYear, startMonth, day, 0, 0, 0, 0, dateTime.Location())
					if currentMonth.Day() == day {
						allNextDaysRulesM = append(allNextDaysRulesM, currentMonth)
					}
					nextMonth := time.Date(startYear, startMonth+1, day, 0, 0, 0, 0, dateTime.Location())
					if nextMonth.Day() == day {
						allNextDaysRulesM = append(allNextDaysRulesM, nextMonth)
					}
					monthAfterNext := time.Date(startYear, startMonth+2, day, 0, 0, 0, 0, dateTime.Location())
					if monthAfterNext.Day() == day {
						allNextDaysRulesM = append(allNextDaysRulesM, monthAfterNext)
					}
					continue
				}
				if day == -1 || day == -2 {
					currentMonth := time.Date(startYear, startMonth+1, day+1, 0, 0, 0, 0, dateTime.Location())
					nextMonth := time.Date(startYear, startMonth+2, day+1, 0, 0, 0, 0, dateTime.Location())
					monthAfterNext := time.Date(startYear, startMonth+3, day+1, 0, 0, 0, 0, dateTime.Location())
					allNextDaysRulesM = append(allNextDaysRulesM, currentMonth, nextMonth, monthAfterNext)
				}
			}
			sort.Slice(allNextDaysRulesM, func(i, j int) bool {
				return allNextDaysRulesM[i].Before(allNextDaysRulesM[j])
			})
			for _, day := range allNextDaysRulesM {
				if day.After(now) && day.After(dateTime) {
					dateTime = day
					break
				}
			}
			return dateTime.Format("20060102"), nil
		}
		var startYear int
		if now.Before(dateTime) {
			startYear = dateTime.Year()
		} else {
			startYear = now.Year()
		}
		var allNextDaysRulesM []time.Time
		for _, day := range rules.days {
			if day < -2 || day > 31 {
				return "", errors.New("неверный формат данных")
			}
			for _, month := range rules.months {
				if month < 0 || month > 12 {
					return "", errors.New("неверный формат данных")
				}
				if day > 0 {
					currentYear := time.Date(startYear, time.Month(month), day, 0, 0, 0, 0, dateTime.Location())
					if currentYear.Day() == day {
						allNextDaysRulesM = append(allNextDaysRulesM, currentYear)
					}
					nextYear := time.Date(startYear+1, time.Month(month), day, 0, 0, 0, 0, dateTime.Location())
					if nextYear.Day() == day {
						allNextDaysRulesM = append(allNextDaysRulesM, nextYear)
					}
					continue
				}
				if day == -1 || day == -2 {
					currentYear := time.Date(startYear, time.Month(month)+1, day+1, 0, 0, 0, 0, dateTime.Location())
					nextYear := time.Date(startYear+1, time.Month(month)+1, day+1, 0, 0, 0, 0, dateTime.Location())
					allNextDaysRulesM = append(allNextDaysRulesM, currentYear, nextYear)
				}
			}
		}
		sort.Slice(allNextDaysRulesM, func(i, j int) bool {
			return allNextDaysRulesM[i].Before(allNextDaysRulesM[j])
		})
		for _, day := range allNextDaysRulesM {
			if day.After(now) && day.After(dateTime) {
				dateTime = day
				break
			}
		}
		return dateTime.Format("20060102"), nil
	default:
		return "", errors.New("неверный формат данных")
	}
}

func CheckRulesAddOrUpdate(input models.Task) (models.Task, models.Err) {
	var dateTime time.Time
	var err error
	var errorAnswer models.Err

	if input.Title == "" {
		errorAnswer.Err = "пустое значение Title"
		return input, errorAnswer
	}

	if input.Date == "" {
		dateTime = time.Now()
		input.Date = dateTime.Format("20060102")
	} else {
		dateTime, err = time.Parse("20060102", input.Date)
		if err != nil {
			errorAnswer.Err = "ошибка в формате даты"
			return input, errorAnswer
		}
	}

	if input.Repeat != "" {
		nextDate, err := NextDate(time.Now(), input.Date, input.Repeat)
		if err != nil {
			errorAnswer.Err = err.Error()
			return input, errorAnswer
		}
		if dateTime.Before(time.Now().Truncate(24 * time.Hour)) {
			input.Date = nextDate
		}
	} else {
		if dateTime.Before(time.Now().Truncate(24 * time.Hour)) {
			input.Date = time.Now().Format("20060102")
		}
	}
	return input, errorAnswer
}
