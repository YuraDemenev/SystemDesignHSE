## Post запрос на register
```
curl -v -X POST -H "Content-Type: application/json" -d "{\"username\":\"testuser1\",\"password\":\"testpass1\"}" http://5.35.37.117:8000/register
```
## Post запрос на login
```
curl -v -X POST -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"password\":\"testpass\"}" http://5.35.37.117:8000/login
```
## Post запрос на создание задачи
```
curl -v -X POST -H "Content-Type: application/json" -d "{\"id\":1, \"text\":\"First task\"}" http://5.35.37.117:8001/tasks
```
## GET запроса на получение всех задач
```
curl -X GET http://5.35.37.117:8001/tasks/list
```