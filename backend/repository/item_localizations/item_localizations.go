// Package item_localizations предоставляет доступ к данным таблицы item_localizations.
package item_localizations

import (
	"database/sql"
	"fmt"

	model "backend/model/item_localizations"
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
// Запись данных выполняется через миграции, поэтому методов изменения нет.
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
