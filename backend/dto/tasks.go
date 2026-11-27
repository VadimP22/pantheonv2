// Package dto содержит DTO (Data Transfer Objects) для HTTP API.
// Пакет глобальный (плоский): один файл на ресурс, что удобно для генерации
// типов (tygo). Каждый DTO назван по правилу <Имя>DTO.
package dto

// TaskDTO описывает задачу в ответах API.
type TaskDTO struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	IconItemName    string `json:"icon_item_name"`
	IsCheckbox      bool   `json:"is_checkbox"`
	Done            bool   `json:"done"`
	Deleted         bool   `json:"deleted"`
	ItemID          string `json:"item_id"`
	ItemCount       int64  `json:"item_count"`
	TargetItemCount int64  `json:"target_item_count"`
	// StockpileID — идентификатор целевого запаса; null, если не задан.
	StockpileID *int64 `json:"stockpile_id"`
}

// TaskPageDTO — страница результатов со счётчиками пагинации.
type TaskPageDTO struct {
	Items      []TaskDTO `json:"items"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}

// UpdateTaskInputDTO — тело запроса на обновление задачи.
// nil-поле = не изменять; не nil = присвоить новое значение.
type UpdateTaskInputDTO struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	IconItemName    *string `json:"icon_item_name"`
	IsCheckbox      *bool   `json:"is_checkbox"`
	Done            *bool   `json:"done"`
	ItemID          *string `json:"item_id"`
	ItemCount       *int64  `json:"item_count"`
	TargetItemCount *int64  `json:"target_item_count"`
	StockpileID     *int64  `json:"stockpile_id"`
}

// CreateTaskInputDTO — тело запроса на создание задачи.
//
// Правила по типу задачи (is_checkbox):
//   - false (ресурс): обязателен item_id; icon_item_name игнорируется и
//     заполняется сервисом из item_id (клиент его не собирает).
//   - true (чекбокс): обязателен icon_item_name; item_id и количества не требуются.
//
// stockpile_id опционален в обоих случаях: nil означает, что целевой склад не задан.
type CreateTaskInputDTO struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	IconItemName    string `json:"icon_item_name"`
	IsCheckbox      bool   `json:"is_checkbox"`
	ItemID          string `json:"item_id"`
	ItemCount       int64  `json:"item_count"`
	TargetItemCount int64  `json:"target_item_count"`
	StockpileID     *int64 `json:"stockpile_id"`
}

// DepositTaskInputDTO — запрос пополнения ресурсной задачи.
type DepositTaskInputDTO struct {
	Count int64 `json:"count"`
}

// DeleteTaskInputDTO — входные данные для удаления задачи (пустой).
type DeleteTaskInputDTO struct{}
