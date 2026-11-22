// Package stockpiles реализует HTTP-контроллер для складских позиций (запасов).
// Контроллер — единственное место, где модель преобразуется в DTO (и наоборот);
// параметры запроса разбираются здесь же. Реализация на стандартной net/http.
package stockpiles

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"backend/dto"
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

// Параметры пагинации списка (без фильтров).
const (
	defaultPage     = 1
	defaultPageSize = 6
	maxPageSize     = 100
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
	query, err := parseListQuery(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	offset := (query.Page - 1) * query.PageSize
	items, err := c.svc.ListPaged(query.PageSize, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("list stockpiles: %w", err))
		return
	}
	total, err := c.svc.Count()
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("count stockpiles: %w", err))
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
	writeJSON(w, http.StatusOK, page)
}

// handleCreate обрабатывает POST /api/stockpiles — создание записи.
func (c *Controller) handleCreate(w http.ResponseWriter, r *http.Request) {
	var in dto.CreateStockpileInputDTO
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
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
	writeJSON(w, http.StatusCreated, toDTO(m))
}

// handleGetByID обрабатывает GET /api/stockpiles/{id}.
func (c *Controller) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	m, err := c.svc.GetByID(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toDTO(m))
}

// handleUpdate обрабатывает PUT /api/stockpiles/{id} — частичное обновление.
func (c *Controller) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var in dto.UpdateStockpileInputDTO
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	m, err := c.svc.Update(id, toUpdateModel(in))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toDTO(m))
}

// handleDelete обрабатывает DELETE /api/stockpiles/{id}. Удаление идемпотентно:
// отсутствие записи не является ошибкой.
func (c *Controller) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := c.svc.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("delete stockpile %d: %w", id, err))
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
	return dto.StockpileDTO{
		ID:        m.ID,
		Name:      m.Name,
		Dimension: m.Dimension,
		IconName:  m.IconName,
		X:         m.X,
		Y:         m.Y,
		Z:         m.Z,
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

// parseListQuery разбирает параметры пагинации page и page_size из query-строки.
// Пустое значение = дефолт; некорректное значение — ошибка.
func parseListQuery(r *http.Request) (dto.ListQueryDTO, error) {
	q := dto.ListQueryDTO{Page: defaultPage, PageSize: defaultPageSize}

	if v := r.URL.Query().Get("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return q, fmt.Errorf("invalid page: %q", v)
		}
		q.Page = n
	}
	if v := r.URL.Query().Get("page_size"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > maxPageSize {
			return q, fmt.Errorf("invalid page_size: %q", v)
		}
		q.PageSize = n
	}
	return q, nil
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
	case errors.Is(err, svc.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}
