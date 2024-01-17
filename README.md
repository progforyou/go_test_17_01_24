# Тестовое задание

Конфигурация файла .env

PORT=8080 //порт

DNS="host=localhost user=nikolai password=nikolai dbname=persons" //база



**Реализованы следующие CRUD функции**

**Список**:

Фильтрация
`curl --request GET \
--url 'http://localhost:8080/person/?filter=%7B%22name%22%3A%20%22Dmitriy%22%7D' \`

Пагинация
`curl --request GET \
--url 'http://localhost:8080/person/?page=1&per_page=2' \`


**Создание**:

`curl --request POST \
--url http://localhost:8080/person/ \
--header 'Content-Type: application/json' \
--data '{
"name": "Dmitriy",
"surname": "Ushakov",
"patronymic": "Vasilevich"
}'`


**Изменение**:

`curl --request PUT \
--url http://localhost:8080/person/1/ \
--header 'Content-Type: application/json' \
--data '{
"name": "Dimitry",
"surname": "Ushakov",
"patronymic": "Vasilevich"
}'`


**Удаление**:

`curl --request DELETE \
--url http://localhost:8080/person/1/`