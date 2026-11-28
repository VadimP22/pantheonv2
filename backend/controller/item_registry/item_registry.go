// Package item_registry реализует HTTP-контроллер реестра предметов Minecraft.
// Контроллер — единственное место, где модель преобразуется в DTO;
// параметры пути разбираются здесь же.
package item_registry

import (
	"errors"
	"fmt"
	"net/http"

	"backend/dto"
	"backend/lib"
	model "backend/model/item_registry"
	svc "backend/services/item_registry"
)

// Паттерн маршрута контроллера со служебным префиксом /api.
// Маршрут /api/loc/... сохранён для обратной совместимости API.
const itemRegistryGetPattern = "GET /api/loc/{lang}/{item}"

// ItemRegistryController обрабатывает HTTP-запросы к реестру предметов.
type ItemRegistryController struct {
	svc *svc.ItemRegistryService
}

// New возвращает новый ItemRegistryController поверх предоставленного сервиса.
func New(s *svc.ItemRegistryService) *ItemRegistryController {
	return &ItemRegistryController{svc: s}
}

// Register регистрирует маршруты контроллера на переданном mux.
func (c *ItemRegistryController) Register(mux *http.ServeMux) {
	mux.HandleFunc(itemRegistryGetPattern, c.handleGetByItemID)
}

// handleGetByItemID обрабатывает GET /api/loc/{lang}/{item} — локализация предмета.
func (c *ItemRegistryController) handleGetByItemID(w http.ResponseWriter, r *http.Request) {
	lang := r.PathValue("lang")
	item := r.PathValue("item")

	m, err := c.svc.GetByItemID(lang, item)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	lib.WriteJSON(w, http.StatusOK, toDTO(m))
}

// toDTO конвертирует модель в ответный DTO.
func toDTO(m model.ItemLocalization) dto.ItemLocalizationDTO {
	return dto.ItemLocalizationDTO{
		ID:       m.ID,
		Lang:     m.Lang,
		ItemID:   m.ItemID,
		ItemName: m.ItemName,
	}
}

// writeServiceError сопоставляет ошибки сервиса с HTTP-статусами.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, svc.ErrNotFound):
		lib.WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, svc.ErrInvalidInput):
		lib.WriteError(w, http.StatusBadRequest, err)
	default:
		lib.WriteError(w, http.StatusInternalServerError, fmt.Errorf("item registry: %w", err))
	}
}
