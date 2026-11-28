// Package item_registry предоставляет доступ к реестру предметов Minecraft
// (таблица item_localizations). Запись данных выполняется через langseed,
// поэтому репозиторий read-only.
package item_registry

import (
	"database/sql"
	"fmt"

	model "backend/model/item_registry"
)

// Имена колонок таблицы item_localizations, используемые в запросах.
const (
	colID       = "id"
	colLang     = "lang"
	colItemID   = "item_id"
	colItemName = "item_name"
)

// listColumns перечисляет все колонки, сканируемые при чтении записи.
var listColumns = []string{
	colID,
	colLang,
	colItemID,
	colItemName,
}

// Repository обеспечивает операции чтения над таблицей item_localizations.
// Запись данных выполняется через langseed, поэтому методов изменения нет.
type Repository struct {
	db *sql.DB
}

// New возвращает новый Repository, работающий через предоставленное соединение с БД.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetByItemID возвращает локализацию по паре «язык + предмет».
// Если запись не найдена, возвращается ошибка, оборачивающая sql.ErrNoRows —
// вызов решает, как её обработать.
func (r *Repository) GetByItemID(lang, itemID string) (*model.ItemLocalization, error) {
	query := buildSelect("WHERE " + colLang + " = ? AND " + colItemID + " = ?")

	var rec model.ItemLocalization
	err := r.db.QueryRow(query, lang, itemID).Scan(&rec.ID, &rec.Lang, &rec.ItemID, &rec.ItemName)
	if err != nil {
		return nil, fmt.Errorf("get item localization lang=%q item=%q: %w", lang, itemID, err)
	}

	return &rec, nil
}

// ExistsByItemID проверяет, что предмет с указанным id присутствует
// в реестре (на любом языке).
func (r *Repository) ExistsByItemID(itemID string) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM item_localizations WHERE item_id = ?)`

	var exists bool
	if err := r.db.QueryRow(query, itemID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check item exists item=%q: %w", itemID, err)
	}
	return exists, nil
}

// buildSelect собирает SELECT-запрос со списком всех колонок и опциональным хвостом
// (WHERE/ORDER BY). Использует listColumns, чтобы списки не расходились с колонками.
func buildSelect(tail string) string {
	const from = "FROM item_localizations"
	if tail == "" {
		return "SELECT " + join(listColumns) + " " + from
	}
	return "SELECT " + join(listColumns) + " " + from + " " + tail
}

// join объединяет имена колонок через запятую.
func join(cols []string) string {
	s := ""
	for i, c := range cols {
		if i > 0 {
			s += ", "
		}
		s += c
	}
	return s
}
