-- Поскольку мы не знаем, какое значение было до удаления,
-- восстанавливаем колонку со значением по умолчанию,
-- которое было при создании таблицы.
ALTER TABLE bots ADD COLUMN auto_publish BOOLEAN NOT NULL DEFAULT FALSE;