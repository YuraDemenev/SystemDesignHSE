Диаграмма контейнеров 
Программный продукт представляет из себя мобильное приложение на телефоны с операционной системой Android. Помимо приложения на телефон, будет разработан микросервис, отвечающий за синхронизацию прогресса пользователя. Он будет сохранять прогресс пользователя, также, аналитик сможет получить от этого микросервиса данные для анализа. 
![Диаграмма_контейнеров для мобильного приложения](C4Diograms_Containers1.jpg "Подсказка")
![Диаграмма_контейнеров для сервера для синхронизации](C4Diograms_Containers2.jpg "Подсказка")
 
Диаграмма компонентов 
![Диаграмма_компонентов для мобильного приложения](C4Diogram_Components1.jpg "Подсказка")
![Диаграмма_компонентов для сервера для синхронизации](C4Diogram_Components2.jpg "Подсказка")
 
Диаграмма последовательностей
![Диаграмма_последовательностей](sequence.jpg "Подсказка")

База данных UML
![База_данных_UML](dbuml_1_jpeg.jpg "Подсказка")

## Принцип KISS (Keep It Simple, Stupid)
Простая функция для проверки ответа пользователя
```
fun isAnswerCorrect(userAnswer: String, correctAnswer: String): Boolean {
    return userAnswer.trim().equals(correctAnswer.trim(), ignoreCase = true)
}
```
Функция проста, читается легко, выполняет ровно одну задачу — сравнивает ответ пользователя с правильным ответом, убирая пробелы и учитывая регистр.
## Принцип YAGNI (You Aren’t Gonna Need It)
Начальный минимальный функционал для уведомлений
```
fun scheduleNotification(timeInMillis: Long, word: String) {
    // Используем AlarmManager для планирования уведомления
    val intent = Intent(context, NotificationReceiver::class.java).apply {
        putExtra("WORD", word)
    }
    val pendingIntent = PendingIntent.getBroadcast(context, 0, intent, PendingIntent.FLAG_UPDATE_CURRENT)
    val alarmManager = context.getSystemService(Context.ALARM_SERVICE) as AlarmManager
    alarmManager.setExact(AlarmManager.RTC_WAKEUP, timeInMillis, pendingIntent)
}
```
Функция планирует уведомление, но без избыточных конфигураций, например, без поддержки повторений или сложных таймеров. Это минимально необходимый функционал.
## Принцип DRY (Don't Repeat Yourself)
Общая функция для получения слов из базы данных
```
fun getWordsByCategory(category: String): List<Word> {
    return wordDao.getWordsByCategory(category)
}

fun getWordsByLevel(level: String): List<Word> {
    return wordDao.getWordsByLevel(level)
}

fun fetchWords(filterType: String, filterValue: String): List<Word> {
    return when (filterType) {
        "CATEGORY" -> getWordsByCategory(filterValue)
        "LEVEL" -> getWordsByLevel(filterValue)
        else -> emptyList()
    }
}
```
Вместо того чтобы писать отдельные функции в каждом месте, была создана обобщенная функция fetchWords, которая вызывает конкретные функции в зависимости от типа фильтра.
# Принцип SOLID
## Single Responsibility Principle (Принцип единственной ответственности)
Суть: У класса должна быть только одна причина для изменения.
```
class WordManager(private val wordRepository: WordRepository) {

    fun getNextWordForReview(): Word? {
        return wordRepository.fetchNextWord()
    }

    fun updateWordProgress(word: Word, success: Boolean) {
        wordRepository.updateWordProgress(word, success)
    }
}
```
WordManager отвечает только за логику управления словами. Другие задачи (например, хранение данных) делегируются репозиторию.
## O — Open/Closed Principle (Принцип открытости/закрытости)
Классы должны быть открыты для расширения, но закрыты для изменения
```
interface ReviewAlgorithm {
    fun calculateNextReviewDate(currentDate: Long, success: Boolean): Long
}

class BasicAlgorithm : ReviewAlgorithm {
    override fun calculateNextReviewDate(currentDate: Long, success: Boolean): Long {
        return if (success) currentDate + ONE_DAY else currentDate + ONE_HOUR
    }
}
```
Можно добавить новые алгоритмы, не меняя существующий код.
## L — Liskov Substitution Principle (Принцип подстановки Барбары Лисков)
Объекты дочернего класса должны заменять объекты родительского класса без нарушения функциональности.
```
interface Notifier {
    fun sendNotification(word: String)
}

class PushNotifier : Notifier {
    override fun sendNotification(word: String) {
        // Отправка push-уведомления
    }
}

fun notifyUser(notifier: Notifier, word: String) {
    notifier.sendNotification(word)
}
```
Функция notifyUser может работать с любым типом Notifier (PushNotifier, EmailNotifier).
## I — Interface Segregation Principle (Принцип разделения интерфейсов)
Интерфейсы должны быть узкоспециализированными.
```
interface WordFetcher {
    fun fetchNextWord(): Word
}

interface WordUpdater {
    fun updateProgress(word: Word, success: Boolean)
}

class WordRepository : WordFetcher, WordUpdater {
    override fun fetchNextWord(): Word {
        // Получение следующего слова
    }

    override fun updateProgress(word: Word, success: Boolean) {
        // Обновление прогресса
    }
}
```
Интерфейсы разделены на WordFetcher и WordUpdater для конкретных задач.
## D — Dependency Inversion Principle (Принцип инверсии зависимостей)
Классы должны зависеть от абстракций, а не от конкретных реализаций.
```
interface NotificationService {
    fun sendNotification(message: String)
}

class PushNotificationService : NotificationService {
    override fun sendNotification(message: String) {
        // Реализация push-уведомления
    }
}

class UserNotifier(private val notificationService: NotificationService) {
    fun notifyUser(word: String) {
        notificationService.sendNotification("Time to review: $word")
    }
}
```