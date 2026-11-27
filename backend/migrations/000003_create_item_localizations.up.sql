CREATE TABLE IF NOT EXISTS item_localizations (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    lang      TEXT NOT NULL,
    item_id   TEXT NOT NULL,
    item_name TEXT NOT NULL,
    UNIQUE (lang, item_id)
);
