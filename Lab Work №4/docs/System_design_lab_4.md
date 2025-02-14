# API 1
## 1 запрос: Регистрация нового пользователя RegisterUser
При первом заходе в приложения, сохраняем в БД на сервере нового пользователя. В запросе отправлям Google id, для регистрации нового пользователя
Метод Post<br>
Endpoint /users/register
Тело запроса JSON:
```
{
    "google_id": "1234567890abcdef"
}
```
Ответ тело запроса JSON:
```
{
    "exist_user": true
}
```

## 2 запрос: Синхронизация пользователя SyncUserData
Если во время 1 открытие приложения, запрос RegisterUser возвращает true, значит такой пользователь есть и так как это первое открытие приложение, нужно синхронизировать его локальные данные с сервером.<br>
Метод Get<br>
Endpoint /users/sync
Тело запроса JSON:
```
{
    "user_id":"550e8400-e29b-41d4-a716-446655440000"
}
```
Ответ тело запроса JSON:
```
{
    "words": [
    {
      "id": 1,
      "english_word": "apple",
      "russian_translation": "яблоко",
      "transcription_id": 101,
      "british_variable": "apple",
      "level_id": 2
    }
  ]
}
```
## 3 запрос: Сохранение статистики пользователя на сервере SaveUserData
В приложение введется статистика пользователя: кол-во выученных слов, кол-во слов, которые пользователь знает,
максимум дней, которые пользователь учит слова, кол-во слов, которые пользователь учит<br>
Метод Post<br>
Endpoint /users/user_data<br>
Тело запроса JSON:
```
{
    "user_id":"550e8400-e29b-41d4-a716-446655440000",
    "count_learned_words": 1000,
    "max_learn_words_days":50,
    "count_learning_words":500,
    "count_known_words":500
}
```
Ответ тело запроса JSON:
```
{
    "result_save":true
}
```
## 4 запрос: Сохранение прогресса пользователя за сессию SaveUserWords
Когда пользователь поучил слова и закрыл приложение, отправляется запрос на сохранение слов в БД на сервере, которые
пользователь выучил за сеанс.<br>
Метод Post<br>
Endpoint /users/user_words<br>
Тело запроса JSON:
```
{
    "user_id":"550e8400-e29b-41d4-a716-446655440000",
    user_words:[
        "id": 1,
        "english_word": "apple",
        "russian_translation": "яблоко",
        "transcription_id": 101,
        "british_variable": "apple",
        "level_id": 2
    ]
}
```
Ответ тело запроса JSON:
```
{
    "result_save":true
}
```
## 5 запрос: Удаление аккаунта пользователя из БД на сервере для синхронизации DeleteAccount<br>
Пользователь может полностью удалить все данные из БД.<br>
Метод DELETE<br>
Endpoint /users/user_delete<br>
Тело запроса JSON:<br>
```
{
    "user_id":"550e8400-e29b-41d4-a716-446655440000",
}
```
Ответ тело запроса JSON:
```
{
    "result":true
}
```
## Так как приложение мобильное, большинство запросов будут внутри приложения не в формате REST API.
## 6 запрос: Получение слов для изучения из локальной БД в приложении
Запрос на получение слов для изучения из локальной БД в приложении<br>
Метод fun<br>
Параметры запроса:<br>
context - Context приложения
db - ссылка на базу данных
Count learning words - кол-во слов, которые пользователь изучает<br>
list of levels - список тем, которые пользователь изучает, чтобы выбрать из них слова (Спорт, Здоровье...)<br>
Ответ:<br>
Pair<MutableList<Words>, HashMap<Int, String>> - Пара (тип данных Kotlin) содержащая список слов для изучения, HashMap содержащая уровень и id слова (Для UI элемента)<br>
## 7 запрос: Сохранение прогресса пользователя в локальном Storage в формате proto data
Метод fun<br>
Параметры запроса:<br>
User - класс пользователя со всеми данными пользователя<br>
ListsProto - список уровней в формамет LevelsProto для прямого сохранения<br>
listOfWordsIdsForRepeat - список words ids для повторения<br>
lastTimeLearned - дата когда последний раз пользователь учил слова<br>
```
class User(
    var userId: String = "",
    var curRepeatDays: Int = 0,
    var maxRepeatDays: Int = 0,
    var countFullLearnedWords: Int = 0,
    var countLearningWords: Int = 0,
    var countLearnedWordsToday: Int = 0,
    var checkLearnedAllWordsToday: Boolean = false,
    var countKnewWords: Int = 0,
    var listOfLevels: MutableList<Levels> = mutableListOf(),
    var checkBritishVariables: Boolean = false,
    var lastTimeLearnedWords: Instant = Instant.now(),
    var listOfWordsForRepeat: List<Words> = listOf()
)
```
## 8 запрос: Получение пользователя из локального хранилища
Метод fun<br>
Параметры запроса:<br>
Нет параметров<br>
Возвращаемый результат:<br>
User - класс пользователя со всеми данными пользователя<br>
# API 2 реализация
Реализация методов на языке Go<br>
## 1 Реализация метода RegisterUser
```
type UserId struct {
	GoogleID string `json:"google_id"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req UserId
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

    res:=false
	existUser := checkExistUser(req)
    if !existUser{
        res = RegisterUserResponse{ExistUser: existUser}
    }
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
```
## 2 Реализация метода SyncUserData
```
type UserId struct {
	GoogleID string `json:"google_id"`
}

