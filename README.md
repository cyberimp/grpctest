# Серверная часть

Сервер реализует CRUD через интерфейс GRPC. Методы передают и получают произвольную структуру.
Сервер хранит данные в Postgres.
Сервер собирает статистику по запросам с момента запуска. Статистика не хранится в базе, а отдаётся по HTTP-хендлеру (один или несколько произвольных, обычный REST, не GRPC).
Клиент можно сделать консольный или как удобнее.

И клиент и сервер на Go.

Документация, тесты, логирование и миграции на усмотрение.