# Cервіс скорочення посилань на Go
Сервіс має два ендпоінти: 
+ Cтворити код для скороченого посилання (POST with json payload {"url":"someURL"}):
```
http://localhost:8033/url
```

+ Отримати повне посилання за кодом:

```
http://localhost:8033/url/{alias}
```

Для збереження даних використовується PostgreSQL. Як бібліотеку для побудови API використано https://github.com/go-chi/chi.

## Пояснення щодо першого запуску

```
# drop any persistent state to make sure you are working with clean install
$ docker compose down -v && docker compose pull
# spin everything up
$ docker compose up -d
# wait while environment initialization is complete
$ docker compose logs -f
```

Після цього має бути отриман доступ до сервісу за адресою http://localhost:8033