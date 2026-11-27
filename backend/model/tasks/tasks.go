// Package tasks описывает модель данных для задач.
package tasks

// UpdateTaskInput описывает данные для обновления задачи.
// nil-поле = не изменять; не nil = присвоить новое значение.
// Идентификатор не входит в input — он передаётся отдельно в метод обновления.
type UpdateTaskInput struct {
	Name            *string
	Description     *string
	IconItemName    *string
	IsCheckbox      *bool
	Done            *bool
	ItemID          *string
	ItemCount       *int64
	TargetItemCount *int64
	StockpileID     *int64
}

// Task описывает запись таблицы tasks.
// Модель конвертируется в DTO в сервисном слое, поэтому JSON-тегов нет.
type Task struct {
	ID           int64
	Name         string
	Description  string
	IconItemName string
	// IsCheckbox — признак «отметки» (0/1 в БД, bool в модели).
	IsCheckbox bool
	// Done — признак выполнения задачи (0/1 в БД, bool в модели).
	Done            bool
	ItemID          string
	ItemCount       int64
	TargetItemCount int64
	// StockpileID — идентификатор целевого запаса (склада);
	// nil означает, что целевой склад не задан.
	StockpileID *int64
}
