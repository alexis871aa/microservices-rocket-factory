# microservices-rocket-factory

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/alexis871aa/4cbdede79d6d29090e053a7c93192093/raw/coverage.json)

Этот репозиторий содержит проект из курса [Микросервисы, как в BigTech 2.0](https://olezhek28.courses/microservices).

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
    - Линтинг кода
    - Проверка безопасности
    - Выполняется автоматическое извлечение версий из Taskfile.yml
