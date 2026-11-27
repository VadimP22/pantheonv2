CREATE TABLE IF NOT EXISTS seeded_languages (
    lang      TEXT PRIMARY KEY,
    seeded_at TEXT NOT NULL DEFAULT (datetime('now'))
);
