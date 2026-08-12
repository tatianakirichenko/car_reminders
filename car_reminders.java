// car_reminders.java — Java версия

import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.time.*;
import java.time.format.*;

class Reminder {
    int id;
    String title;
    String date;
    int mileage;
    String description;
    String priority;
    boolean done;

    Reminder(int id, String title, String date, int mileage, String description, String priority, boolean done) {
        this.id = id;
        this.title = title;
        this.date = date;
        this.mileage = mileage;
        this.description = description;
        this.priority = priority;
        this.done = done;
    }

    String toJson() {
        return String.format("{\"id\":%d,\"title\":\"%s\",\"date\":\"%s\",\"mileage\":%d,\"description\":\"%s\",\"priority\":\"%s\",\"done\":%b}",
                id, title, date, mileage, description, priority, done);
    }

    static Reminder fromJson(String json) {
        // Упрощённый парсинг для демонстрации
        return null;
    }
}

public class car_reminders {
    private static final String DATA_FILE = "reminders.json";
    private static List<Reminder> reminders = new ArrayList<>();
    private static Scanner scanner = new Scanner(System.in);

    public static void main(String[] args) {
        load();
        while (true) {
            System.out.println("\n\u001B[36m🔔 Дневник напоминаний (Java)\u001B[0m");
            System.out.println("1. Добавить напоминание");
            System.out.println("2. Показать все напоминания");
            System.out.println("3. Показать предстоящие напоминания");
            System.out.println("4. Поиск напоминаний");
            System.out.println("5. Удалить напоминание");
            System.out.println("6. Отметить как выполненное");
            System.out.println("7. Статистика");
            System.out.println("8. Выход");
            System.out.print("Выберите действие: ");
            String choice = scanner.nextLine();
            switch (choice) {
                case "1": addReminder(); break;
                case "2": listAll(); break;
                case "3": listUpcoming(); break;
                case "4": searchReminders(); break;
                case "5": deleteReminder(); break;
                case "6": markDone(); break;
                case "7": showStats(); break;
                case "8": System.out.println("До свидания!"); return;
                default: System.out.println("\u001B[31mНеверный выбор.\u001B[0m");
            }
        }
    }

    private static void load() {
        try {
            String content = new String(Files.readAllBytes(Paths.get(DATA_FILE)));
            // Упрощённо, в реальности нужен JSON-парсер
            reminders = new ArrayList<>();
        } catch (IOException e) {
            reminders = new ArrayList<>();
        }
    }

    private static void save() {
        try {
            StringBuilder sb = new StringBuilder("[");
            for (int i = 0; i < reminders.size(); i++) {
                sb.append(reminders.get(i).toJson());
                if (i < reminders.size() - 1) sb.append(",");
            }
            sb.append("]");
            Files.write(Paths.get(DATA_FILE), sb.toString().getBytes());
        } catch (IOException e) {
            System.out.println("Ошибка сохранения.");
        }
    }

    private static void addReminder() {
        System.out.print("Название: ");
        String title = scanner.nextLine();
        System.out.print("Дата (ГГГГ-ММ-ДД): ");
        String date = scanner.nextLine();
        System.out.print("Пробег (км): ");
        int mileage = Integer.parseInt(scanner.nextLine());
        System.out.print("Описание: ");
        String desc = scanner.nextLine();
        System.out.print("Приоритет (низкий/средний/высокий): ");
        String priority = scanner.nextLine().toLowerCase();
        if (!priority.equals("низкий") && !priority.equals("средний") && !priority.equals("высокий")) {
            priority = "средний";
        }
        int id = reminders.size() + 1;
        reminders.add(new Reminder(id, title, date, mileage, desc, priority, false));
        save();
        System.out.println("\u001B[32m✅ Напоминание добавлено (ID: " + id + ")\u001B[0m");
    }

    private static void listAll() {
        if (reminders.isEmpty()) {
            System.out.println("\u001B[33mНет напоминаний.\u001B[0m");
            return;
        }
        System.out.printf("\u001B[36m%-4s %-20s %-12s %-10s %-10s %-12s\u001B[0m\n", "ID", "Название", "Дата", "Пробег", "Приоритет", "Статус");
        System.out.println("-".repeat(80));
        for (Reminder r : reminders) {
            String status = r.done ? "✅ Выполнено" : "⏳ Ожидает";
            String color = r.priority.equals("низкий") ? "\u001B[32m" :
                           r.priority.equals("средний") ? "\u001B[33m" : "\u001B[31m";
            System.out.printf("%-4d %-20s %-12s %-10d %s%-10s\u001B[0m %-12s\n", r.id, r.title, r.date, r.mileage, color, r.priority, status);
        }
    }

    private static void listUpcoming() {
        System.out.print("Количество дней (по умолч. 7): ");
        String daysStr = scanner.nextLine();
        int days = daysStr.isEmpty() ? 7 : Integer.parseInt(daysStr);
        LocalDate now = LocalDate.now();
        System.out.printf("\u001B[36m📅 Напоминания на ближайшие %d дней:\u001B[0m\n", days);
        for (Reminder r : reminders) {
            if (r.done) continue;
            try {
                LocalDate d = LocalDate.parse(r.date);
                if (d.isAfter(now) && d.isBefore(now.plusDays(days))) {
                    long diff = d.toEpochDay() - now.toEpochDay();
                    System.out.printf("  %d: %s — %s (осталось %d дн)\n", r.id, r.title, r.date, diff);
                }
            } catch (Exception e) {}
        }
    }

    private static void searchReminders() {
        System.out.print("Введите ключевое слово: ");
        String keyword = scanner.nextLine().toLowerCase();
        boolean found = false;
        for (Reminder r : reminders) {
            if (r.title.toLowerCase().contains(keyword) || r.description.toLowerCase().contains(keyword)) {
                System.out.printf("%d: %s | %s | %s | %s\n", r.id, r.title, r.date, r.priority, r.done ? "✅" : "⏳");
                found = true;
            }
        }
        if (!found) System.out.println("\u001B[33mНичего не найдено.\u001B[0m");
    }

    private static void deleteReminder() {
        listAll();
        System.out.print("Введите ID для удаления: ");
        int id = Integer.parseInt(scanner.nextLine());
        boolean removed = reminders.removeIf(r -> r.id == id);
        if (removed) {
            save();
            System.out.println("\u001B[32m✅ Напоминание удалено.\u001B[0m");
        } else {
            System.out.println("\u001B[31m❌ Напоминание не найдено.\u001B[0m");
        }
    }

    private static void markDone() {
        listAll();
        System.out.print("Введите ID для отметки как выполненное: ");
        int id = Integer.parseInt(scanner.nextLine());
        for (Reminder r : reminders) {
            if (r.id == id) {
                r.done = true;
                save();
                System.out.println("\u001B[32m✅ Напоминание отмечено как выполненное.\u001B[0m");
                return;
            }
        }
        System.out.println("\u001B[31m❌ Напоминание не найдено.\u001B[0m");
    }

    private static void showStats() {
        int total = reminders.size();
        int done = 0, overdue = 0;
        LocalDate now = LocalDate.now();
        for (Reminder r : reminders) {
            if (r.done) done++;
            else {
                try {
                    if (LocalDate.parse(r.date).isBefore(now)) overdue++;
                } catch (Exception e) {}
            }
        }
        System.out.println("\u001B[36m📊 Статистика:\u001B[0m");
        System.out.println("  Всего напоминаний: " + total);
        System.out.println("  Выполнено: " + done);
        System.out.println("  Ожидает: " + (total - done));
        System.out.println("  Просрочено: " + overdue);
    }
}
