<?php
// car_reminders.php — PHP версия

$dataFile = 'reminders.json';

function loadReminders() {
    global $dataFile;
    if (file_exists($dataFile)) {
        $json = file_get_contents($dataFile);
        return json_decode($json, true) ?: [];
    }
    return [];
}

function saveReminders($reminders) {
    global $dataFile;
    file_put_contents($dataFile, json_encode($reminders, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
}

$reminders = loadReminders();

function color($text, $code) {
    return "\033[{$code}m{$text}\033[0m";
}

while (true) {
    echo "\n" . color("🔔 Дневник напоминаний (PHP)", '36') . "\n";
    echo "1. Добавить напоминание\n";
    echo "2. Показать все напоминания\n";
    echo "3. Показать предстоящие напоминания\n";
    echo "4. Поиск напоминаний\n";
    echo "5. Удалить напоминание\n";
    echo "6. Отметить как выполненное\n";
    echo "7. Статистика\n";
    echo "8. Выход\n";
    echo "Выберите действие: ";
    $choice = trim(fgets(STDIN));

    switch ($choice) {
        case '1':
            echo "Название: ";
            $title = trim(fgets(STDIN));
            echo "Дата (ГГГГ-ММ-ДД): ";
            $date = trim(fgets(STDIN));
            echo "Пробег (км): ";
            $mileage = (int) trim(fgets(STDIN));
            echo "Описание: ";
            $desc = trim(fgets(STDIN));
            echo "Приоритет (низкий/средний/высокий): ";
            $priority = strtolower(trim(fgets(STDIN)));
            if (!in_array($priority, ['низкий', 'средний', 'высокий'])) $priority = 'средний';
            $id = count($reminders) + 1;
            $reminders[] = [
                'id' => $id,
                'title' => $title,
                'date' => $date,
                'mileage' => $mileage,
                'description' => $desc,
                'priority' => $priority,
                'done' => false
            ];
            saveReminders($reminders);
            echo color("✅ Напоминание добавлено (ID: $id)\n", '32');
            break;

        case '2':
            if (empty($reminders)) {
                echo color("Нет напоминаний.\n", '33');
            } else {
                printf(color("%-4s %-20s %-12s %-10s %-10s %-12s\n", '36'), "ID", "Название", "Дата", "Пробег", "Приоритет", "Статус");
                echo str_repeat("-", 80) . "\n";
                foreach ($reminders as $r) {
                    $status = $r['done'] ? "✅ Выполнено" : "⏳ Ожидает";
                    $color = $r['priority'] == 'низкий' ? '32' : ($r['priority'] == 'средний' ? '33' : '31');
                    printf("%-4d %-20s %-12s %-10d %s%-10s\033[0m %-12s\n", $r['id'], $r['title'], $r['date'], $r['mileage'], color("", $color), $r['priority'], $status);
                }
            }
            break;

        case '3':
            echo "Количество дней (по умолч. 7): ";
            $daysStr = trim(fgets(STDIN));
            $days = $daysStr ? (int)$daysStr : 7;
            $now = new DateTime();
            $found = false;
            echo color("📅 Напоминания на ближайшие $days дней:\n", '36');
            foreach ($reminders as $r) {
                if ($r['done']) continue;
                try {
                    $d = new DateTime($r['date']);
                    $diff = $now->diff($d);
                    if ($d >= $now && $d <= (clone $now)->modify("+$days days")) {
                        echo "  {$r['id']}: {$r['title']} — {$r['date']} (осталось {$diff->days} дн)\n";
                        $found = true;
                    }
                } catch (Exception $e) {}
            }
            if (!$found) echo color("Нет напоминаний на ближайшие $days дней.\n", '32');
            break;

        case '4':
            echo "Введите ключевое слово: ";
            $keyword = strtolower(trim(fgets(STDIN)));
            $found = false;
            foreach ($reminders as $r) {
                if (stripos($r['title'], $keyword) !== false || stripos($r['description'], $keyword) !== false) {
                    echo "{$r['id']}: {$r['title']} | {$r['date']} | {$r['priority']} | " . ($r['done'] ? '✅' : '⏳') . "\n";
                    $found = true;
                }
            }
            if (!$found) echo color("Ничего не найдено.\n", '33');
            break;

        case '5':
            if (empty($reminders)) { echo color("Нет напоминаний.\n", '33'); break; }
            foreach ($reminders as $r) {
                echo "{$r['id']}: {$r['title']} — {$r['date']}\n";
            }
            echo "Введите ID для удаления: ";
            $id = (int) trim(fgets(STDIN));
            $index = array_search($id, array_column($reminders, 'id'));
            if ($index !== false) {
                array_splice($reminders, $index, 1);
                saveReminders($reminders);
                echo color("✅ Напоминание удалено.\n", '32');
            } else {
                echo color("❌ Напоминание не найдено.\n", '31');
            }
            break;

        case '6':
            if (empty($reminders)) { echo color("Нет напоминаний.\n", '33'); break; }
            foreach ($reminders as $r) {
                echo "{$r['id']}: {$r['title']} — {$r['date']} (" . ($r['done'] ? '✅ выполнено' : '⏳ ожидает') . ")\n";
            }
            echo "Введите ID для отметки как выполненное: ";
            $id = (int) trim(fgets(STDIN));
            foreach ($reminders as &$r) {
                if ($r['id'] == $id) {
                    $r['done'] = true;
                    saveReminders($reminders);
                    echo color("✅ Напоминание отмечено как выполненное.\n", '32');
                    break 2;
                }
            }
            echo color("❌ Напоминание не найдено.\n", '31');
            break;

        case '7':
            $total = count($reminders);
            $done = count(array_filter($reminders, function($r) { return $r['done']; }));
            $pending = $total - $done;
            $now = new DateTime();
            $overdue = 0;
            foreach ($reminders as $r) {
                if (!$r['done']) {
                    try {
                        if (new DateTime($r['date']) < $now) $overdue++;
                    } catch (Exception $e) {}
                }
            }
            echo color("📊 Статистика:\n", '36');
            echo "  Всего напоминаний: $total\n";
            echo "  Выполнено: $done\n";
            echo "  Ожидает: $pending\n";
            echo "  Просрочено: $overdue\n";
            break;

        case '8':
            echo "До свидания!\n";
            exit(0);

        default:
            echo color("Неверный выбор.\n", '31');
    }
}
?>