func SyncUserDataHandler(w http.ResponseWriter, r *http.Request) {
	var req UserId
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Заглушка: получаем слова из БД
	words := getUserWords(req)

	res := SyncUserDataResponse{Words: words}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
```
## 3 Реализация метода SaveUserData
```
    type SaveUserDataRequest struct {
        UserID             string `json:"user_id"`
        CountLearnedWords  int    `json:"count_learned_words"`
        MaxLearnWordsDays  int    `json:"max_learn_words_days"`
        CountLearningWords int    `json:"count_learning_words"`
        CountKnownWords    int    `json:"count_known_words"`
    }

	var req SaveUserDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

    err:=saveUserData(req)
    if err!=nil{
        http.Error(w, "Invalid request", http.StatusBadRequest)
		return
    }

	res := SaveUserResponse{ResultSave: true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
```
## 4 Реализация метода SaveUserWords
```
type SaveUserWordsRequest struct {
	UserID    string `json:"user_id"`
	UserWords []Word `json:"user_words"`
}

    var req SaveUserWordsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err:=SaveUserWords(req)
    if err!=nil{
        http.Error(w, "Invalid request", http.StatusBadRequest)
		return
    }

	res := SaveUserResponse{ResultSave: true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
```
## 5 Реализация метода DeleteAccount
```
    type UserId struct {
	    GoogleID string `json:"google_id"`
    }

    func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err:=deleteAccount(req)
    if err!=nil{
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

	res := DeleteAccountResponse{Result: true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
```
## 6 Реализация метода getWordsForLearn
```
 override suspend fun getWordsForLearn(
        context: Context,
        db: MainDB,
        flowLevelsModel: FlowLevelsModel,
        countLearningWords: Int
    ): Pair<MutableList<Words>, HashMap<Int, String>> {
        val myScope = CoroutineScope(Dispatchers.IO)
        lateinit var listOfLevels: List<Levels>

        //Получаем уровни
        myScope.launch {
            val hashSet = flowLevelsModel.data.value
            if (hashSet != null) {
                if (hashSet.size == 0) {
                    Log.e("Main page model get words for learn.", "hash set size is zero")
                    throw Exception()
                }
                listOfLevels = db.getDao().getLevelsByNamesMultipleQueries(hashSet)
            } else {
                Log.e("Main page model", "getWordsForLearn, hash set is null")
            }
        }.join()

        //Получаем из list levels ids чтобы их использовать в следующем запросе
        val arrayLevelsIds: Array<Int> = Array(listOfLevels.size) { 0 }
        listOfLevels.forEachIndexed { index, level ->
            if (level.id != null) {
                arrayLevelsIds[index] = level.id
            }
        }

        //Получаем случайный список слов с levels_id
        lateinit var words: List<Words>
        myScope.launch {
            words =
                db.getDao().getWordsByLevelsIdsMultiplyQueries(arrayLevelsIds, countLearningWords)
        }.join()

        //Создаем hash map, чтобы хранить названия уровней по ids
        val hashMap = HashMap<Int, String>()
        listOfLevels.forEach { element ->
            if (element.id != null) {
                hashMap[element.id] = element.name
            } else {
                Log.e(
                    "Main page contract",
                    "getWordsForLearn, listOfLevels has element with null id"
                )
                throw Exception()
            }
        }
        val pair = Pair(words.toMutableList(), hashMap)
        return pair
    }
```
## 7 Реализация метода setUserProtoData
```
private suspend fun setUserProtoData(
        context: Context,
        user: User,
        listOfLevelsBuilders: MutableList<LevelsProto>,
        listOfWordsIdsForRepeat: List<Int>,
        lastTimeLearned: Timestamp
    ) {
        //Полностью обновляем пользователя
        context.userParamsDataStore.updateData { userPorto ->
            userPorto.toBuilder()
                //Очищаем user proto data
                .clearListOfWordsIdsForRepeat()
                .clearListOfLevels()
                //Меняем user proto data
                .setUserId(user.userId)
                .setCurRepeatDays(user.curRepeatDays)
                .setMaxRepeatDays(user.maxRepeatDays)
                .setCountFullLearnedWords(user.countFullLearnedWords)
                .setCountLearningWords(user.countLearningWords)
                .setCountLearnedWordsToday(user.countLearnedWordsToday)
                .setCheckLearnedAllWordsToday(user.checkLearnedAllWordsToday)
                .setCountKnewWords(user.countKnewWords)
                .addAllListOfLevels(listOfLevelsBuilders)
                .setCheckBritishVariables(false)
                .setLastTimeLearnedWords(lastTimeLearned)
                .addAllListOfWordsIdsForRepeat(listOfWordsIdsForRepeat)
                .build()
        }
    }
```
## 8 Реализация метода getUserProtoData
```
   private fun getUserProtoData(context: Context, db: MainDB): Flow<User> {
        //Получаем данные из хранилища Proto DataStore
        val userFlow: Flow<User> = context.userParamsDataStore.data.map { userProto ->
            User(
                userProto.userId,
                userProto.curRepeatDays,
                userProto.maxRepeatDays,
                userProto.countFullLearnedWords,
                userProto.countLearningWords,
                userProto.countLearnedWordsToday,
                userProto.checkLearnedAllWordsToday,
                userProto.countKnewWords,
                userProto.listOfLevelsList.map { levelsProto ->
                    convertLevelsProtoToLevels(levelsProto)
                }.toMutableList(),
                userProto.checkBritishVariables,
                Instant.ofEpochSecond(
                    userProto.lastTimeLearnedWords.seconds,
                    userProto.lastTimeLearnedWords.nanos.toLong()
                ),
                getWordsByIds(userProto.listOfWordsIdsForRepeatList, db)

            )
        }
        return userFlow
    }
```