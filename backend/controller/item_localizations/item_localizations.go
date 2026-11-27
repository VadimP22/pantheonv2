// Package item_localizations реализует HTTP-контроллер для локализованных
// названий предметов Minecraft. Контроллер — единственное место, где модель
// преобразуется в DTO; параметры пути разбираются здесь же.
package item_localizations

import (
	"errors"
	"fmt"
	"net/http"

	"backend/dto"
	"backend/lib"
	model "backend/model/item_localizations"
	svc "backend/services/item_localizations"
)

// Паттерн маршрута контроллера со служебным префиксом /api.
const itemLocalizationsGetPattern = "GET /api/loc/{lang}/{item}"

// ItemLocalizationsController обрабатывает HTTP-запросы к ресурсу item_localizations.
type ItemLocalizationsController struct {
	svc *svc.ItemLocalizationsService
}

// New возвращает новый ItemLocalizationsController поверх предоставленного сервиса.
func New(s *svc.ItemLocalizationsService) *ItemLocalizationsController {
	return &ItemLocalizationsController{svc: s}
}

// Register регистрирует маршруты контроллера на переданном mux.
func (c *ItemLocalizationsController) Register(mux *http.ServeMux) {
	mux.HandleFunc(itemLocalizationsGetPattern, c.handleGetByItemID)
}

// handleGetByItemID обрабатывает GET /api/loc/{lang}/{item} — локализация предмета.
func (c *ItemLocalizationsController) handleGetByItemID(w http.ResponseWriter, r *http.Request) {
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
		lib.WriteError(w, http.StatusInternalServerError, fmt.Errorf("item localizations: %w", err))
	}
}
