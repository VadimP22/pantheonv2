// Package item_localizations описывает модель данных для локализованных
// названий предметов Minecraft.
package item_localizations

// ItemLocalization описывает запись таблицы item_localizations.
// Пара (Lang, ItemID) уникальна: у предмета одно название на язык.
// Модель конвертируется в DTO в контроллере, поэтому JSON-тегов нет.
type ItemLocalization struct {
	ID       int64
	Lang     string
	ItemID   string
	ItemName string
}
