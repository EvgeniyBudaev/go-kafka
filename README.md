Разворачивание локального кластера для разработки
Требуется запускать zookeeper раньше Kafka. Т.е. сперва:
```
docker-compose up zookeeper
```
потом:
```
docker-compose up kafka-ui kafka-1 kafka-2 kafka-3
```
Для подключения нужно знать адреса всех брокеров. Для данного файла это:
 * 127.0.0.1:9095
 * 127.0.0.1:9096
 * 127.0.0.1:9097

Инициализация зависимостей
```
go mod init github.com/EvgeniyBudaev/kafka-go/app
```

kafka-go
https://github.com/segmentio/kafka-go
```
go get github.com/segmentio/kafka-go
```

Fiber
https://github.com/gofiber/fiber
```
go get -u github.com/gofiber/fiber/v2
```

Создание docker network
```
docker network create web-network
```