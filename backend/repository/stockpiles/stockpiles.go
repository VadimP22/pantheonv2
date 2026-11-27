// Package stockpiles предоставляет доступ к данным таблицы stockpiles.
package stockpiles

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	model "backend/model/stockpiles"
)

// Имена колонок таблицы stockpiles и ресурсов, используемые в CRUD-запросах.
const (
	colID        = "id"
	colName      = "name"
	colDimension = "dimension"
	colIconName  = "icon_name"
	colX         = "x"
	colY         = "y"
	colZ         = "z"
	colUpdatedAt = "updated_at"
	// colItemCount и colItemTypeCount — вычисляемые колонки (не столбцы таблицы);
	// используются как алиасы коррелированных подзапросов.
	colItemCount     = "item_count"
	colItemTypeCount = "item_type_count"
	// colItemID — колонка item_id таблицы stockpile_resources.
	colItemID = "item_id"
)

// listColumns перечисляет все колонки (включая вычисляемые), сканируемые при
// чтении записи. Порядок обязан совпадать с порядком выбора в buildSelect.
var listColumns = []string{
	colID,
	colName,
	colDimension,
	colIconName,
	colX,
	colY,
	colZ,
	colUpdatedAt,
	colItemCount,
	colItemTypeCount,
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
	const query = `INSERT INTO stockpiles (name, dimension, icon_name, x, y, z, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now().UTC()

	res, err := r.db.Exec(query,
		s.Name, s.Dimension, s.IconName, s.X, s.Y, s.Z, formatTime(now),
	)
	if err != nil {
		return fmt.Errorf("insert stockpile: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get stockpile last insert id: %w", err)
	}

	s.ID = id
	s.UpdatedAt = now
	return nil
}

// GetByID возвращает запись по идентификатору.
// Если запись не найдена, возвращается sql.ErrNoRows — вызов решает, как обработать.
func (r *Repository) GetByID(id int64) (*model.Stockpile, error) {
	query := buildSelect("WHERE " + colID + " = ?")

	s, err := scanStockpile(r.db.QueryRow(query, id))
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
		s, err := scanStockpile(rows)
		if err != nil {
			return nil, err
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
		s, err := scanStockpile(rows)
		if err != nil {
			return nil, err
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

// Update обновляет у записи id только те поля input, которые не nil,
// и всегда обновляет время последнего изменения.
// Возвращает true, если запись существовала и была изменена, false — если записи
// не было (SQLite RowsAffected = 0).
// Если в input нет ни одного не-nil поля, возвращается ErrNoUpdateFields.
func (r *Repository) Update(id int64, in model.UpdateStockpileInput) (bool, error) {
	sets := make([]string, 0, 7)
	args := make([]any, 0, 8)

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

	// Время последнего изменения обновляется при любом успешном изменении записи.
	sets = append(sets, colUpdatedAt+" = ?")
	args = append(args, formatTime(time.Now()))

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

// Touch обновляет время последнего изменения записи.
// Отсутствие записи не является ошибкой.
func (r *Repository) Touch(id int64) error {
	const query = `UPDATE stockpiles SET ` + colUpdatedAt + ` = ? WHERE id = ?`

	if _, err := r.db.Exec(query, formatTime(time.Now()), id); err != nil {
		return fmt.Errorf("touch stockpile %d: %w", id, err)
	}
	return nil
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
// (WHERE/ORDER BY). Обычные колонки перечисляются по имени, а вычисляемые
// (item_count, item_type_count) — коррелированными подзапросами с теми же
// алиасами. Порядок колонок обязан совпадать с порядком в listColumns.
func buildSelect(tail string) string {
	const from = "FROM stockpiles"
	cols := append([]string{
		colID, colName, colDimension, colIconName, colX, colY, colZ, colUpdatedAt,
	}, itemCountExpr(), itemTypeCountExpr())
	if tail == "" {
		return "SELECT " + join(cols) + " " + from
	}
	return "SELECT " + join(cols) + " " + from + " " + tail
}

// itemCountExpr возвращает коррелированный подзапрос суммарного количества
// единиц предметов запаса. Для запаса без ресурсов даёт 0 (не NULL).
func itemCountExpr() string {
	return "(" +
		"SELECT COALESCE(SUM(r." + colItemCount + "), 0) " +
		"FROM stockpile_resources r WHERE r.stockpile_id = stockpiles.id" +
		") AS " + colItemCount
}

// itemTypeCountExpr возвращает коррелированный подзапрос количества различных
// типов предметов запаса. Для запаса без ресурсов даёт 0 (не NULL).
func itemTypeCountExpr() string {
	return "(" +
		"SELECT COALESCE(COUNT(DISTINCT r." + colItemID + "), 0) " +
		"FROM stockpile_resources r WHERE r.stockpile_id = stockpiles.id" +
		") AS " + colItemTypeCount
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

// scanStockpile читает одну строку выборки stockpiles, включая сырое значение
// updated_at, и разбирает время в RFC3339. Пустое значение трактуется как
// нулевое время (записи до заполнения миграцией).
func scanStockpile(scanner interface{ Scan(dest ...any) error }) (model.Stockpile, error) {
	var (
		s         model.Stockpile
		updatedAt string
	)
	if err := scanner.Scan(&s.ID, &s.Name, &s.Dimension, &s.IconName, &s.X, &s.Y, &s.Z, &updatedAt, &s.ItemCount, &s.ItemTypeCount); err != nil {
		return model.Stockpile{}, fmt.Errorf("scan stockpile: %w", err)
	}
	if updatedAt != "" {
		t, err := time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return model.Stockpile{}, fmt.Errorf("parse stockpile updated_at %q: %w", updatedAt, err)
		}
		s.UpdatedAt = t
	}
	return s, nil
}

// formatTime приводит время к UTC и форматирует в RFC3339 для хранения в БД.
func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
