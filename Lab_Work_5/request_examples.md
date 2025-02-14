## Post запрос на register
```
curl -v -X POST -H "Content-Type: application/json" -d "{\"username\":\"testuser1\",\"password\":\"testpass1\"}" http://localhost:8000/register
```
## Post запрос на login
```
curl -v -X POST -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"password\":\"testpass\"}" http://localhost:8000/login
```
## Post запрос на создание задачи
```
curl -v -X POST -H "Content-Type: application/json" -d "{\"id\":1, \"text\":\"First task\"}" http://localhost:8001/tasks
```
## GET запроса на получение всех задач
```
curl -X GET http://localhost:8001/tasks/list
```
//