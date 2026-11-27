// Package stockpiles реализует HTTP-контроллер для складских позиций (запасов).
// Контроллер — единственное место, где модель преобразуется в DTO (и наоборот);
// параметры запроса разбираются здесь же. Реализация на стандартной net/http.
package stockpiles

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"backend/dto"
	"backend/lib"
	model "backend/model/stockpiles"
	svc "backend/services/stockpiles"
)

// Паттерны маршрутов контроллера со служебным префиксом /api.
// Именование: controllername + Handlername + Pattern.
const (
	stockpilesListPattern   = "GET /api/stockpiles"
	stockpilesCreatePattern = "POST /api/stockpiles"
	stockpilesGetPattern    = "GET /api/stockpiles/{id}"
	stockpilesUpdatePattern = "PUT /api/stockpiles/{id}"
	stockpilesDeletePattern = "DELETE /api/stockpiles/{id}"
)

// Controller обрабатывает HTTP-запросы к ресурсу stockpiles.
type Controller struct {
	svc *svc.Service
}

// New возвращает новый Controller поверх предоставленного сервиса.
func New(s *svc.Service) *Controller {
	return &Controller{svc: s}
}

// Register регистрирует маршруты контроллера на переданном mux.
func (c *Controller) Register(mux *http.ServeMux) {
	mux.HandleFunc(stockpilesListPattern, c.handleList)
	mux.HandleFunc(stockpilesCreatePattern, c.handleCreate)
	mux.HandleFunc(stockpilesGetPattern, c.handleGetByID)
	mux.HandleFunc(stockpilesUpdatePattern, c.handleUpdate)
	mux.HandleFunc(stockpilesDeletePattern, c.handleDelete)
}

// handleList обрабатывает GET /api/stockpiles — страница записей (без фильтров).
func (c *Controller) handleList(w http.ResponseWriter, r *http.Request) {
	query, err := lib.ParseListQuery(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	offset := (query.Page - 1) * query.PageSize
	items, err := c.svc.ListPaged(query.PageSize, offset)
	if err != nil {
		lib.WriteError(w, http.StatusInternalServerError, fmt.Errorf("list stockpiles: %w", err))
		return
	}
	total, err := c.svc.Count()
	if err != nil {
		lib.WriteError(w, http.StatusInternalServerError, fmt.Errorf("count stockpiles: %w", err))
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(query.PageSize) - 1) / int64(query.PageSize))
	}

	page := dto.StockpilePageDTO{
		Items:      toDTOList(items),
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}
	lib.WriteJSON(w, http.StatusOK, page)
}

// handleCreate обрабатывает POST /api/stockpiles — создание записи.
func (c *Controller) handleCreate(w http.ResponseWriter, r *http.Request) {
	var in dto.CreateStockpileInputDTO
	if err := lib.DecodeJSON(r, &in); err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := c.svc.Create(model.Stockpile{
		Name:      in.Name,
		Dimension: in.Dimension,
		IconName:  in.IconName,
		X:         in.X,
		Y:         in.Y,
		Z:         in.Z,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	lib.WriteJSON(w, http.StatusCreated, toDTO(m))
}

// handleGetByID обрабатывает GET /api/stockpiles/{id}.
func (c *Controller) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := lib.ParseID(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := c.svc.GetByID(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	lib.WriteJSON(w, http.StatusOK, toDTO(m))
}

// handleUpdate обрабатывает PUT /api/stockpiles/{id} — частичное обновление.
func (c *Controller) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := lib.ParseID(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var in dto.UpdateStockpileInputDTO
	if err := lib.DecodeJSON(r, &in); err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := c.svc.Update(id, toUpdateModel(in))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	lib.WriteJSON(w, http.StatusOK, toDTO(m))
}

// handleDelete обрабатывает DELETE /api/stockpiles/{id}. Удаление идемпотентно:
// отсутствие записи не является ошибкой.
func (c *Controller) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := lib.ParseID(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := c.svc.Delete(id); err != nil {
		lib.WriteError(w, http.StatusInternalServerError, fmt.Errorf("delete stockpile %d: %w", id, err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- конвертация и хелперы ---

// toDTOList конвертирует список моделей в список ответных DTO.
func toDTOList(items []model.Stockpile) []dto.StockpileDTO {
	out := make([]dto.StockpileDTO, len(items))
	for i, m := range items {
		out[i] = toDTO(m)
	}
	return out
}

// toDTO конвертирует модель в ответный DTO.
func toDTO(m model.Stockpile) dto.StockpileDTO {
	updatedAt := ""
	if !m.UpdatedAt.IsZero() {
		updatedAt = m.UpdatedAt.UTC().Format(time.RFC3339)
	}

	return dto.StockpileDTO{
		ID:            m.ID,
		Name:          m.Name,
		Dimension:     m.Dimension,
		IconName:      m.IconName,
		X:             m.X,
		Y:             m.Y,
		Z:             m.Z,
		UpdatedAt:     updatedAt,
		ItemCount:     m.ItemCount,
		ItemTypeCount: m.ItemTypeCount,
	}
}

// toUpdateModel конвертирует DTO обновления в модель обновления.
func toUpdateModel(in dto.UpdateStockpileInputDTO) model.UpdateStockpileInput {
	return model.UpdateStockpileInput{
		Name:      in.Name,
		Dimension: in.Dimension,
		IconName:  in.IconName,
		X:         in.X,
		Y:         in.Y,
		Z:         in.Z,
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
		lib.WriteError(w, http.StatusInternalServerError, err)
	}
}
