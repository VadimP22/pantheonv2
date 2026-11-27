CREATE TABLE IF NOT EXISTS tasks (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    name              TEXT    NOT NULL,
    description       TEXT    NOT NULL DEFAULT '',
    icon_item_name    TEXT    NOT NULL,
    is_checkbox       INTEGER NOT NULL DEFAULT 0 CHECK (is_checkbox IN (0, 1)),
    item_id           TEXT    NOT NULL DEFAULT '',
    item_count        INTEGER NOT NULL DEFAULT 0,
    target_item_count INTEGER NOT NULL DEFAULT 0,
    stockpile_id      INTEGER REFERENCES stockpiles(id) ON DELETE SET NULL
);
