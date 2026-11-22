// Package dto содержит DTO (Data Transfer Objects) для HTTP API.
// Пакет глобальный (плоский): один файл на ресурс, что удобно для генерации
// типов (tygo). Каждый DTO назван по правилу <Имя>DTO.
package dto

// StockpileResourceDTO — ресурс склада в ответах API.
type StockpileResourceDTO struct {
	ID          int64  `json:"id"`
	StockpileID int64  `json:"stockpile_id"`
	ItemID      string `json:"item_id"`
	ItemCount   int64  `json:"item_count"`
}

// ChangeStockpileResourceInputDTO — тело deposit/withdraw. stockpile_id приходит
// из path-переменной, поэтому в теле его нет.
type ChangeStockpileResourceInputDTO struct {
	ItemID string `json:"item_id"`
	Count  int64  `json:"count"`
}
