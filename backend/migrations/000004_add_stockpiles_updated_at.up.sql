-- updated_at хранит время последнего изменения записи в RFC3339 (UTC).
-- SQLite не позволяет ADD COLUMN с неконстантным DEFAULT, поэтому добавляем
-- колонку с пустой строкой и затем заполняем её текущим временем.
ALTER TABLE stockpiles ADD COLUMN updated_at TEXT NOT NULL DEFAULT '';

UPDATE stockpiles SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ', 'now') WHERE updated_at = '';
