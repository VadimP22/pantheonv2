// Package stockpile_resources реализует HTTP-контроллер для ресурсов склада.
// Контроллер — единственное место, где модель преобразуется в DTO (и наоборот);
// параметры запроса разбираются здесь же. Реализация на стандартной net/http.
package stockpile_resources

import (
	"errors"
	"fmt"
	"net/http"

	"backend/dto"
	"backend/lib"
	model "backend/model/stockpile_resources"
	svc "backend/services/stockpile_resources"
)

// Паттерны маршрутов контроллера со служебным префиксом /api.
const (
	resourcesListPattern     = "GET /api/stockpiles/{id}/resources"
	resourcesDepositPattern  = "POST /api/stockpiles/{id}/resources/deposit"
	resourcesWithdrawPattern = "POST /api/stockpiles/{id}/resources/withdraw"
)

// Controller обрабатывает HTTP-запросы к ресурсам склада.
type Controller struct {
	svc *svc.Service
}

// New возвращает новый Controller поверх предоставленного сервиса.
func New(s *svc.Service) *Controller {
	return &Controller{svc: s}
}

// Register регистрирует маршруты контроллера на переданном mux.
func (c *Controller) Register(mux *http.ServeMux) {
	mux.HandleFunc(resourcesListPattern, c.handleList)
	mux.HandleFunc(resourcesDepositPattern, c.handleDeposit)
	mux.HandleFunc(resourcesWithdrawPattern, c.handleWithdraw)
}

// handleList обрабатывает GET /api/stockpiles/{id}/resources — ресурсы склада,
// отсортированные по убыванию количества.
func (c *Controller) handleList(w http.ResponseWriter, r *http.Request) {
	id, err := lib.ParseID(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	items, err := c.svc.ListByStockpile(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	lib.WriteJSON(w, http.StatusOK, toDTOList(items))
}

// handleDeposit обрабатывает POST /api/stockpiles/{id}/resources/deposit —
// пополнение ресурса. Возвращает актуальное состояние записи.
func (c *Controller) handleDeposit(w http.ResponseWriter, r *http.Request) {
	id, err := lib.ParseID(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var in dto.ChangeStockpileResourceInputDTO
	if err := lib.DecodeJSON(r, &in); err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := c.svc.Deposit(id, in.ItemID, in.Count)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	lib.WriteJSON(w, http.StatusOK, toDTO(m))
}

// handleWithdraw обрабатывает POST /api/stockpiles/{id}/resources/withdraw —
// списание ресурса. При нехватке ресурса сервис возвращает ErrInsufficient,
// который маппится в 409 Conflict.
func (c *Controller) handleWithdraw(w http.ResponseWriter, r *http.Request) {
	id, err := lib.ParseID(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var in dto.ChangeStockpileResourceInputDTO
	if err := lib.DecodeJSON(r, &in); err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := c.svc.Withdraw(id, in.ItemID, in.Count); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- конвертация и хелперы ---

// toDTOList конвертирует список моделей в список ответных DTO.
func toDTOList(items []model.StockpileResource) []dto.StockpileResourceDTO {
	out := make([]dto.StockpileResourceDTO, len(items))
	for i, m := range items {
		out[i] = toDTO(m)
	}
	return out
}

// toDTO конвертирует модель в ответный DTO.
func toDTO(m model.StockpileResource) dto.StockpileResourceDTO {
	return dto.StockpileResourceDTO{
		ID:          m.ID,
		StockpileID: m.StockpileID,
		ItemID:      m.ItemID,
		ItemCount:   m.ItemCount,
	}
}

// writeServiceError сопоставляет ошибки сервиса с HTTP-статусами.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, svc.ErrNotFound):
		lib.WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, svc.ErrInsufficient):
		lib.WriteError(w, http.StatusConflict, err)
	case errors.Is(err, svc.ErrInvalidInput):
		lib.WriteError(w, http.StatusBadRequest, err)
	default:
		lib.WriteError(w, http.StatusInternalServerError, fmt.Errorf("stockpile resources: %w", err))
	}
}
