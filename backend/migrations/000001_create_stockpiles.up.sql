CREATE TABLE IF NOT EXISTS stockpiles (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT    NOT NULL,
    dimension  TEXT    NOT NULL,
    icon_name  TEXT    NOT NULL,
    x          REAL    NOT NULL,
    y          REAL    NOT NULL,
    z          REAL    NOT NULL
);
