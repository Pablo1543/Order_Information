<div align="center">

# Order Information System

**Микросервис для просмотра информацией о заказах с REST API**  
**Разработан на Go (Golang)**  

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-28.5.1-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Gin](https://img.shields.io/badge/Gin%20Framework-1.9.1-009688?style=for-the-badge&logo=go&logoColor=white)](https://gin-gonic.com/)
[![NATS](https://img.shields.io/badge/NATS%20Streaming-2.10+-27AE60?style=for-the-badge)](https://docs.nats.io/nats-streaming-concepts/intro)
[![Docker Compose](https://img.shields.io/badge/Docker%20Compose-2.40.0-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docs.docker.com/compose/)

• [Особенности](#особенности) • [Установка](#локальная-установка) • [Запуск через контейнер](#запуск-через-контейнер) • [Добавление заказа](#добавление-нового-заказа) • [Проверки](#проверить-можно-следующими-способами)

</div>

## Демонстрация работы

<div align="center">

### 🎥 Видео-демонстрация работы Web-сайта

[![Demo Video](https://img.shields.io/badge/🎥_Смотреть_видео-FF0000?style=for-the-badge&logo=youtube&logoColor=white)](https://drive.google.com/file/d/1WGPNU-y9FjZeEahE84x1qdbCBcarm1i3/view?usp=sharing)

*Демонстрация работы сайта - Google Drive*

</div>

## Особенности

### Архитектура
- **RESTful API** с Gin Framework
- **Чистая архитектура** с разделением слоев
- **Подключение к PostgreSQL** с миграциями
- **Валидация данных** и обработка ошибок

### Технологии
| Компонент | Технология | Назначение |
|-----------|------------|------------|
| **Backend** | Go 1.21+ | Высокопроизводительный язык |
| **Web Framework** | Gin Gonic | Быстрый HTTP фреймворк |
| **Database** | PostgreSQL | Надежное хранение данных |
| **ORM** | pgx + чисты SQL | Эффективный доступ к данным |

### Производительность
- **Высокая производительность** благодаря Go
- **Минимальное потребление памяти**
- **Быстрое время ответа** API

### Предварительные требования
- Go 1.21+
- PostgreSQL 16+
- Docker (опционально)

## Локальная установка

```bash
# Клонируем репозиторий
git clone https://github.com/Pablo1543/Order_Information.git
cd Order_Information

# Устанавливаем зависимости
go mod download

# Настройка базы данных
# Создайте базу данных PostgreSQL и настройте connection string

# Запускаем приложение
go run cmd/api/main.go
```

## Запуск через контейнер

```bash
# Запуск с Docker Compose
docker-compose up --build
```

## Добавление нового заказа

```bash
# Запустить новый терминал и перейти в каталог
cd ./Order_Information

# Создать JSON файл с заказом
cat > new_order.json << 'EOF'
{
  "order_uid": "test_new_order",
  "track_number": "WBILMTESTNEW",
  "entry": "WBIL",
  "delivery": {
    "name": "John Doe",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "test_transaction_new",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637900000,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 99349340,
      "track_number": "WBILMTESTNEW",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
EOF

# Скопировать файл в работающий контейнер app
docker cp new_order.json order_service:/app/

# Запустить publisher из контейнера app
docker exec order_service /app/publisher /app/new_order.json
```

## Проверить можно следующими способами

```bash
# Посмотреть логи app
docker logs order_service
# Проверить в браузере
curl http://localhost:8080/
```
## Остановка работы сайта

```bash
# Остановить контейнер
docker-compose down -v
```
