// Package stockpiles описывает модель данных для запасов (складских позиций).
package stockpiles

import "time"

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
	// UpdatedAt — время последнего изменения записи (UTC).
	UpdatedAt time.Time
	// ItemCount — суммарное количество единиц предметов, находящихся на запасе
	// (вычисляется на лету из таблицы stockpile_resources).
	ItemCount int64
	// ItemTypeCount — количество различных типов предметов на запасе
	// (вычисляется на лету из таблицы stockpile_resources).
	ItemTypeCount int64
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
