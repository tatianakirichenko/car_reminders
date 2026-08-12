// car_reminders.js — JavaScript версия

const fs = require('fs');
const readline = require('readline');

const DATA_FILE = 'reminders.json';
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

let reminders = [];

function load() {
    if (fs.existsSync(DATA_FILE)) {
        try {
            reminders = JSON.parse(fs.readFileSync(DATA_FILE, 'utf8'));
        } catch (e) {
            reminders = [];
        }
    }
}

function save() {
    fs.writeFileSync(DATA_FILE, JSON.stringify(reminders, null, 2));
}

function ask(question) {
    return new Promise(resolve => rl.question(question, resolve));
}

function color(text, code) {
    return `\x1b[${code}m${text}\x1b[0m`;
}

async function main() {
    load();
    while (true) {
        console.log(`\n${color('🔔 Дневник напоминаний (JavaScript)', '36')}`);
        console.log("1. Добавить напоминание");
        console.log("2. Показать все напоминания");
        console.log("3. Показать предстоящие напоминания");
        console.log("4. Поиск напоминаний");
        console.log("5. Удалить напоминание");
        console.log("6. Отметить как выполненное");
        console.log("7. Статистика");
        console.log("8. Выход");
        const choice = await ask("Выберите действие: ");
        switch (choice.trim()) {
            case "1": await addReminder(); break;
            case "2": listAll(); break;
            case "3": await listUpcoming(); break;
            case "4": await searchReminders(); break;
            case "5": await deleteReminder(); break;
            case "6": await markDone(); break;
            case "7": stats(); break;
            case "8": console.log("До свидания!"); rl.close(); return;
            default: console.log(color("Неверный выбор.", "31"));
        }
    }
}

async function addReminder() {
    const title = await ask("Название: ");
    const date = await ask("Дата (ГГГГ-ММ-ДД): ");
    const mileage = parseInt(await ask("Пробег (км): "));
    const desc = await ask("Описание: ");
    let priority = (await ask("Приоритет (низкий/средний/высокий): ")).toLowerCase();
    if (!["низкий", "средний", "высокий"].includes(priority)) priority = "средний";
    const id = reminders.length + 1;
    reminders.push({ id, title, date, mileage, description: desc, priority, done: false });
    save();
    console.log(color(`✅ Напоминание добавлено (ID: ${id})`, "32"));
}

function listAll() {
    if (reminders.length === 0) {
        console.log(color("Нет напоминаний.", "33"));
        return;
    }
    console.log(color(`${'ID'.padEnd(4)} ${'Название'.padEnd(20)} ${'Дата'.padEnd(12)} ${'Пробег'.padEnd(10)} ${'Приоритет'.padEnd(10)} ${'Статус'.padEnd(12)}`, "36"));
    console.log("-".repeat(80));
    for (const r of reminders) {
        const status = r.done ? "✅ Выполнено" : "⏳ Ожидает";
        const priorityColor = r.priority === "низкий" ? "32" : r.priority === "средний" ? "33" : "31";
        console.log(`${String(r.id).padEnd(4)} ${r.title.padEnd(20)} ${r.date.padEnd(12)} ${String(r.mileage).padEnd(10)} ${color(r.priority.padEnd(10), priorityColor)} ${status}`);
    }
}

async function listUpcoming() {
    const days = parseInt(await ask("Количество дней (по умолч. 7): ")) || 7;
    const now = new Date();
    const upcoming = reminders.filter(r => {
        if (r.done) return false;
        const d = new Date(r.date);
        return d >= now && d <= new Date(now.getTime() + days * 24 * 60 * 60 * 1000);
    });
    if (upcoming.length === 0) {
        console.log(color(`Нет напоминаний на ближайшие ${days} дней.`, "32"));
        return;
    }
    console.log(color(`📅 Напоминания на ближайшие ${days} дней:`, "36"));
    for (const r of upcoming) {
        const diff = Math.ceil((new Date(r.date) - now) / (1000 * 60 * 60 * 24));
        console.log(`  ${r.id}: ${r.title} — ${r.date} (осталось ${diff} дн)`);
    }
}

async function searchReminders() {
    const keyword = (await ask("Введите ключевое слово: ")).toLowerCase();
    const results = reminders.filter(r => r.title.toLowerCase().includes(keyword) || r.description.toLowerCase().includes(keyword));
    if (results.length === 0) {
        console.log(color("Ничего не найдено.", "33"));
    } else {
        results.forEach(r => console.log(`${r.id}: ${r.title} | ${r.date} | ${r.priority} | ${r.done ? '✅' : '⏳'}`));
    }
}

async function deleteReminder() {
    listAll();
    const id = parseInt(await ask("Введите ID для удаления: "));
    const index = reminders.findIndex(r => r.id === id);
    if (index !== -1) {
        reminders.splice(index, 1);
        save();
        console.log(color("✅ Напоминание удалено.", "32"));
    } else {
        console.log(color("❌ Напоминание не найдено.", "31"));
    }
}

async function markDone() {
    listAll();
    const id = parseInt(await ask("Введите ID для отметки как выполненное: "));
    const r = reminders.find(r => r.id === id);
    if (r) {
        r.done = true;
        save();
        console.log(color("✅ Напоминание отмечено как выполненное.", "32"));
    } else {
        console.log(color("❌ Напоминание не найдено.", "31"));
    }
}

function stats() {
    const total = reminders.length;
    const done = reminders.filter(r => r.done).length;
    const pending = total - done;
    const now = new Date();
    const overdue = reminders.filter(r => !r.done && new Date(r.date) < now).length;
    console.log(color("📊 Статистика:", "36"));
    console.log(`  Всего напоминаний: ${total}`);
    console.log(`  Выполнено: ${done}`);
    console.log(`  Ожидает: ${pending}`);
    console.log(`  Просрочено: ${overdue}`);
}

main().catch(console.error);
