# ТЗ для CSV сервиса (для практического проекта #2)

## Реализовать REST API для хранения информации о продукции. Сущность "продукт" имеет следующий вид
```sql
CREATE TABLE products (
   id serial not null unique,
   name varchar(255) not null unique,
   price integer not null default 0
)
```

## Эндпоинты  
- **POST /products -** создание продукта
- **PUT /products/:id -** апдейт продукта по айди
- **DELETE /products/:id -** удаление продукта по айди
- **GET /products -** получение CSV файла со списком всех продуктов

CSV файл должен выглядеть следующим образом:  
```
PRODUCT NAME;PRICE
```

## Дополнительно
- Добавить аутентификацию
