-- +goose Up
-- создаёт расширение uuid-ossp для работы с UUID, not exists не падает, если расширение уже есть
create extension if not exists "uuid-ossp";

-- +goose Down
-- удаляет расширение при откате, если оно было установлено
drop extension if exists "uuid-ossp";
