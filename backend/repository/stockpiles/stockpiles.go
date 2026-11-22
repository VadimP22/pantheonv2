// Package stockpiles предоставляет доступ к данным таблицы stockpiles.
package stockpiles

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	model "backend/model/stockpiles"
)

// Имена текстовых колонок таблицы stockpiles, используемые в CRUD-запросах.
const (
	colID        = "id"
	colName      = "name"
	colDimension = "dimension"
	colIconName  = "icon_name"
	colX         = "x"
	colY         = "y"
	colZ         = "z"
)

// listColumns перечисляет все колонки, сканируемые при чтении записи.
var listColumns = []string{
	colID,
	colName,
	colDimension,
	colIconName,
	colX,
	colY,
	colZ,
}

// Repository обеспечивает CRUD-операции над таблицей stockpiles.
type Repository struct {
	db *sql.DB
}

// New возвращает новый Repository, работающий через предоставленное соединение с БД.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create вставляет новую запись и заполняет поле ID переданной структуры.
// Поле ID на входе игнорируется (присваивается БД).
func (r *Repository) Create(s *model.Stockpile) error {
	const query = `INSERT INTO stockpiles (name, dimension, icon_name, x, y, z)
		VALUES (?, ?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query,
		s.Name, s.Dimension, s.IconName, s.X, s.Y, s.Z,
	)
	if err != nil {
		return fmt.Errorf("insert stockpile: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get stockpile last insert id: %w", err)
	}

	s.ID = id
	return nil
}

// GetByID возвращает запись по идентификатору.
// Если запись не найдена, возвращается sql.ErrNoRows — вызов решает, как обработать.
func (r *Repository) GetByID(id int64) (*model.Stockpile, error) {
	query := buildSelect("WHERE " + colID + " = ?")

	var s model.Stockpile
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.Name, &s.Dimension, &s.IconName, &s.X, &s.Y, &s.Z)
	if err != nil {
		return nil, fmt.Errorf("get stockpile by id %d: %w", id, err)
	}

	return &s, nil
}

// List возвращает все записи, упорядоченные по id.
// При пустой таблице возвращается пустой не-nil слайс.
func (r *Repository) List() ([]model.Stockpile, error) {
	query := buildSelect("ORDER BY " + colID)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list stockpiles: %w", err)
	}
	defer rows.Close()

	items := make([]model.Stockpile, 0)
	for rows.Next() {
		var s model.Stockpile
		if err := rows.Scan(&s.ID, &s.Name, &s.Dimension, &s.IconName, &s.X, &s.Y, &s.Z); err != nil {
			return nil, fmt.Errorf("scan stockpile: %w", err)
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stockpiles: %w", err)
	}

	return items, nil
}

// Count возвращает общее количество записей в таблице stockpiles.
func (r *Repository) Count() (int64, error) {
	const query = `SELECT COUNT(*) FROM stockpiles`

	var n int64
	if err := r.db.QueryRow(query).Scan(&n); err != nil {
		return 0, fmt.Errorf("count stockpiles: %w", err)
	}
	return n, nil
}

// ListPaged возвращает страницу записей, упорядоченных по id, начиная с offset,
// в количестве не более limit. При пустой странице возвращается пустой не-nil слайс.
func (r *Repository) ListPaged(limit, offset int) ([]model.Stockpile, error) {
	query := buildSelect("ORDER BY " + colID + " LIMIT ? OFFSET ?")

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list stockpiles page: %w", err)
	}
	defer rows.Close()

	items := make([]model.Stockpile, 0)
	for rows.Next() {
		var s model.Stockpile
		if err := rows.Scan(&s.ID, &s.Name, &s.Dimension, &s.IconName, &s.X, &s.Y, &s.Z); err != nil {
			return nil, fmt.Errorf("scan stockpile: %w", err)
		}
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stockpiles: %w", err)
	}

	return items, nil
}

// ErrNoUpdateFields возвращается, когда в input нет ни одного не-nil поля,
// поэтому обновлять нечего.
var ErrNoUpdateFields = errors.New("stockpiles: no fields to update")

// Update обновляет у записи id только те поля input, которые не nil.
// Возвращает true, если запись существовала и была изменена, false — если записи
// не было либо значения оказались идентичны текущим (SQLite RowsAffected = 0).
// Если в input нет ни одного не-nil поля, возвращается ErrNoUpdateFields.
func (r *Repository) Update(id int64, in model.UpdateStockpileInput) (bool, error) {
	sets := make([]string, 0, 6)
	args := make([]any, 0, 7)

	if in.Name != nil {
		sets = append(sets, colName+" = ?")
		args = append(args, *in.Name)
	}
	if in.Dimension != nil {
		sets = append(sets, colDimension+" = ?")
		args = append(args, *in.Dimension)
	}
	if in.IconName != nil {
		sets = append(sets, colIconName+" = ?")
		args = append(args, *in.IconName)
	}
	if in.X != nil {
		sets = append(sets, colX+" = ?")
		args = append(args, *in.X)
	}
	if in.Y != nil {
		sets = append(sets, colY+" = ?")
		args = append(args, *in.Y)
	}
	if in.Z != nil {
		sets = append(sets, colZ+" = ?")
		args = append(args, *in.Z)
	}

	if len(sets) == 0 {
		return false, ErrNoUpdateFields
	}

	query := "UPDATE stockpiles SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, id)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return false, fmt.Errorf("update stockpile %d: %w", id, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("update stockpile %d rows affected: %w", id, err)
	}

	return n > 0, nil
}

// Delete удаляет запись по идентификатору.
// Возвращает ошибку только при проблемах с БД; отсутствие записи не является ошибкой.
func (r *Repository) Delete(id int64) error {
	const query = `DELETE FROM stockpiles WHERE id = ?`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete stockpile %d: %w", id, err)
	}

	// RowsAffected учитываем лишь для отладки; метод не сигнализирует об отсутствии записи.
	if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("delete stockpile %d rows affected: %w", id, err)
	}

	return nil
}

// buildSelect собирает SELECT-запрос со списком всех колонок и опциональным хвостом
// (WHERE/ORDER BY). Использует listColumns, чтобы списки не расходились с колонками.
func buildSelect(tail string) string {
	const from = "FROM stockpiles"
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
