// Package tasks предоставляет доступ к данным таблицы tasks.
package tasks

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	model "backend/model/tasks"
)

// Имена колонок таблицы tasks, используемые в запросах на чтение.
const (
	colID              = "id"
	colName            = "name"
	colDescription     = "description"
	colIconItemName    = "icon_item_name"
	colIsCheckbox      = "is_checkbox"
	colDone            = "done"
	colItemID          = "item_id"
	colItemCount       = "item_count"
	colTargetItemCount = "target_item_count"
	colStockpileID     = "stockpile_id"
)

// listColumns перечисляет все колонки, сканируемые при чтении записи.
// Порядок обязан совпадать с порядком выбора в buildSelect.
var listColumns = []string{
	colID,
	colName,
	colDescription,
	colIconItemName,
	colIsCheckbox,
	colDone,
	colItemID,
	colItemCount,
	colTargetItemCount,
	colStockpileID,
}

// Repository обеспечивает операции чтения над таблицей tasks.
type Repository struct {
	db *sql.DB
}

// New возвращает новый Repository, работающий через предоставленное соединение с БД.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create вставляет новую запись в таблицу tasks и возвращает присвоенный id.
// Поле is_checkbox хранится как 0/1 (SQLite INTEGER); драйвер конвертирует bool.
// stockpile_id может быть nil — в БД запишется NULL.
func (r *Repository) Create(name, description, iconItemName string, isCheckbox bool, done bool, itemID string, itemCount, targetItemCount int64, stockpileID *int64) (int64, error) {
	const query = `INSERT INTO tasks (name, description, icon_item_name, is_checkbox, done, item_id, item_count, target_item_count, stockpile_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	var id int64
	if err := r.db.QueryRow(query, name, description, iconItemName, isCheckbox, done, itemID, itemCount, targetItemCount, stockpileID).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}
	return id, nil
}

// GetByID возвращает запись по идентификатору.
// Если запись не найдена, возвращается sql.ErrNoRows — вызов решает, как обработать.
func (r *Repository) GetByID(id int64) (*model.Task, error) {
	query := buildSelect("WHERE " + colID + " = ?")

	t, err := scanTask(r.db.QueryRow(query, id))
	if err != nil {
		return nil, fmt.Errorf("get task by id %d: %w", id, err)
	}

	return &t, nil
}

// List возвращает все записи, упорядоченные по id.
// При пустой таблице возвращается пустой не-nil слайс.
func (r *Repository) List() ([]model.Task, error) {
	query := buildSelect("ORDER BY " + colID)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	items := make([]model.Task, 0)
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return items, nil
}

// Count возвращает общее количество записей в таблице tasks.
func (r *Repository) Count() (int64, error) {
	const query = `SELECT COUNT(*) FROM tasks`

	var n int64
	if err := r.db.QueryRow(query).Scan(&n); err != nil {
		return 0, fmt.Errorf("count tasks: %w", err)
	}
	return n, nil
}

// AddItemCount прибавляет delta к item_count задачи id.
// Возвращает false, если задача не найдена.
func (r *Repository) AddItemCount(id int64, delta int64) (bool, error) {
	res, err := r.db.Exec("UPDATE tasks SET "+colItemCount+" = "+colItemCount+" + ? WHERE "+colID+" = ?", delta, id)
	if err != nil {
		return false, fmt.Errorf("tasks: add item count: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("tasks: add item count rows affected: %w", err)
	}
	return n > 0, nil
}

// ErrNoUpdateFields возвращается, когда в input нет ни одного не-nil поля,
// поэтому обновлять нечего.
var ErrNoUpdateFields = errors.New("tasks: no fields to update")

// Update обновляет у записи id только те поля input, которые не nil.
// Возвращает true, если запись существовала и была изменена, false — если записи
// не было (SQLite RowsAffected = 0).
// Если в input нет ни одного не-nil поля, возвращается ErrNoUpdateFields.
func (r *Repository) Update(id int64, in model.UpdateTaskInput) (bool, error) {
	sets := make([]string, 0, 9)
	args := make([]any, 0, 10)

	if in.Name != nil {
		sets = append(sets, colName+" = ?")
		args = append(args, *in.Name)
	}
	if in.Description != nil {
		sets = append(sets, colDescription+" = ?")
		args = append(args, *in.Description)
	}
	if in.IconItemName != nil {
		sets = append(sets, colIconItemName+" = ?")
		args = append(args, *in.IconItemName)
	}
	if in.IsCheckbox != nil {
		sets = append(sets, colIsCheckbox+" = ?")
		args = append(args, *in.IsCheckbox)
	}
	if in.Done != nil {
		sets = append(sets, colDone+" = ?")
		args = append(args, *in.Done)
	}
	if in.ItemID != nil {
		sets = append(sets, colItemID+" = ?")
		args = append(args, *in.ItemID)
	}
	if in.ItemCount != nil {
		sets = append(sets, colItemCount+" = ?")
		args = append(args, *in.ItemCount)
	}
	if in.TargetItemCount != nil {
		sets = append(sets, colTargetItemCount+" = ?")
		args = append(args, *in.TargetItemCount)
	}
	if in.StockpileID != nil {
		sets = append(sets, colStockpileID+" = ?")
		args = append(args, *in.StockpileID)
	}

	if len(sets) == 0 {
		return false, ErrNoUpdateFields
	}

	query := "UPDATE tasks SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, id)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return false, fmt.Errorf("update task %d: %w", id, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("update task %d rows affected: %w", id, err)
	}

	return n > 0, nil
}

// ListPaged возвращает страницу записей, упорядоченных по id, начиная с offset,
// в количестве не более limit. При пустой странице возвращается пустой не-nil слайс.
func (r *Repository) ListPaged(limit, offset int) ([]model.Task, error) {
	query := buildSelect("ORDER BY " + colID + " LIMIT ? OFFSET ?")

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list tasks page: %w", err)
	}
	defer rows.Close()

	items := make([]model.Task, 0)
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return items, nil
}

// buildSelect собирает SELECT-запрос со списком всех колонок и опциональным
// хвостом (WHERE/ORDER BY). Порядок колонок обязан совпадать с порядком
// в listColumns.
func buildSelect(tail string) string {
	base := "SELECT " + join(listColumns) + " FROM tasks"
	if tail == "" {
		return base
	}
	return base + " " + tail
}

// join объединяет элементы списка через запятую.
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

// scanTask читает одну строку выборки tasks. Признак is_checkbox хранится
// в БД как 0/1 и сканируется в bool; stockpile_id может быть NULL
// (целевой склад не задан) и читается через sql.NullInt64.
func scanTask(scanner interface{ Scan(dest ...any) error }) (model.Task, error) {
	var (
		t       model.Task
		stockID sql.NullInt64
	)
	if err := scanner.Scan(
		&t.ID, &t.Name, &t.Description, &t.IconItemName, &t.IsCheckbox, &t.Done,
		&t.ItemID, &t.ItemCount, &t.TargetItemCount, &stockID,
	); err != nil {
		return model.Task{}, fmt.Errorf("scan task: %w", err)
	}
	if stockID.Valid {
		v := stockID.Int64
		t.StockpileID = &v
	}
	return t, nil
}
