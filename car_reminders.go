// car_reminders.go — Go версия

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Reminder struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Date        string `json:"date"`
	Mileage     int    `json:"mileage"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Done        bool   `json:"done"`
}

type Manager struct {
	Reminders []Reminder
	file      string
}

func NewManager(file string) *Manager {
	m := &Manager{file: file}
	m.load()
	return m
}

func (m *Manager) load() {
	data, err := os.ReadFile(m.file)
	if err != nil {
		m.Reminders = []Reminder{}
		return
	}
	json.Unmarshal(data, &m.Reminders)
}

func (m *Manager) save() {
	data, _ := json.MarshalIndent(m.Reminders, "", "  ")
	os.WriteFile(m.file, data, 0644)
}

func (m *Manager) add(title, date string, mileage int, desc, priority string) int {
	id := len(m.Reminders) + 1
	m.Reminders = append(m.Reminders, Reminder{
		ID:          id,
		Title:       title,
		Date:        date,
		Mileage:     mileage,
		Description: desc,
		Priority:    priority,
		Done:        false,
	})
	m.save()
	return id
}

func (m *Manager) listAll() {
	if len(m.Reminders) == 0 {
		fmt.Println("\u001B[33mНет напоминаний.\u001B[0m")
		return
	}
	fmt.Printf("\u001B[36m%-4s %-20s %-12s %-10s %-10s %-12s\u001B[0m\n", "ID", "Название", "Дата", "Пробег", "Приоритет", "Статус")
	fmt.Println(strings.Repeat("-", 80))
	for _, r := range m.Reminders {
		status := "⏳ Ожидает"
		if r.Done {
			status = "✅ Выполнено"
		}
		priorityColor := ""
		switch r.Priority {
		case "низкий":
			priorityColor = "\u001B[32m"
		case "средний":
			priorityColor = "\u001B[33m"
		case "высокий":
			priorityColor = "\u001B[31m"
		}
		fmt.Printf("%-4d %-20s %-12s %-10d %s%-10s\u001B[0m %-12s\n", r.ID, r.Title, r.Date, r.Mileage, priorityColor, r.Priority, status)
	}
}

func (m *Manager) listUpcoming(days int) {
	now := time.Now()
	upcoming := []Reminder{}
	for _, r := range m.Reminders {
		if r.Done {
			continue
		}
		t, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			continue
		}
		if t.After(now) && t.Before(now.AddDate(0, 0, days)) {
			upcoming = append(upcoming, r)
		}
	}
	if len(upcoming) == 0 {
		fmt.Printf("\u001B[32mНет напоминаний на ближайшие %d дней.\u001B[0m\n", days)
		return
	}
	fmt.Printf("\u001B[36m📅 Напоминания на ближайшие %d дней:\u001B[0m\n", days)
	for _, r := range upcoming {
		t, _ := time.Parse("2006-01-02", r.Date)
		diff := int(t.Sub(now).Hours() / 24)
		fmt.Printf("  %d: %s — %s (осталось %d дн)\n", r.ID, r.Title, r.Date, diff)
	}
}

func (m *Manager) search(keyword string) {
	found := false
	for _, r := range m.Reminders {
		if strings.Contains(strings.ToLower(r.Title), strings.ToLower(keyword)) ||
			strings.Contains(strings.ToLower(r.Description), strings.ToLower(keyword)) {
			fmt.Printf("%d: %s | %s | %s | %s\n", r.ID, r.Title, r.Date, r.Priority, map[bool]string{true: "✅", false: "⏳"}[r.Done])
			found = true
		}
	}
	if !found {
		fmt.Println("\u001B[33mНичего не найдено.\u001B[0m")
	}
}

func (m *Manager) delete(id int) bool {
	for i, r := range m.Reminders {
		if r.ID == id {
			m.Reminders = append(m.Reminders[:i], m.Reminders[i+1:]...)
			m.save()
			return true
		}
	}
	return false
}

func (m *Manager) markDone(id int) bool {
	for i, r := range m.Reminders {
		if r.ID == id {
			m.Reminders[i].Done = true
			m.save()
			return true
		}
	}
	return false
}

func (m *Manager) stats() {
	total := len(m.Reminders)
	done := 0
	overdue := 0
	now := time.Now()
	for _, r := range m.Reminders {
		if r.Done {
			done++
		} else {
			t, err := time.Parse("2006-01-02", r.Date)
			if err == nil && t.Before(now) {
				overdue++
			}
		}
	}
	pending := total - done
	fmt.Println("\u001B[36m📊 Статистика:\u001B[0m")
	fmt.Printf("  Всего напоминаний: %d\n", total)
	fmt.Printf("  Выполнено: %d\n", done)
	fmt.Printf("  Ожидает: %d\n", pending)
	fmt.Printf("  Просрочено: %d\n", overdue)
}

func main() {
	manager := NewManager("reminders.json")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n\u001B[36m🔔 Дневник напоминаний (Go)\u001B[0m")
		fmt.Println("1. Добавить напоминание")
		fmt.Println("2. Показать все напоминания")
		fmt.Println("3. Показать предстоящие напоминания")
		fmt.Println("4. Поиск напоминаний")
		fmt.Println("5. Удалить напоминание")
		fmt.Println("6. Отметить как выполненное")
		fmt.Println("7. Статистика")
		fmt.Println("8. Выход")
		fmt.Print("Выберите действие: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			fmt.Print("Название: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			fmt.Print("Дата (ГГГГ-ММ-ДД): ")
			date, _ := reader.ReadString('\n')
			date = strings.TrimSpace(date)
			fmt.Print("Пробег (км): ")
			mileageStr, _ := reader.ReadString('\n')
			mileage, _ := strconv.Atoi(strings.TrimSpace(mileageStr))
			fmt.Print("Описание: ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)
			fmt.Print("Приоритет (низкий/средний/высокий): ")
			priority, _ := reader.ReadString('\n')
			priority = strings.TrimSpace(strings.ToLower(priority))
			if priority != "низкий" && priority != "средний" && priority != "высокий" {
				priority = "средний"
			}
			id := manager.add(title, date, mileage, desc, priority)
			fmt.Printf("\u001B[32m✅ Напоминание добавлено (ID: %d)\u001B[0m\n", id)
		case "2":
			manager.listAll()
		case "3":
			fmt.Print("Количество дней (по умолч. 7): ")
			daysStr, _ := reader.ReadString('\n')
			days := 7
			if daysStr = strings.TrimSpace(daysStr); daysStr != "" {
				days, _ = strconv.Atoi(daysStr)
			}
			manager.listUpcoming(days)
		case "4":
			fmt.Print("Введите ключевое слово: ")
			keyword, _ := reader.ReadString('\n')
			keyword = strings.TrimSpace(keyword)
			manager.search(keyword)
		case "5":
			manager.listAll()
			fmt.Print("Введите ID для удаления: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			if manager.delete(id) {
				fmt.Println("\u001B[32m✅ Напоминание удалено.\u001B[0m")
			} else {
				fmt.Println("\u001B[31m❌ Напоминание не найдено.\u001B[0m")
			}
		case "6":
			manager.listAll()
			fmt.Print("Введите ID для отметки как выполненное: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			if manager.markDone(id) {
				fmt.Println("\u001B[32m✅ Напоминание отмечено как выполненное.\u001B[0m")
			} else {
				fmt.Println("\u001B[31m❌ Напоминание не найдено.\u001B[0m")
			}
		case "7":
			manager.stats()
		case "8":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("\u001B[31mНеверный выбор.\u001B[0m")
		}
	}
}
