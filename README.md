**Запуск проекта**
1. Копируем переменные
```shell
cp .env.example .env
```
2. Запускаем проект
```shell
docker-compose up -d
```

3. API запросы

***Получение новостей***
```bash
curl --location --request GET 'http://localhost:8080/news/' \
--header 'Content-Type: application/json'
```
***Получение конкретной новости (с комментариями)***
```bash
curl --location --request GET 'http://localhost:8080/news/{id}' \
--header 'Content-Type: application/json'
```
***Получение комментариев новости***
```bash
curl --location --request GET 'http://localhost:8080/news/{id}/comments/' \
--header 'Content-Type: application/json'
```
***Добавление комментария к новости***
```bash
curl --location --request POST 'http://localhost:8080/news/{id}/comments/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "newsId": 1,
    "content": "asdsadqwert",
    "parentId": 0
}'
```