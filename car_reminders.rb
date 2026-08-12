# car_reminders.rb — Ruby версия

require 'json'
require 'date'

DATA_FILE = 'reminders.json'

class Reminder
  attr_accessor :id, :title, :date, :mileage, :description, :priority, :done

  def initialize(id, title, date, mileage, description, priority, done = false)
    @id = id
    @title = title
    @date = date
    @mileage = mileage
    @description = description
    @priority = priority
    @done = done
  end

  def to_h
    { id: @id, title: @title, date: @date, mileage: @mileage,
      description: @description, priority: @priority, done: @done }
  end

  def self.from_h(data)
    new(data[:id], data[:title], data[:date], data[:mileage],
        data[:description], data[:priority], data[:done] || false)
  end
end

class ReminderManager
  attr_reader :reminders

  def initialize
    @reminders = []
    load
  end

  def load
    if File.exist?(DATA_FILE)
      begin
        data = JSON.parse(File.read(DATA_FILE), symbolize_names: true)
        @reminders = data.map { |r| Reminder.from_h(r) }
      rescue
        @reminders = []
      end
    end
  end

  def save
    File.write(DATA_FILE, JSON.pretty_generate(@reminders.map(&:to_h)))
  end

  def add(title, date, mileage, description, priority)
    id = @reminders.size + 1
    @reminders << Reminder.new(id, title, date, mileage, description, priority)
    save
    id
  end

  def list_all
    if @reminders.empty?
      puts "\e[33mНет напоминаний.\e[0m"
      return
    end
    printf "\e[36m%-4s %-20s %-12s %-10s %-10s %-12s\e[0m\n", "ID", "Название", "Дата", "Пробег", "Приоритет", "Статус"
    puts "-" * 80
    @reminders.each do |r|
      status = r.done ? "✅ Выполнено" : "⏳ Ожидает"
      color = case r.priority
              when "низкий" then "\e[32m"
              when "средний" then "\e[33m"
              when "высокий" then "\e[31m"
              else ""
              end
      puts "%-4d %-20s %-12s %-10d %s%-10s\e[0m %-12s" % [r.id, r.title, r.date, r.mileage, color, r.priority, status]
    end
  end

  def list_upcoming(days = 7)
    now = Date.today
    upcoming = @reminders.select do |r|
      !r.done && Date.parse(r.date) >= now && Date.parse(r.date) <= now + days
    rescue
      false
    end
    if upcoming.empty?
      puts "\e[32mНет напоминаний на ближайшие #{days} дней.\e[0m"
      return
    end
    puts "\e[36m📅 Напоминания на ближайшие #{days} дней:\e[0m"
    upcoming.each do |r|
      diff = (Date.parse(r.date) - now).to_i
      puts "  #{r.id}: #{r.title} — #{r.date} (осталось #{diff} дн)"
    end
  end

  def search(keyword)
    results = @reminders.select { |r| r.title.downcase.include?(keyword.downcase) || r.description.downcase.include?(keyword.downcase) }
    if results.empty?
      puts "\e[33mНичего не найдено.\e[0m"
    else
      results.each { |r| puts "#{r.id}: #{r.title} | #{r.date} | #{r.priority} | #{r.done ? '✅' : '⏳'}" }
    end
  end

  def delete(id)
    found = @reminders.find { |r| r.id == id }
    if found
      @reminders.delete(found)
      save
      true
    else
      false
    end
  end

  def mark_done(id)
    found = @reminders.find { |r| r.id == id }
    if found
      found.done = true
      save
      true
    else
      false
    end
  end

  def stats
    total = @reminders.size
    done = @reminders.count(&:done)
    pending = total - done
    now = Date.today
    overdue = @reminders.count { |r| !r.done && Date.parse(r.date) < now rescue false }
    puts "\e[36m📊 Статистика:\e[0m"
    puts "  Всего напоминаний: #{total}"
    puts "  Выполнено: #{done}"
    puts "  Ожидает: #{pending}"
    puts "  Просрочено: #{overdue}"
  end
end

def main
  manager = ReminderManager.new
  loop do
    puts "\n\e[36m🔔 Дневник напоминаний (Ruby)\e[0m"
    puts "1. Добавить напоминание"
    puts "2. Показать все напоминания"
    puts "3. Показать предстоящие напоминания"
    puts "4. Поиск напоминаний"
    puts "5. Удалить напоминание"
    puts "6. Отметить как выполненное"
    puts "7. Статистика"
    puts "8. Выход"
    print "Выберите действие: "
    choice = gets.chomp
    case choice
    when "1"
      print "Название: "
      title = gets.chomp
      print "Дата (ГГГГ-ММ-ДД): "
      date = gets.chomp
      print "Пробег (км): "
      mileage = gets.chomp.to_i
      print "Описание: "
      desc = gets.chomp
      print "Приоритет (низкий/средний/высокий): "
      priority = gets.chomp.downcase
      priority = "средний" unless ["низкий", "средний", "высокий"].include?(priority)
      id = manager.add(title, date, mileage, desc, priority)
      puts "\e[32m✅ Напоминание добавлено (ID: #{id})\e[0m"
    when "2"
      manager.list_all
    when "3"
      print "Количество дней (по умолч. 7): "
      days_str = gets.chomp
      days = days_str.empty? ? 7 : days_str.to_i
      manager.list_upcoming(days)
    when "4"
      print "Введите ключевое слово: "
      keyword = gets.chomp
      manager.search(keyword)
    when "5"
      manager.list_all
      print "Введите ID для удаления: "
      id = gets.chomp.to_i
      if manager.delete(id)
        puts "\e[32m✅ Напоминание удалено.\e[0m"
      else
        puts "\e[31m❌ Напоминание не найдено.\e[0m"
      end
    when "6"
      manager.list_all
      print "Введите ID для отметки как выполненное: "
      id = gets.chomp.to_i
      if manager.mark_done(id)
        puts "\e[32m✅ Напоминание отмечено как выполненное.\e[0m"
      else
        puts "\e[31m❌ Напоминание не найдено.\e[0m"
      end
    when "7"
      manager.stats
    when "8"
      puts "До свидания!"
      break
    else
      puts "\e[31mНеверный выбор.\e[0m"
    end
  end
end

main if __FILE__ == $0
