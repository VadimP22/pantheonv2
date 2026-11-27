// Package item_localizations реализует сервисный слой для локализованных
// названий предметов Minecraft. Сервис валидирует вход и транслирует ошибки
// репозитория в семантические ошибки уровня сервиса.
package item_localizations

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	model "backend/model/item_localizations"
	repo "backend/repository/item_localizations"
)

// Ошибки уровня сервиса.
var (
	// ErrNotFound возвращается, когда локализация для пары «язык + предмет» не найдена.
	ErrNotFound = errors.New("item_localizations: not found")
	// ErrInvalidInput возвращается при пустом lang или item.
	ErrInvalidInput = errors.New("item_localizations: invalid input")
)

// ItemLocalizationsService обеспечивает чтение локализаций предметов.
type ItemLocalizationsService struct {
	repo *repo.Repository
}

// New возвращает новый ItemLocalizationsService поверх предоставленного репозитория.
func New(r *repo.Repository) *ItemLocalizationsService {
	return &ItemLocalizationsService{repo: r}
}

// GetByItemID возвращает локализацию по паре «язык + предмет».
// Если запись не найдена, возвращается ErrNotFound.
func (s *ItemLocalizationsService) GetByItemID(lang, itemID string) (model.ItemLocalization, error) {
	lang = strings.TrimSpace(lang)
	itemID = strings.TrimSpace(itemID)
	if lang == "" || itemID == "" {
		return model.ItemLocalization{}, fmt.Errorf("%w: lang and item are required", ErrInvalidInput)
	}

	res, err := s.repo.GetByItemID(lang, itemID)
	switch {
	case err == nil:
		return *res, nil
	case errors.Is(err, sql.ErrNoRows):
		return model.ItemLocalization{}, ErrNotFound
	default:
		return model.ItemLocalization{}, err
	}
}
