// Package stockpile_resources реализует HTTP-контроллер для ресурсов склада.
// Контроллер — единственное место, где модель преобразуется в DTO (и наоборот);
// параметры запроса разбираются здесь же. Реализация на стандартной net/http.
package stockpile_resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"backend/dto"
	model "backend/model/stockpile_resources"
	svc "backend/services/stockpile_resources"
)

// Паттерны маршрутов контроллера со служебным префиксом /api.
// Именование: controllername + Handlername + Pattern.
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
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	items, err := c.svc.ListByStockpile(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toDTOList(items))
}

// handleDeposit обрабатывает POST /api/stockpiles/{id}/resources/deposit —
// пополнение ресурса. Возвращает актуальное состояние записи.
func (c *Controller) handleDeposit(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var in dto.ChangeStockpileResourceInputDTO
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	m, err := c.svc.Deposit(id, in.ItemID, in.Count)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toDTO(m))
}

// handleWithdraw обрабатывает POST /api/stockpiles/{id}/resources/withdraw —
// списание ресурса. При нехватке ресурса сервис возвращает ErrInsufficient,
// который маппится в 409 Conflict.
func (c *Controller) handleWithdraw(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var in dto.ChangeStockpileResourceInputDTO
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
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

// parseID извлекает и валидирует id из path-переменной {id}.
func parseID(r *http.Request) (int64, error) {
	s := r.PathValue("id")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid id: %q", s)
	}
	return id, nil
}

// decodeJSON декодирует тело запроса; пустое тело считается ошибкой (EOF).
func decodeJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json body: %w", err)
	}
	return nil
}

// writeJSON записывает тело в формате JSON с заданным статусом.
// Если кодирование не удалось (тело уже частично записано), статус не меняем.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError записывает сообщение об ошибке в виде {"error":"..."}.
func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// writeServiceError сопоставляет ошибки сервиса с HTTP-статусами.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, svc.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, svc.ErrInsufficient):
		writeError(w, http.StatusConflict, err)
	case errors.Is(err, svc.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}
