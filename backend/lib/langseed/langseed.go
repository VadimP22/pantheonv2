// Package langseed заполняет таблицу item_localizations названиями
// ванильных предметов Minecraft из JSON-файлов, встроенных в бинарник.
package langseed

import (
	"database/sql"
	_ "embed" // для директив //go:embed
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// ruJSON — русские локализации (ключи вида "item.minecraft.stick"),
// встроенные в бинарник: внешние файлы при деплое не нужны.
//
//go:embed data/ru.json
var ruJSON []byte

// itemKeyRe выделяет ID из ключей вида "item.minecraft.stick"
// или "block.minecraft.stone" (namespace только minecraft — другие
// префиксы файла, например "item.modifiers.*", относятся к UI-строкам).
var itemKeyRe = regexp.MustCompile(`^(?:item|block)\.minecraft\.([a-z0-9_/]+)$`)

// Run засеивает item_localizations для языков, ещё не отмеченных
// в seeded_languages.
func Run(db *sql.DB) error {
	seeded, err := isSeeded(db, "ru")
	if err != nil {
		return err
	}
	if seeded {
		log.Printf("seed ru: already seeded, skip")
		return nil
	}
	return seedRU(db)
}

// seedRU парсит ru.json и вставляет записи одним транзакционным пакетом,
// в конце отмечая язык в seeded_languages. Флаг ставится в той же
// транзакции: либо записи добавлены и язык отмечен, либо ничего.
func seedRU(db *sql.DB) error {
	var kv map[string]string
	if err := json.Unmarshal(ruJSON, &kv); err != nil {
		return fmt.Errorf("parse embedded ru.json: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin seed tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"INSERT OR IGNORE INTO item_localizations (lang, item_id, item_name) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare seed insert: %w", err)
	}
	defer stmt.Close()

	total, inserted := 0, 0
	// Сначала вставляем ключи "item.*", затем "block.*": у нескольких
	// предметов (например, wheat) есть оба ключа с разными названиями,
	// а INSERT OR IGNORE сохраняет первое попавшееся. Для предмета
	// склада каноничное имя — именно из ключа item.*.
	for _, prefix := range []string{"item.", "block."} {
		for key, name := range kv {
			if !strings.HasPrefix(key, prefix) {
				continue
			}
			m := itemKeyRe.FindStringSubmatch(key)
			if m == nil || name == "" {
				continue
			}
			total++
			itemID := "minecraft:" + m[1] // "minecraft:stone"
			res, err := stmt.Exec("ru", itemID, name)
			if err != nil {
				return fmt.Errorf("seed insert lang=ru item=%q: %w", itemID, err)
			}
			if n, _ := res.RowsAffected(); n > 0 {
				inserted++
			}
		}
	}

	if _, err := tx.Exec(
		"INSERT OR IGNORE INTO seeded_languages (lang) VALUES ('ru')"); err != nil {
		return fmt.Errorf("mark seeded lang ru: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed tx: %w", err)
	}

	log.Printf("seed ru: %d записей, вставлено новых %d", total, inserted)
	return nil
}

// isSeeded проверяет, отмечен ли язык как засеянный в seeded_languages.
func isSeeded(db *sql.DB, lang string) (bool, error) {
	var one int
	err := db.QueryRow(
		"SELECT 1 FROM seeded_languages WHERE lang = ?", lang).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check seeded lang %q: %w", lang, err)
	}
	return true, nil
}
