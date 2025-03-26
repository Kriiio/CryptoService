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
+ Мониторинг с помощью prometheus
+ Трейсинг с помощью OpenTelemetry

## Запуск приложения

    __Опции makefile__
    1. Main
        + _run_ - для запуска приложения
        + _build_ - для сброки приложения
        + _test_ - для запуска unit-тустов
    2. Tools
        + _lint_ - запуск линтера golanchi-lint
    3. Docker commands
        + _docker-build_ - для сборки приложения
        + _docker-up_ - для запуска контейнеров

