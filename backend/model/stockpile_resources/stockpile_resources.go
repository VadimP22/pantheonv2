// Package stockpile_resources описывает модель данных для ресурсов склада.
package stockpile_resources

// StockpileResource описывает запись таблицы stockpile_resources.
// Модель конвертируется в DTO в сервисном слое, поэтому JSON-тегов нет.
type StockpileResource struct {
	ID          int64
	StockpileID int64
	ItemID      string
	ItemCount   int64
}

// UpdateStockpileResourceInput описывает данные для обновления ресурса склада.
// nil-поле = не изменять; не nil = присвоить новое значение.
// Идентификатор не входит в input — он передаётся отдельно в метод обновления.
type UpdateStockpileResourceInput struct {
	StockpileID *int64
	ItemID      *string
	ItemCount   *int64
}
