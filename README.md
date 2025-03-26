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

Для запуска приложения создан Makefile 

__Опции makefile__
1. Main
    + `run` - для запуска приложения
        ```make run```
    + `build` - для сброки приложения
        ```make build```
   + `test` - для запуска unit-тустов
        ```make test```
2. Tools
    + `lint` - запуск линтера golanchi-lint
        ```make limit```
3. Docker commands
    + `docker-build` - для сборки приложения
        ```make docker-build```
    + `docker-up` - для запуска контейнеров
        ```make docker-up```

