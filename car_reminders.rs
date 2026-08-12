// car_reminders.rs — Rust версия

use serde::{Deserialize, Serialize};
use std::fs;
use std::io::{self, Write};
use chrono::{NaiveDate, Utc};

#[derive(Serialize, Deserialize, Clone)]
struct Reminder {
    id: usize,
    title: String,
    date: String,
    mileage: u32,
    description: String,
    priority: String,
    done: bool,
}

struct Manager {
    reminders: Vec<Reminder>,
    file: String,
}

impl Manager {
    fn new(file: &str) -> Self {
        let mut m = Manager { reminders: Vec::new(), file: file.to_string() };
        m.load();
        m
    }

    fn load(&mut self) {
        if let Ok(data) = fs::read_to_string(&self.file) {
            if let Ok(reminders) = serde_json::from_str(&data) {
                self.reminders = reminders;
                return;
            }
        }
        self.reminders = Vec::new();
    }

    fn save(&self) {
        let data = serde_json::to_string_pretty(&self.reminders).unwrap();
        fs::write(&self.file, data).unwrap();
    }

    fn add(&mut self, title: String, date: String, mileage: u32, description: String, priority: String) -> usize {
        let id = self.reminders.len() + 1;
        self.reminders.push(Reminder {
            id,
            title,
            date,
            mileage,
            description,
            priority,
            done: false,
        });
        self.save();
        id
    }

    fn list_all(&self) {
        if self.reminders.is_empty() {
            println!("\x1b[33mНет напоминаний.\x1b[0m");
            return;
        }
        println!("\x1b[36m{:<4} {:<20} {:<12} {:<10} {:<10} {:<12}\x1b[0m", "ID", "Название", "Дата", "Пробег", "Приоритет", "Статус");
        println!("{}", "-".repeat(80));
        for r in &self.reminders {
            let status = if r.done { "✅ Выполнено" } else { "⏳ Ожидает" };
            let color = match r.priority.as_str() {
                "низкий" => "\x1b[32m",
                "средний" => "\x1b[33m",
                "высокий" => "\x1b[31m",
                _ => "",
            };
            println!("{:<4} {:<20} {:<12} {:<10} {}{:<10}\x1b[0m {:<12}", r.id, r.title, r.date, r.mileage, color, r.priority, status);
        }
    }

    fn list_upcoming(&self, days: i64) {
        let now = chrono::Local::now().date_naive();
        let mut found = false;
        println!("\x1b[36m📅 Напоминания на ближайшие {} дней:\x1b[0m", days);
        for r in &self.reminders {
            if r.done { continue; }
            if let Ok(d) = NaiveDate::parse_from_str(&r.date, "%Y-%m-%d") {
                if d >= now && d <= now + chrono::Days::new(days as u64) {
                    let diff = (d - now).num_days();
                    println!("  {}: {} — {} (осталось {} дн)", r.id, r.title, r.date, diff);
                    found = true;
                }
            }
        }
        if !found {
            println!("\x1b[32mНет напоминаний на ближайшие {} дней.\x1b[0m", days);
        }
    }

    fn search(&self, keyword: &str) {
        let kw = keyword.to_lowercase();
        let results: Vec<&Reminder> = self.reminders.iter()
            .filter(|r| r.title.to_lowercase().contains(&kw) || r.description.to_lowercase().contains(&kw))
            .collect();
        if results.is_empty() {
            println!("\x1b[33mНичего не найдено.\x1b[0m");
        } else {
            for r in results {
                println!("{}: {} | {} | {} | {}", r.id, r.title, r.date, r.priority, if r.done { "✅" } else { "⏳" });
            }
        }
    }

    fn delete(&mut self, id: usize) -> bool {
        let pos = self.reminders.iter().position(|r| r.id == id);
        if let Some(idx) = pos {
            self.reminders.remove(idx);
            self.save();
            true
        } else {
            false
        }
    }

    fn mark_done(&mut self, id: usize) -> bool {
        for r in &mut self.reminders {
            if r.id == id {
                r.done = true;
                self.save();
                return true;
            }
        }
        false
    }

