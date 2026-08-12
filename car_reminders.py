# car_reminders.py — Python версия

import json
import os
from datetime import datetime, timedelta
from colorama import init, Fore, Style

init(autoreset=True)
DATA_FILE = "reminders.json"

class Reminder:
    def __init__(self, id, title, date, mileage, description, priority, done=False):
        self.id = id
        self.title = title
        self.date = date
        self.mileage = mileage
        self.description = description
        self.priority = priority
        self.done = done

    def to_dict(self):
        return {
            "id": self.id,
            "title": self.title,
            "date": self.date,
            "mileage": self.mileage,
            "description": self.description,
            "priority": self.priority,
            "done": self.done
        }

    @classmethod
    def from_dict(cls, data):
        return cls(data["id"], data["title"], data["date"], data["mileage"],
                   data["description"], data["priority"], data.get("done", False))

class ReminderManager:
    def __init__(self):
        self.reminders = []
        self.load()

    def load(self):
        if os.path.exists(DATA_FILE):
            try:
                with open(DATA_FILE, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                    self.reminders = [Reminder.from_dict(r) for r in data]
            except:
                self.reminders = []

    def save(self):
        with open(DATA_FILE, 'w', encoding='utf-8') as f:
            json.dump([r.to_dict() for r in self.reminders], f, indent=2, ensure_ascii=False)

    def add(self, title, date, mileage, description, priority):
        id = len(self.reminders) + 1
        reminder = Reminder(id, title, date, mileage, description, priority)
        self.reminders.append(reminder)
        self.save()
        return id

    def list_all(self):
        if not self.reminders:
            print(Fore.YELLOW + "Нет напоминаний.")
            return
        print(Fore.CYAN + f"{'ID':<4} {'Название':<20} {'Дата':<12} {'Пробег':<10} {'Приоритет':<10} {'Статус':<12}")
        print("-" * 80)
        for r in self.reminders:
            status = "✅ Выполнено" if r.done else "⏳ Ожидает"
            priority_color = Fore.GREEN if r.priority == "низкий" else Fore.YELLOW if r.priority == "средний" else Fore.RED
            print(f"{r.id:<4} {r.title:<20} {r.date:<12} {r.mileage:<10} {priority_color}{r.priority:<10}{Style.RESET_ALL} {status:<12}")

    def list_upcoming(self, days=7):
        now = datetime.now()
        upcoming = []
        for r in self.reminders:
            if r.done:
                continue
            try:
                dt = datetime.strptime(r.date, "%Y-%m-%d")
                if dt >= now and dt <= now + timedelta(days=days):
                    upcoming.append(r)
            except:
                continue
        if not upcoming:
            print(Fore.GREEN + f"Нет напоминаний на ближайшие {days} дней.")
            return
        print(Fore.CYAN + f"📅 Напоминания на ближайшие {days} дней:")
        for r in upcoming:
            print(f"  {r.id}: {r.title} — {r.date} (осталось {(datetime.strptime(r.date, '%Y-%m-%d') - now).days} дн)")

    def search(self, keyword):
        results = [r for r in self.reminders if keyword.lower() in r.title.lower() or keyword.lower() in r.description.lower()]
        if not results:
            print(Fore.YELLOW + "Ничего не найдено.")
        else:
            for r in results:
                print(f"{r.id}: {r.title} | {r.date} | {r.priority} | {'✅' if r.done else '⏳'}")

    def delete(self, id):
        for i, r in enumerate(self.reminders):
            if r.id == id:
                del self.reminders[i]
                self.save()
                return True
        return False

    def mark_done(self, id):
        for r in self.reminders:
            if r.id == id:
                r.done = True
                self.save()
                return True
        return False

    def stats(self):
        total = len(self.reminders)
        done = sum(1 for r in self.reminders if r.done)
        pending = total - done
        overdue = 0
        now = datetime.now()
        for r in self.reminders:
            if not r.done:
                try:
                    if datetime.strptime(r.date, "%Y-%m-%d") < now:
                        overdue += 1
                except:
                    continue
        print(Fore.CYAN + "📊 Статистика:")
        print(f"  Всего напоминаний: {total}")
        print(f"  Выполнено: {done}")
        print(f"  Ожидает: {pending}")
        print(f"  Просрочено: {overdue}")

def main():
    manager = ReminderManager()
    while True:
        print(Fore.CYAN + "\n🔔 Дневник напоминаний (Python)")
        print("1. Добавить напоминание")
        print("2. Показать все напоминания")
        print("3. Показать предстоящие напоминания")
        print("4. Поиск напоминаний")
        print("5. Удалить напоминание")
        print("6. Отметить как выполненное")
        print("7. Статистика")
        print("8. Выход")
        choice = input("Выберите действие: ").strip()
        if choice == "1":
            title = input("Название: ")
            date = input("Дата (ГГГГ-ММ-ДД): ")
            mileage = int(input("Пробег (км): "))
            desc = input("Описание: ")
            priority = input("Приоритет (низкий/средний/высокий): ").lower()
            if priority not in ["низкий", "средний", "высокий"]:
                priority = "средний"
            id = manager.add(title, date, mileage, desc, priority)
            print(Fore.GREEN + f"✅ Напоминание добавлено (ID: {id})")
        elif choice == "2":
            manager.list_all()
        elif choice == "3":
            days = input("Количество дней (по умолч. 7): ")
            days = int(days) if days.strip() else 7
            manager.list_upcoming(days)
        elif choice == "4":
            keyword = input("Введите ключевое слово: ")
            manager.search(keyword)
        elif choice == "5":
            manager.list_all()
            id = int(input("Введите ID для удаления: "))
            if manager.delete(id):
                print(Fore.GREEN + "✅ Напоминание удалено.")
            else:
                print(Fore.RED + "❌ Напоминание не найдено.")
        elif choice == "6":
            manager.list_all()
            id = int(input("Введите ID для отметки как выполненное: "))
            if manager.mark_done(id):
                print(Fore.GREEN + "✅ Напоминание отмечено как выполненное.")
            else:
                print(Fore.RED + "❌ Напоминание не найдено.")
        elif choice == "7":
            manager.stats()
        elif choice == "8":
            print("До свидания!")
            break
        else:
            print(Fore.RED + "Неверный выбор.")

if __name__ == "__main__":
    main()
