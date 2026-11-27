// Package dto содержит DTO (Data Transfer Objects) для HTTP API.
// Пакет глобальный (плоский): один файл на ресурс, что удобно для генерации
// типов (tygo). Каждый DTO назван по правилу <Имя>DTO.
package dto

// ItemLocalizationDTO — локализованное название предмета в ответах API.
type ItemLocalizationDTO struct {
	ID       int64  `json:"id"`
	Lang     string `json:"lang"`
	ItemID   string `json:"item_id"`
	ItemName string `json:"item_name"`
}