    fn stats(&self) {
        let total = self.reminders.len();
        let done = self.reminders.iter().filter(|r| r.done).count();
        let pending = total - done;
        let now = chrono::Local::now().date_naive();
        let overdue = self.reminders.iter()
            .filter(|r| !r.done && NaiveDate::parse_from_str(&r.date, "%Y-%m-%d").map_or(false, |d| d < now))
            .count();
        println!("\x1b[36m📊 Статистика:\x1b[0m");
        println!("  Всего напоминаний: {}", total);
        println!("  Выполнено: {}", done);
        println!("  Ожидает: {}", pending);
        println!("  Просрочено: {}", overdue);
    }
}

fn main() {
    let mut manager = Manager::new("reminders.json");
    loop {
        println!("\n\x1b[36m🔔 Дневник напоминаний (Rust)\x1b[0m");
        println!("1. Добавить напоминание");
        println!("2. Показать все напоминания");
        println!("3. Показать предстоящие напоминания");
        println!("4. Поиск напоминаний");
        println!("5. Удалить напоминание");
        println!("6. Отметить как выполненное");
        println!("7. Статистика");
        println!("8. Выход");
        print!("Выберите действие: ");
        io::stdout().flush().unwrap();
        let mut choice = String::new();
        io::stdin().read_line(&mut choice).unwrap();
        match choice.trim() {
            "1" => {
                print!("Название: ");
                io::stdout().flush().unwrap();
                let mut title = String::new();
                io::stdin().read_line(&mut title).unwrap();
                let title = title.trim().to_string();
                print!("Дата (ГГГГ-ММ-ДД): ");
                io::stdout().flush().unwrap();
                let mut date = String::new();
                io::stdin().read_line(&mut date).unwrap();
                let date = date.trim().to_string();
                print!("Пробег (км): ");
                io::stdout().flush().unwrap();
                let mut mileage_str = String::new();
                io::stdin().read_line(&mut mileage_str).unwrap();
                let mileage: u32 = mileage_str.trim().parse().unwrap();
                print!("Описание: ");
                io::stdout().flush().unwrap();
                let mut desc = String::new();
                io::stdin().read_line(&mut desc).unwrap();
                let desc = desc.trim().to_string();
                print!("Приоритет (низкий/средний/высокий): ");
                io::stdout().flush().unwrap();
                let mut priority = String::new();
                io::stdin().read_line(&mut priority).unwrap();
                let priority = priority.trim().to_lowercase();
                let priority = if priority == "низкий" || priority == "средний" || priority == "высокий" { priority } else { "средний".to_string() };
                let id = manager.add(title, date, mileage, desc, priority);
                println!("\x1b[32m✅ Напоминание добавлено (ID: {})\x1b[0m", id);
            }
            "2" => manager.list_all(),
            "3" => {
                print!("Количество дней (по умолч. 7): ");
                io::stdout().flush().unwrap();
                let mut days_str = String::new();
                io::stdin().read_line(&mut days_str).unwrap();
                let days = if days_str.trim().is_empty() { 7 } else { days_str.trim().parse().unwrap() };
                manager.list_upcoming(days);
            }
            "4" => {
                print!("Введите ключевое слово: ");
                io::stdout().flush().unwrap();
                let mut keyword = String::new();
                io::stdin().read_line(&mut keyword).unwrap();
                manager.search(keyword.trim());
            }
            "5" => {
                manager.list_all();
                print!("Введите ID для удаления: ");
                io::stdout().flush().unwrap();
                let mut id_str = String::new();
                io::stdin().read_line(&mut id_str).unwrap();
                let id: usize = id_str.trim().parse().unwrap();
                if manager.delete(id) {
                    println!("\x1b[32m✅ Напоминание удалено.\x1b[0m");
                } else {
                    println!("\x1b[31m❌ Напоминание не найдено.\x1b[0m");
                }
            }
            "6" => {
                manager.list_all();
                print!("Введите ID для отметки как выполненное: ");
                io::stdout().flush().unwrap();
                let mut id_str = String::new();
                io::stdin().read_line(&mut id_str).unwrap();
                let id: usize = id_str.trim().parse().unwrap();
                if manager.mark_done(id) {
                    println!("\x1b[32m✅ Напоминание отмечено как выполненное.\x1b[0m");
                } else {
                    println!("\x1b[31m❌ Напоминание не найдено.\x1b[0m");
                }
            }
            "7" => manager.stats(),
            "8" => {
                println!("До свидания!");
                break;
            }
            _ => println!("\x1b[31mНеверный выбор.\x1b[0m"),
        }
    }
}
