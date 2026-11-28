// Package item_registry описывает модель данных реестра предметов Minecraft.
// Физически записи хранятся в таблице item_localizations (имя предмета на язык).
package item_registry

// ItemLocalization описывает запись таблицы item_localizations.
// Пара (Lang, ItemID) уникальна: у предмета одно название на язык.
// Модель конвертируется в DTO в контроллере, поэтому JSON-тегов нет.
type ItemLocalization struct {
	ID       int64
	Lang     string
	ItemID   string
	ItemName string
}
