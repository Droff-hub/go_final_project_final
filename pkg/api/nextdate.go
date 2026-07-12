package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// NextDate вычисляет следующую дату выполнения задачи.
// Параметры:
//   now     – текущее время (от которого ищем следующую дату)
//   dstart  – исходная дата в формате DateFormat
//   repeat  – правило повторения
// Возвращает: следующую дату в формате DateFormat и ошибку.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Пустое правило – ошибка
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не может быть пустым")
	}

	// Парсим начальную дату
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	// Разбиваем правило на части
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила")
	}

	switch parts[0] {
	case "d":
		// правило "d <число>"
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d: требуется число")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверное число дней: %v", err)
		}
		if days < 1 || days > 400 {
			return "", fmt.Errorf("число дней должно быть от 1 до 400")
		}
		// Сдвигаем дату на days дней до тех пор, пока она не станет больше now
		for {
			date = date.AddDate(0, 0, days)
			if !date.Before(now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "y":
		// ежегодное повторение
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат правила y")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "w":
		// правило недель (сделал для задания со звёздочкой)
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила w: требуется список дней")
		}
		daysStr := strings.Split(parts[1], ",")
		allowed := make([]bool, 8) // индексы 1..7 (пн=1, вс=7)
		for _, s := range daysStr {
			d, err := strconv.Atoi(s)
			if err != nil || d < 1 || d > 7 {
				return "", fmt.Errorf("неверный день недели: %s", s)
			}
			allowed[d] = true
		}
		// Ищем ближайший допустимый день недели, начиная со следующего дня от now
		// Начинаем проверку с tomorrow (чтобы дата была строго больше now)
		check := now.AddDate(0, 0, 1)
		for {
			// Проверяем, совпадает ли день недели (0-воскресенье, 1-понедельник...)
			weekday := int(check.Weekday()) // 0=Sun, 1=Mon, ..., 6=Sat
			// Преобразуем в наш формат: 1=Mon, 7=Sun
			dow := weekday
			if dow == 0 {
				dow = 7
			}
			if allowed[dow] {
				// Проверяем, что дата >= исходной dstart. Нужно, чтобы дата была > now,
				// мы начали с tomorrow, так что условие выполнено.
				return check.Format(DateFormat), nil
			}
			check = check.AddDate(0, 0, 1)
		}

	case "m":
		// правило месяцев (тоже, для задания со звёздочкой)
		if len(parts) < 2 {
			return "", fmt.Errorf("неверный формат правила m: требуется список дней")
		}
		// Парсим дни
		daysStr := strings.Split(parts[1], ",")
		allowedDays := make([]bool, 32) // 1..31
		for _, s := range daysStr {
			d, err := strconv.Atoi(s)
			if err != nil || d < -2 || d == 0 || d > 31 {
				return "", fmt.Errorf("неверный день месяца: %s", s)
			}
			if d < 0 {
				// Для отрицательных значений нужно будет вычислять с конца месяца
				// сохраняем как есть, обработаем позже
				// Пока просто сохраняем в мап для дальнейшей обработки
				allowedDays[d] = true
			} else {
				allowedDays[d] = true
			}
		}
		// Месяцы 
		var allowedMonths []int
		if len(parts) >= 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, s := range monthsStr {
				m, err := strconv.Atoi(s)
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("неверный месяц: %s", s)
				}
				allowedMonths = append(allowedMonths, m)
			}
		} else {
			// Если месяцы не указаны,значит разрешены все
			for m := 1; m <= 12; m++ {
				allowedMonths = append(allowedMonths, m)
			}
		}
		// Создаём карту для быстрой проверки месяцев
		monthsMap := make(map[int]bool)
		for _, m := range allowedMonths {
			monthsMap[m] = true
		}

		// Ищем следующую дату, начиная с tomorrow
		check := now.AddDate(0, 0, 1)
		for {
			year, month, day := check.Date()
			// Проверяем месяц
			if !monthsMap[int(month)] {
				check = check.AddDate(0, 1, 0)
				continue
			}
			// Проверяем день
			// Если день отрицательный, вычисляем с конца месяца
			needMatch := false
			for d, ok := range allowedDays {
				if !ok {
					continue
				}
				var targetDay int
				if d < 0 {
					// последний, предпоследний и т.д. и т.п.
					lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
					targetDay = lastDay + d + 1 // например, -1 → lastDay, -2 → lastDay-1
				} else {
					targetDay = d
				}
				if day == targetDay {
					needMatch = true
					break
				}
			}
			if needMatch {
				return check.Format(DateFormat), nil
			}
			check = check.AddDate(0, 0, 1)
		}

	default:
		return "", fmt.Errorf("неподдерживаемое правило: %s", parts[0])
	}
}

// nextDateHandler обрабатывает GET /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Если now не передан, используем текущую дату
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат now", http.StatusBadRequest)
			return
		}
	}

	// Вызываем функцию
	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}