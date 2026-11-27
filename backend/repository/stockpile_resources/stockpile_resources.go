// Package stockpile_resources предоставляет доступ к данным таблицы stockpile_resources.
package stockpile_resources

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	model "backend/model/stockpile_resources"
)

// Имена текстовых колонок таблицы stockpile_resources, используемые в CRUD-запросах.
const (
	colID          = "id"
	colStockpileID = "stockpile_id"
	colItemID      = "item_id"
	colItemCount   = "item_count"
)

// listColumns перечисляет все колонки, сканируемые при чтении записи.
var listColumns = []string{
	colID,
	colStockpileID,
	colItemID,
	colItemCount,
}

// Repository обеспечивает CRUD-операции над таблицей stockpile_resources.
type Repository struct {
	db *sql.DB
}

// New возвращает новый Repository, работающий через предоставленное соединение с БД.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create вставляет новую запись и заполняет поле ID переданной структуры.
// Поле ID на входе игнорируется (присваивается БД).
func (r *Repository) Create(rec *model.StockpileResource) error {
	const query = `INSERT INTO stockpile_resources (stockpile_id, item_id, item_count)
		VALUES (?, ?, ?)`

	res, err := r.db.Exec(query, rec.StockpileID, rec.ItemID, rec.ItemCount)
	if err != nil {
		return fmt.Errorf("insert stockpile resource: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get stockpile resource last insert id: %w", err)
	}

	rec.ID = id
	return nil
}

// GetByID возвращает запись по идентификатору.
// Если запись не найдена, возвращается sql.ErrNoRows — вызов решает, как обработать.
func (r *Repository) GetByID(id int64) (*model.StockpileResource, error) {
	query := buildSelect("WHERE " + colID + " = ?")

	var rec model.StockpileResource
	err := r.db.QueryRow(query, id).Scan(&rec.ID, &rec.StockpileID, &rec.ItemID, &rec.ItemCount)
	if err != nil {
		return nil, fmt.Errorf("get stockpile resource by id %d: %w", id, err)
	}

	return &rec, nil
}

// GetByStockpileAndItem возвращает ресурс по паре (stockpile_id, item_id).
// Если пара не найдена, возвращается sql.ErrNoRows — вызов решает, как обработать.
func (r *Repository) GetByStockpileAndItem(stockpileID int64, itemID string) (*model.StockpileResource, error) {
	query := buildSelect("WHERE " + colStockpileID + " = ? AND " + colItemID + " = ?")

	var rec model.StockpileResource
	err := r.db.QueryRow(query, stockpileID, itemID).Scan(&rec.ID, &rec.StockpileID, &rec.ItemID, &rec.ItemCount)
	if err != nil {
		return nil, fmt.Errorf("get stockpile resource stockpile=%d item=%q: %w", stockpileID, itemID, err)
	}

	return &rec, nil
}

// ListByStockpile возвращает все ресурсы указанного склада, упорядоченные по id.
// При отсутствии записей возвращается пустой не-nil слайс.
func (r *Repository) ListByStockpile(stockpileID int64) ([]model.StockpileResource, error) {
	query := buildSelect("WHERE " + colStockpileID + " = ? ORDER BY " + colID)

	rows, err := r.db.Query(query, stockpileID)
	if err != nil {
		return nil, fmt.Errorf("list stockpile resources: %w", err)
	}
	defer rows.Close()

	items := make([]model.StockpileResource, 0)
	for rows.Next() {
		var rec model.StockpileResource
		if err := rows.Scan(&rec.ID, &rec.StockpileID, &rec.ItemID, &rec.ItemCount); err != nil {
			return nil, fmt.Errorf("scan stockpile resource: %w", err)
		}
		items = append(items, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stockpile resources: %w", err)
	}

	return items, nil
}

// ErrNoUpdateFields возвращается, когда в input нет ни одного не-nil поля,
// поэтому обновлять нечего.
var ErrNoUpdateFields = errors.New("stockpile_resources: no fields to update")

// Update обновляет у записи id только те поля input, которые не nil.
// Возвращает true, если запись существовала и была изменена, false — если записи
// не было либо значения оказались идентичны текущим (SQLite RowsAffected = 0).
// Если в input нет ни одного не-nil поля, возвращается ErrNoUpdateFields.
func (r *Repository) Update(id int64, in model.UpdateStockpileResourceInput) (bool, error) {
	sets := make([]string, 0, 3)
	args := make([]any, 0, 4)

	if in.StockpileID != nil {
		sets = append(sets, colStockpileID+" = ?")
		args = append(args, *in.StockpileID)
	}
	if in.ItemID != nil {
		sets = append(sets, colItemID+" = ?")
		args = append(args, *in.ItemID)
	}
	if in.ItemCount != nil {
		sets = append(sets, colItemCount+" = ?")
		args = append(args, *in.ItemCount)
	}

	if len(sets) == 0 {
		return false, ErrNoUpdateFields
	}

	query := "UPDATE stockpile_resources SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, id)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return false, fmt.Errorf("update stockpile resource %d: %w", id, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("update stockpile resource %d rows affected: %w", id, err)
	}

	return n > 0, nil
}

// Delete удаляет запись по идентификатору.
// Возвращает ошибку только при проблемах с БД; отсутствие записи не является ошибкой.
func (r *Repository) Delete(id int64) error {
	const query = `DELETE FROM stockpile_resources WHERE id = ?`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete stockpile resource %d: %w", id, err)
	}

	// RowsAffected учитываем лишь для отладки; метод не сигнализирует об отсутствии записи.
	if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("delete stockpile resource %d rows affected: %w", id, err)
	}

	return nil
}

// buildSelect собирает SELECT-запрос со списком всех колонок и опциональным хвостом
// (WHERE/ORDER BY). Использует listColumns, чтобы списки не расходились с колонками.
func buildSelect(tail string) string {
	const from = "FROM stockpile_resources"
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
