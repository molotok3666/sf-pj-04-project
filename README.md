**Запуск проекта**
1. Копируем переменные
```shell
cp .env.example .env
```
2. Запускаем проект
```shell
docker-compose up -d
```
3. Стучимся по URL
```link
localhost
```
4. API доступно по URL (count - количество новостей)
```link
http://localhost/news/{count}
```
