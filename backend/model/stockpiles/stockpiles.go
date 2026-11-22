// Package stockpiles описывает модель данных для запасов (складских позиций).
package stockpiles

// Stockpile описывает запись таблицы stockpiles.
// Модель конвертируется в DTO в сервисном слое, поэтому JSON-тегов нет.
type Stockpile struct {
	ID        int64
	Name      string
	Dimension string
	IconName  string
	X         float64
	Y         float64
	Z         float64
}

// UpdateStockpileInput описывает данные для обновления складской позиции.
// nil-поле = не изменять; не nil = присвоить новое значение.
// Идентификатор не входит в input — он передаётся отдельно в метод обновления.
type UpdateStockpileInput struct {
	Name      *string
	Dimension *string
	IconName  *string
	X         *float64
	Y         *float64
	Z         *float64
}
