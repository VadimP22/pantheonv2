// Package item_registry реализует сервисный слой реестра предметов Minecraft:
// чтение локализованных названий и проверка существования предмета.
// Сервис валидирует вход и транслирует ошибки репозитория в семантические
// ошибки уровня сервиса.
package item_registry

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	model "backend/model/item_registry"
	repo "backend/repository/item_registry"
)

// Ошибки уровня сервиса.
var (
	// ErrNotFound возвращается, когда запись для пары «язык + предмет» не найдена.
	ErrNotFound = errors.New("item_registry: not found")
	// ErrInvalidInput возвращается при пустом lang или item.
	ErrInvalidInput = errors.New("item_registry: invalid input")
)

// ItemRegistryService обеспечивает чтение реестра предметов.
type ItemRegistryService struct {
	repo *repo.Repository
}

// New возвращает новый ItemRegistryService поверх предоставленного репозитория.
func New(r *repo.Repository) *ItemRegistryService {
	return &ItemRegistryService{repo: r}
}

// GetByItemID возвращает локализацию по паре «язык + предмет».
// Если запись не найдена, возвращается ErrNotFound.
func (s *ItemRegistryService) GetByItemID(lang, itemID string) (model.ItemLocalization, error) {
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
