// car_reminders.cs — C# версия

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;

class Reminder
{
    public int Id { get; set; }
    public string Title { get; set; }
    public string Date { get; set; }
    public int Mileage { get; set; }
    public string Description { get; set; }
    public string Priority { get; set; }
    public bool Done { get; set; }
}

class Program
{
    private static List<Reminder> reminders = new List<Reminder>();
    private const string DataFile = "reminders.json";

    static void Main()
    {
        Load();
        while (true)
        {
            Console.WriteLine("\n\u001B[36m🔔 Дневник напоминаний (C#)\u001B[0m");
            Console.WriteLine("1. Добавить напоминание");
            Console.WriteLine("2. Показать все напоминания");
            Console.WriteLine("3. Показать предстоящие напоминания");
            Console.WriteLine("4. Поиск напоминаний");
            Console.WriteLine("5. Удалить напоминание");
            Console.WriteLine("6. Отметить как выполненное");
            Console.WriteLine("7. Статистика");
            Console.WriteLine("8. Выход");
            Console.Write("Выберите действие: ");
            string choice = Console.ReadLine();
            switch (choice)
            {
                case "1": AddReminder(); break;
                case "2": ListAll(); break;
                case "3": ListUpcoming(); break;
                case "4": SearchReminders(); break;
                case "5": DeleteReminder(); break;
                case "6": MarkDone(); break;
                case "7": Stats(); break;
                case "8": Console.WriteLine("До свидания!"); return;
                default: Console.WriteLine("\u001B[31mНеверный выбор.\u001B[0m"); break;
            }
        }
    }

    static void Load()
    {
        if (File.Exists(DataFile))
        {
            try
            {
                string json = File.ReadAllText(DataFile);
                reminders = JsonSerializer.Deserialize<List<Reminder>>(json) ?? new List<Reminder>();
            }
            catch { reminders = new List<Reminder>(); }
        }
    }

    static void Save()
    {
        string json = JsonSerializer.Serialize(reminders, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(DataFile, json);
    }

    static void AddReminder()
    {
        Console.Write("Название: ");
        string title = Console.ReadLine();
        Console.Write("Дата (ГГГГ-ММ-ДД): ");
        string date = Console.ReadLine();
        Console.Write("Пробег (км): ");
        int mileage = int.Parse(Console.ReadLine());
        Console.Write("Описание: ");
        string desc = Console.ReadLine();
        Console.Write("Приоритет (низкий/средний/высокий): ");
        string priority = Console.ReadLine().ToLower();
        if (priority != "низкий" && priority != "средний" && priority != "высокий") priority = "средний";
        int id = reminders.Count + 1;
        reminders.Add(new Reminder { Id = id, Title = title, Date = date, Mileage = mileage, Description = desc, Priority = priority, Done = false });
        Save();
        Console.WriteLine($"\u001B[32m✅ Напоминание добавлено (ID: {id})\u001B[0m");
    }

    static void ListAll()
    {
        if (reminders.Count == 0)
        {
            Console.WriteLine("\u001B[33mНет напоминаний.\u001B[0m");
            return;
        }
        Console.WriteLine($"\u001B[36m{"ID",-4} {"Название",-20} {"Дата",-12} {"Пробег",-10} {"Приоритет",-10} {"Статус",-12}\u001B[0m");
        Console.WriteLine(new string('-', 80));
        foreach (var r in reminders)
        {
            string status = r.Done ? "✅ Выполнено" : "⏳ Ожидает";
            string color = r.Priority == "низкий" ? "\u001B[32m" : r.Priority == "средний" ? "\u001B[33m" : "\u001B[31m";
            Console.WriteLine($"{r.Id,-4} {r.Title,-20} {r.Date,-12} {r.Mileage,-10} {color}{r.Priority,-10}\u001B[0m {status,-12}");
        }
    }

    static void ListUpcoming()
    {
        Console.Write("Количество дней (по умолч. 7): ");
        string daysStr = Console.ReadLine();
        int days = string.IsNullOrEmpty(daysStr) ? 7 : int.Parse(daysStr);
        var now = DateTime.Now.Date;
        var upcoming = reminders.Where(r => !r.Done && DateTime.Parse(r.Date) >= now && DateTime.Parse(r.Date) <= now.AddDays(days)).ToList();
        if (upcoming.Count == 0)
        {
            Console.WriteLine($"\u001B[32mНет напоминаний на ближайшие {days} дней.\u001B[0m");
            return;
        }
        Console.WriteLine($"\u001B[36m📅 Напоминания на ближайшие {days} дней:\u001B[0m");
        foreach (var r in upcoming)
        {
            int diff = (DateTime.Parse(r.Date) - now).Days;
            Console.WriteLine($"  {r.Id}: {r.Title} — {r.Date} (осталось {diff} дн)");
        }
    }

    static void SearchReminders()
    {
        Console.Write("Введите ключевое слово: ");
        string keyword = Console.ReadLine().ToLower();
        var results = reminders.Where(r => r.Title.ToLower().Contains(keyword) || r.Description.ToLower().Contains(keyword)).ToList();
        if (results.Count == 0)
            Console.WriteLine("\u001B[33mНичего не найдено.\u001B[0m");
        else
            results.ForEach(r => Console.WriteLine($"{r.Id}: {r.Title} | {r.Date} | {r.Priority} | {(r.Done ? "✅" : "⏳")}"));
    }

    static void DeleteReminder()
    {
        ListAll();
        Console.Write("Введите ID для удаления: ");
        int id = int.Parse(Console.ReadLine());
        var item = reminders.FirstOrDefault(r => r.Id == id);
        if (item != null)
        {
            reminders.Remove(item);
            Save();
            Console.WriteLine("\u001B[32m✅ Напоминание удалено.\u001B[0m");
        }
        else
            Console.WriteLine("\u001B[31m❌ Напоминание не найдено.\u001B[0m");
    }

    static void MarkDone()
    {
        ListAll();
        Console.Write("Введите ID для отметки как выполненное: ");
        int id = int.Parse(Console.ReadLine());
        var item = reminders.FirstOrDefault(r => r.Id == id);
        if (item != null)
        {
            item.Done = true;
            Save();
            Console.WriteLine("\u001B[32m✅ Напоминание отмечено как выполненное.\u001B[0m");
        }
        else
            Console.WriteLine("\u001B[31m❌ Напоминание не найдено.\u001B[0m");
    }

    static void Stats()
    {
        int total = reminders.Count;
        int done = reminders.Count(r => r.Done);
        int pending = total - done;
        var now = DateTime.Now.Date;
        int overdue = reminders.Count(r => !r.Done && DateTime.Parse(r.Date) < now);
        Console.WriteLine("\u001B[36m📊 Статистика:\u001B[0m");
        Console.WriteLine($"  Всего напоминаний: {total}");
        Console.WriteLine($"  Выполнено: {done}");
        Console.WriteLine($"  Ожидает: {pending}");
        Console.WriteLine($"  Просрочено: {overdue}");
    }
}
