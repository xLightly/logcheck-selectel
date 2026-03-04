# logcheck — Go-линтер для проверки соглашений о логировании

Go-линтер, совместимый с [golangci-lint](https://golangci-lint.run/), который анализирует сообщения в логах и проверяет их на соответствие установленным соглашениям.

## Поддерживаемые логгеры

- `log/slog` (стандартная библиотека Go)
- `go.uber.org/zap`

## Правила

| № | Правило | Автоисправление |
|---|---------|----------------|
| 1 | Сообщения логов должны начинаться со строчной буквы | ✅ |
| 2 | Сообщения логов должны быть только на английском (без не-ASCII символов) | — |
| 3 | Сообщения логов не должны содержать специальные символы или эмодзи | ✅ |
| 4 | Сообщения логов не должны содержать ключевые слова, связанные с чувствительными данными | — |

### Примеры

```go
// ❌ Плохо
slog.Info("Starting server on port 8080")
slog.Error("ошибка подключения к базе данных")
slog.Info("server started! 🚀")
slog.Info("user password: " + password)

// ✅ Хорошо
slog.Info("starting server on port 8080")
slog.Error("failed to connect to database")
slog.Info("server started")
slog.Info("user authenticated successfully")
```

## Установка

### Как отдельная утилита

```bash
go install github.com/logcheck/cmd/logcheck@latest
```

Запуск:

```bash
logcheck ./...
```

### Как плагин для golangci-lint

1. Соберите плагин (должна использоваться та же версия Go, что и в вашем golangci-lint):

```bash
go build -buildmode=plugin -o logcheck.so ./plugin/
```

2. Добавьте в `.golangci.yml`:

```yaml
linters-settings:
  custom:
    logcheck:
      path: ./logcheck.so
      description: "Checks log messages for convention compliance"
      original-url: github.com/logcheck
      settings:
        config: ".logcheck.yml"

linters:
  enable:
    - logcheck
```

3. Запустите:

```bash
golangci-lint run ./...
```

## Конфигурация

Создайте файл `.logcheck.yml` в корне проекта:

```yaml
rules:
  lowercase_start: true
  english_only: true
  no_special_chars: true
  no_sensitive: true

# Пользовательские паттерны для чувствительных данных
sensitive_patterns:
  - "my_internal_token"
  - "secret_field"

# Разрешённые специальные символы (если нужно)
allowed_special_chars: ""
```

Передача пути к конфигу через флаг CLI:

```bash
logcheck -config .logcheck.yml ./...
```

## Запуск тестов

```bash
# Юнит-тесты для отдельных правил
go test -v ./pkg/rules/...

# Интеграционные тесты анализатора (использует testdata/)
go test -v ./pkg/analyzer/...

# Все тесты проекта
go test -v ./...
```

## Структура проекта

```
logcheck/
├── cmd/logcheck/          # Точка входа CLI (standalone)
│   └── main.go
├── pkg/
│   ├── analyzer/          # Ядро анализатора + конфигурация
│   │   ├── analyzer.go
│   │   ├── analyzer_test.go
│   │   └── config.go
│   └── rules/             # Реализации правил
│       ├── lowercase.go
│       ├── english.go
│       ├── special_chars.go
│       ├── sensitive.go
│       └── rules_test.go
├── plugin/                # Точка входа для плагина golangci-lint
│   └── plugin.go
├── testdata/src/          # Тестовые фикстуры для analysistest
│   ├── a/                 # Тесты правила lowercase
│   ├── b/                 # Тесты правила english-only
│   ├── c/                 # Тесты правила special_chars
│   └── d/                 # Тесты правила sensitive
├── .logcheck.yml          # Пример конфигурации
├── .golangci.yml          # Пример конфига golangci-lint
├── .github/workflows/     # GitHub Actions CI
├── .gitlab-ci.yml         # GitLab CI
├── go.mod
└── README.md
```

## Как это работает

Линтер использует фреймворк `golang.org/x/tools/go/analysis` для обхода AST исходных файлов Go:

1. Для каждого вызова функции проверяется, является ли он методом логгера из поддерживаемого пакета
2. Извлекается первый строковый аргумент (сообщение лога)
3. К сообщению применяются включённые правила валидации
4. При возможности генерируются `SuggestedFix` для автоисправления

> Выражения конкатенации строк вида `"prefix: " + variable` анализируются частично — проверяется только литеральная часть.

## Требования

- Go **1.22+**
- golangci-lint **v1.57+** (для поддержки плагинов)

## Лицензия

MIT
