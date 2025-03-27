# CryptoService

Микросервис представляет собой GRPC сервер 

## Функционал сервиса:
+ Получение из стороннего API криптобиржи(_Kraken_) информации пары **USDT-USD**
    Так как биржа _Garantex_ из ТЗ больше недоступна, а остальные биржи не предоставляют информацию по паре **USDT-RUB**, было принято решение о курсе к _Доллару_
+ Сохранение в базу данных _PostgreSQL_ данных (**asks и bids**) с фиксацией времени 

    **таблица**
    | id  |   time   | ask_price | ask_volume | ask_time | bid_price | bid_volume | bid_time |
    |:---:|:--------:|:---------:|:----------:|:--------:|:---------:|:----------:|:--------:|
    | int |Timestamp |   Float   |    Float   |Timestamp |   Float   |   Float    |Timestamp |

+ Отправка клиенту ask и bid цены с меткой времени
+ Graceful shutdown для мягкого завершения сервиса
+ HealthCheck метод `Ping` для проверки работоспособности сервера(отправляет ответ `Pong`)
+ Мониторинг с помощью prometheus
+ Трейсинг с помощью OpenTelemetry

## Запуск приложения

Для запуска приложения создан Makefile 

__Опции makefile__
1. Main
    + `run` - для запуска приложения
    
        ```bash
        make run
        ```
    + `build` - для сброки приложения

        ```bash
        make build
        ```
   + `test` - для запуска unit-тустов

        ```bash
        make test
        ```
2. Tools
    + `lint` - запуск линтера golanchi-lint

        ```bash
        make limit
        ```
3. Docker commands
    + `docker-build` - для сборки приложения

        ```bash
        make docker-build
        ```
    + `docker-up` - для запуска контейнеров

        ```bash
        make docker-up
        ```

