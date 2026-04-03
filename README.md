# GophKeeper Client - Руководство пользователя

## Содержание
1. [Введение](#введение)
2. [Запуск сервера](#запуск-сервера)
3. [Установка](#установка)
4. [Начало работы](#начало-работы)
5. [Аутентификация](#аутентификация)
6. [Управление секретами](#управление-секретами)
7. [Примеры использования](#примеры-использования)
8. [Команды справки](#команды-справки)

---

## Введение

**GophKeeper** — это безопасный CLI-клиент для управления конфиденциальной информацией. Приложение позволяет:

- Регистрироваться и авторизироваться на сервере
- Создавать, редактировать, удалять и просматривать секреты
- Хранить пароли, текстовые данные, данные банковских карт и файлы
- Синхронизировать данные между несколькими устройствами
- Безопасно передавать данные через защищённое TLS-соединение

### Поддерживаемые платформы

- Linux (x86_64, ARM64)
- macOS (Apple Silicon, Intel)
- Windows (x86_64)

---

## Запуск сервера

### Установите docker и docker compose

Выполните в корне проекта

```base
docker-compose up
```

## Установка

### 1. Скачайте клиент

Загрузите бинарный файл для вашей платформы с GitHub Releases:

```bash
# Linux
wget https://github.com/flash1nho/GophKeeper/tree/main/releases/download/v1.0.0/gophkeeper-client-linux-amd64
chmod +x gophkeeper-client-linux-amd64

# macOS (Apple Silicon)
wget https://github.com/flash1nho/GophKeeper/tree/main/releases/download/v1.0.0/gophkeeper-client-darwin-arm64
chmod +x gophkeeper-client-darwin-arm64

# Windows
https://github.com/flash1nho/GophKeeper/tree/main/releases/download/v1.0.0/gophkeeper-client-windows-amd64.exe
```

### 2. Добавьте в PATH

```bash
# Linux/macOS
mv gophkeeper-client-linux-amd64 /usr/local/bin/gophkeeper
export PATH=$PATH:/usr/local/bin

# Windows
# Добавьте папку с .exe в переменную окружения PATH
```

### 3. Проверьте установку

```bash
gophkeeper --version
```

---

## Начало работы

### 1. Инициализация конфигурации

При первом запуске создайте конфигурационный файл `.env`:

```bash
# token авторизации создается автоматически в ~/.gophkeeper.yaml
gophkeeper users register --login my_login --password my_password --secret my_secret_word
```

### 2. Переменные окружения

Создайте файл `config/.env`:

```bash
# .env
DATABASE_DSN=postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable
GRPC_SERVER_ADDRESS=localhost:3200
MASTER_KEY=SUPER_SECRET_MASTER_KEY
```

---

## Аутентификация

### Регистрация нового пользователя

```bash
gophkeeper users register \
  --login your_login \
  --password your_password \
  --secret your_secret_word
```

**Параметры:**
- `--login` (обязательно) — уникальный логин пользователя
- `--password` (обязательно) — пароль
- `--secret` (обязательно) — секретное слово для восстановления доступа

**Пример:**
```bash
gophkeeper users register --login john_doe --password MySecurePass123 --secret MySecretWord42
```

**Результат:**
```
✅ Успешная регистрация!
```

Токен сохраняется в `~/.gophkeeper.yaml`

### Вход в систему

```bash
gophkeeper users login \
  --login your_login \
  --password your_password
```

**Параметры:**
- `--login` (обязательно) — ваш логин
- `--password` (обязательно) — ваш пароль

**Пример:**
```bash
gophkeeper users login --login john_doe --password MySecurePass123
```

**Результат:**
```
✅ Успешный вход!
```

---

## Управление секретами

### Общая структура команд

```bash
gophkeeper secrets <тип> <действие> [параметры]
```

### Типы секретов

1. **text** — произвольные текстовые данные
2. **cred** — пара логин/пароль
3. **card** — данные банковской карты
4. **file** — бинарные файлы

### Действия

Для каждого типа доступны следующие команды:

| Команда | Описание |
|---------|---------|
| `create` | Создать новый секрет |
| `get` | Просмотреть секрет по ID |
| `list` | Показать список всех секретов |
| `update` | Обновить существующий секрет |
| `delete` | Удалить секрет |
| `upload` | Загрузить файл |
| `download` | Скачать файл |

---

## Примеры использования

### Текстовые данные (text)

#### Создание

```bash
gophkeeper secrets text create \
  --content "Текст заметки"
```

**Параметры:**
- `--content` — содержимое текстовых данных

**Пример:**
```bash
gophkeeper secrets text create --content "My important notes and thoughts"
```

#### Просмотр списка

```bash
gophkeeper secrets text list
```

**Вывод:**
```
Список секретов:
---
id: 1
content: My important notes and thoughts
type: Text
created_at: 2026-04-02T10:30:00Z
updated_at: 2026-04-02T10:30:00Z
---
```

#### Получение по ID

```bash
gophkeeper secrets text get --id 1
```

#### Обновление

```bash
gophkeeper secrets text update \
  --id 1 \
  --content "Updated notes"
```

#### Удаление

```bash
gophkeeper secrets text delete --id 1
```

---

### Пары логин/пароль (cred)

#### Создание

```bash
gophkeeper secrets cred create \
  --login my_email@example.com \
  --password MySecurePassword123
```

**Параметры:**
- `--login` — логин или email
- `--password` — пароль

**Пример:**
```bash
gophkeeper secrets cred create \
  --login john@example.com \
  --password SecurePass2024
```

#### Просмотр списка

```bash
gophkeeper secrets cred list
```

#### Получение по ID

```bash
gophkeeper secrets cred get --id 1
```

#### Обновление

```bash
gophkeeper secrets cred update \
  --id 1 \
  --login john@example.com \
  --password NewPassword456
```

#### Удаление

```bash
gophkeeper secrets cred delete --id 1
```

---

### Данные банковских карт (card)

#### Создание

```bash
gophkeeper secrets card create \
  --number 4111111111111111 \
  --holder "JOHN DOE" \
  --expiry "12/26" \
  --cvv 123
```

**Параметры:**
- `--number` — номер карты
- `--holder` — имя держателя
- `--expiry` — дата истечения (MM/YY)
- `--cvv` — код безопасности

**Пример:**
```bash
gophkeeper secrets card create \
  --number 4532015112830366 \
  --holder "IVAN PETROV" \
  --expiry "08/25" \
  --cvv 456
```

#### Просмотр списка

```bash
gophkeeper secrets card list
```

#### Получение по ID

```bash
gophkeeper secrets card get --id 2
```

#### Обновление

```bash
gophkeeper secrets card update \
  --id 2 \
  --number 4532015112830366 \
  --holder "IVAN PETROV" \
  --expiry "08/27"
```

#### Удаление

```bash
gophkeeper secrets card delete --id 2
```

---

### Загрузка файлов (file upload)

#### Загрузка нового файла

```bash
gophkeeper secrets file upload --path /path/to/your/file.pdf
```

**Параметры:**
- `--path` (обязательно) — полный путь к файлу

**Особенности:**
- Поддерживается возобновление загрузки при прерывании
- Максимальный размер файла зависит от сервера
- Загружаемые файлы шифруются на сервере

**Пример:**
```bash
gophkeeper secrets file upload --path ~/Documents/passport.pdf
```

**Вывод:**
```
🚀 Начинаю загрузку: passport.pdf (всего 2048576 байт)
📤 Загрузка: 50.00% [1024288 / 2048576 байт]
✅ Файл загружен
---
id: 3
file_name: passport.pdf
type: File
created_at: 2026-04-02T11:00:00Z
updated_at: 2026-04-02T11:00:00Z
---
```

#### Список загруженных файлов

```bash
gophkeeper secrets file list
```

#### Загрузка файла с возобновлением

Если соединение разорвалось во время загрузки, просто повторите команду:

```bash
gophkeeper secrets file upload --path ~/Documents/large_file.zip
```

Клиент автоматически:
1. Обнаружит существующий локальный файл
2. Запросит статус загрузки на сервере
3. Продолжит с того же места

**Вывод:**
```
🔄 Найден локальный фрагмент (1024288 байт), запрашиваю докачку...
🚀 Начинаю загрузку...
📥 Прогресс: 75.50% [1546240 / 2048576 байт]
```

#### Скачивание файла

```bash
gophkeeper secrets file download \
  --id 3 \
  --out ~/Downloads/passport.pdf
```

**Параметры:**
- `--id` (обязательно) — ID файла
- `--out` (обязательно) — путь для сохранения

**Особенности:**
- Поддерживается докачка при прерывании
- Файл расшифровывается при скачивании

**Пример:**
```bash
gophkeeper secrets file download --id 3 --out ~/Documents/my_file.pdf
```

**Вывод:**
```
🔄 Найден локальный фрагмент (1024288 байт), запрашиваю докачку...
🚀 Начинаю загрузку...
📥 Прогресс: 100.00% (2048576 / 2048576 байт)
✅ Файл успешно скачан
```

#### Просмотр информации о файле

```bash
gophkeeper secrets file get --id 3
```

#### Удаление файла

```bash
gophkeeper secrets file delete --id 3
```

---

## Команды справки

### Получить справку по приложению

```bash
gophkeeper --help
```

**Вывод:**
```
GophKeeper client

Usage:
  gophkeeper [command]

Available Commands:
  completion        Generate the autocompletion script for the specified shell
  help              Help about any command
  secrets           Менеджер хранения данных
  users             Менеджер регистрации и авторизации

Flags:
  -h, --help      help for gophkeeper
  -v, --version   version for gophkeeper

Use "gophkeeper [command] --help" for more information about a command.
```

### Справка по командам users

```bash
gophkeeper users --help
```

### Справка по команде регистрации

```bash
gophkeeper users register --help
```

### Справка по управлению секретами

```bash
gophkeeper secrets --help
gophkeeper secrets text --help
gophkeeper secrets cred --help
gophkeeper secrets card --help
gophkeeper secrets file --help
```

---

### Безопасность

1. **Токен** хранится в зашифрованном виде в `~/.gophkeeper.yaml`
2. **Все соединения** с сервером используют TLS
3. **Данные** шифруются на сервере с использованием master key
4. **Пароли** никогда не передаются в открытом виде

---

## Типичные сценарии

### Сценарий 1: Сохранение пароля от аккаунта

```bash
# 1. Войдите в систему
gophkeeper users login --login john_doe --password MyPassword

# 2. Сохраните пароль
gophkeeper secrets cred create \
  --login john@gmail.com \
  --password secretpassword123

# 3. Просмотрите список
gophkeeper secrets cred list

# 4. Получите детали по ID
gophkeeper secrets cred get --id 1
```

### Сценарий 2: Загрузка и скачивание конфиденциального документа

```bash
# 1. Загрузите документ
gophkeeper secrets file upload --path ~/Documents/contract.pdf

# 2. Список файлов
gophkeeper secrets file list

# 3. На другом устройстве, после входа, скачайте:
gophkeeper secrets file download --id 1 --out ~/Downloads/contract.pdf
```

### Сценарий 3: Синхронизация между устройствами

```bash
# Устройство 1:
gophkeeper users login --login john_doe --password MyPassword
gophkeeper secrets cred list  # Все ваши секреты загружены

# Устройство 2:
gophkeeper users login --login john_doe --password MyPassword
gophkeeper secrets cred list  # Те же секреты!
```

---

## Обработка ошибок

### Ошибка: "Токен не найден"

```
❌ Токен не найден. Сначала выполните вход (login)!
```

**Решение:**
```bash
gophkeeper users login --login your_login --password your_password
```

### Ошибка: "Не удалось подключиться"

```
❌ не удалось подключиться
```

**Причины:**
- Сервер не запущен
- Неверный адрес сервера в `GRPC_SERVER_ADDRESS`
- Проблемы с сетью

**Решение:**
1. Убедитесь, что сервер запущен
2. Проверьте переменную окружения `GRPC_SERVER_ADDRESS`

### Ошибка: "Сертификаты не найдены"

```
❌ сертификаты не найдены
```

**Решение:**
1. Убедитесь, что папка `internal/certs/` находится в корневой директории
2. Файлы `ca.crt`, `client.crt`, `client.key` должны быть на месте

---

## Версия и информация о сборке

```bash
gophkeeper --version
```

**Вывод содержит:**
- Версию клиента
- Дату сборки
- Хеш коммита

Пример:
```
GophKeeper Client v1.0.0
Build Date: 2026-04-02
Commit: 7bde3d81edf
```

---

## Контакты и поддержка

- **GitHub:** https://github.com/flash1nho/GophKeeper
- **Issues:** https://github.com/flash1nho/GophKeeper/issues

---

### Новое в этой версии

- ✅ Полная поддержка текстовых данных
- ✅ Управление парами логин/пароль
- ✅ Хранение данных банковских карт
- ✅ Загрузка и скачивание файлов с возобновлением
- ✅ TLS-шифрование соединений
- ✅ JWT-аутентификация
- ✅ GitHub Actions для CI/CD
