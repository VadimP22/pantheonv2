CREATE TABLE IF NOT EXISTS stockpile_resources (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    stockpile_id INTEGER NOT NULL REFERENCES stockpiles(id) ON DELETE CASCADE,
    item_id      TEXT    NOT NULL,
    item_count   INTEGER NOT NULL DEFAULT 0
);
