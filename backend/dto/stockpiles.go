// Package dto содержит DTO (Data Transfer Objects) для HTTP API.
// Пакет глобальный (плоский): один файл на ресурс, что удобно для генерации
// типов (tygo). Каждый DTO назван по правилу <Имя>DTO.
package dto

// StockpileDTO описывает складскую позицию в ответах API.
type StockpileDTO struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Dimension     string  `json:"dimension"`
	IconName      string  `json:"icon_name"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Z             float64 `json:"z"`
	UpdatedAt     string  `json:"updated_at"`
	ItemCount     int64   `json:"item_count"`
	ItemTypeCount int64   `json:"item_type_count"`
}

// CreateStockpileInputDTO — тело запроса на создание складской позиции.
type CreateStockpileInputDTO struct {
	Name      string  `json:"name"`
	Dimension string  `json:"dimension"`
	IconName  string  `json:"icon_name"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Z         float64 `json:"z"`
}

// UpdateStockpileInputDTO — тело запроса на обновление складской позиции.
// nil-поле = не изменять; не nil = присвоить новое значение.
type UpdateStockpileInputDTO struct {
	Name      *string  `json:"name"`
	Dimension *string  `json:"dimension"`
	IconName  *string  `json:"icon_name"`
	X         *float64 `json:"x"`
	Y         *float64 `json:"y"`
	Z         *float64 `json:"z"`
}

// ListQueryDTO — разобранные параметры пагинации из query-строки (без фильтров).
type ListQueryDTO struct {
	Page     int
	PageSize int
}

// StockpilePageDTO — страница результатов со счётчиками пагинации.
type StockpilePageDTO struct {
	Items      []StockpileDTO `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}
