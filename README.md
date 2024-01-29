# Cервіс скорочення посилань на Go
Сервіс має два ендпоінти: створити код для скороченого посилання, отримати повне посилання за кодом. Для збереження даних використовується PostgreSQL. Як бібліотеку для побудови API використано https://github.com/go-chi/chi.

## Пояснення щодо запуску
Щоб коректно запустити програму необхідно у кореневій папці проєкту створити файл збереження локальних змінних середовища із назвою "local.env".
У самому файлі необхідно вказати наступні змінні:

'''
CONFIG_PATH=./config // шлях до каталогу, що містить файл конфігурації
CONFIG_NAME=local // назва файлу конфігурації
POSTGRES_PASSWORD=password // пароль для користувача postgres для з'єднання з базою даних
'''