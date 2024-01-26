# тестовое задание для Effective Mobile

Для запуска необходимо передать параметр DATABASE_DSN в env

пример DATABASE_DSN=postgres://user:password@localhost/DB_name?sslmode=disable

после запуска доступны 4 ручки
___
GET {address}:{port}/api/v1/person/

Для фильтрации используются имена полей (id, name, surname, patronymic, age, gender)

пример {address}:{port}/api/v1/person/?name=Petr

Для пагинации используется параметр limit

пример {address}:{port}/api/v1/person/?limit=5
___
POST {address}:{port}/api/v1/person/

Для вставки нужно отправить заполненный объект

пример

{
"name": "gregeg",
"surname": "ewqeqwe",
"patronymic": "" // необязательно
}
___
PUT {address}:{port}/api/v1/person/

Для обновления нужно отправить изменённый объект

{
"id": 1,
"name": "новое имя",
"surname": "новая фамилия",
"patronymic": "Vasilevich",
"age": 30,
"gender": "male",
"country": [
{
"country_id": "GT",
"probability": 0.1
},
{
"country_id": "DS",
"probability": 0.123
},
{
"country_id": "SK",
"probability": 0.055
},
{
"country_id": "UA",
"probability": 0.011
},
{
"country_id": "AT",
"probability": 0.005
}
]
}
___
DELETE {address}:{port}/api/v1/person/delete/

Для удаления нужно передать ID 

пример {address}:{port}/api/v1/person/delete/2