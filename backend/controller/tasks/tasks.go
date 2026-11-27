// Package tasks реализует HTTP-контроллер для задач.
// Контроллер — единственное место, где модель преобразуется в DTO (и наоборот);
// параметры запроса разбираются здесь же. Реализация на стандартной net/http.
package tasks

import (
	"errors"
	"fmt"
	"net/http"

	"backend/dto"
	"backend/lib"
	model "backend/model/tasks"
	svc "backend/services/tasks"
)

// Паттерны маршрутов контроллера со служебным префиксом /api.
const (
	tasksListPattern    = "GET /api/tasks"
	tasksCreatePattern  = "POST /api/tasks"
	tasksGetPattern     = "GET /api/tasks/{id}"
	tasksUpdatePattern  = "PUT /api/tasks/{id}"
	tasksDepositPattern = "POST /api/tasks/{id}/deposit"
)

// Controller обрабатывает HTTP-запросы к ресурсу tasks.
type Controller struct {
	svc *svc.Service
}

// New возвращает новый Controller поверх предоставленного сервиса.
func New(s *svc.Service) *Controller {
	return &Controller{svc: s}
}

// Register регистрирует маршруты контроллера на переданном mux.
func (c *Controller) Register(mux *http.ServeMux) {
	mux.HandleFunc(tasksListPattern, c.handleList)
	mux.HandleFunc(tasksCreatePattern, c.handleCreate)
	mux.HandleFunc(tasksGetPattern, c.handleGetByID)
	mux.HandleFunc(tasksUpdatePattern, c.handleUpdate)
	mux.HandleFunc(tasksDepositPattern, c.handleDeposit)
}

// handleCreate обрабатывает POST /api/tasks — создание записи.
func (c *Controller) handleCreate(w http.ResponseWriter, r *http.Request) {
	var in dto.CreateTaskInputDTO
	if err := lib.DecodeJSON(r, &in); err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := c.svc.Create(in)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	m, err := c.svc.GetByID(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	lib.WriteJSON(w, http.StatusCreated, toDTO(m))
}

// handleList обрабатывает GET /api/tasks — страница задач (без фильтров).
func (c *Controller) handleList(w http.ResponseWriter, r *http.Request) {
	query, err := lib.ParseListQuery(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	offset := (query.Page - 1) * query.PageSize
	items, err := c.svc.ListPaged(query.PageSize, offset)
	if err != nil {
		lib.WriteError(w, http.StatusInternalServerError, fmt.Errorf("list tasks: %w", err))
		return
	}
	total, err := c.svc.Count()
	if err != nil {
		lib.WriteError(w, http.StatusInternalServerError, fmt.Errorf("count tasks: %w", err))
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(query.PageSize) - 1) / int64(query.PageSize))
	}

	page := dto.TaskPageDTO{
		Items:      toDTOList(items),
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}
	lib.WriteJSON(w, http.StatusOK, page)
}

// handleGetByID обрабатывает GET /api/tasks/{id}.
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

// toDTO конвертирует модель в ответный DTO.
func toDTO(m model.Task) dto.TaskDTO {
	return dto.TaskDTO{
		ID:              m.ID,
		Name:            m.Name,
		Description:     m.Description,
		IconItemName:    m.IconItemName,
		IsCheckbox:      m.IsCheckbox,
		Done:            m.Done,
		ItemID:          m.ItemID,
		ItemCount:       m.ItemCount,
		TargetItemCount: m.TargetItemCount,
		StockpileID:     m.StockpileID,
	}
}

// toDTOList конвертирует список моделей в список ответных DTO.
func toDTOList(items []model.Task) []dto.TaskDTO {
	out := make([]dto.TaskDTO, len(items))
	for i, m := range items {
		out[i] = toDTO(m)
	}
	return out
}

// toUpdateModel конвертирует DTO обновления в модель обновления.
func toUpdateModel(in dto.UpdateTaskInputDTO) model.UpdateTaskInput {
	return model.UpdateTaskInput{
		Name:            in.Name,
		Description:     in.Description,
		IconItemName:    in.IconItemName,
		IsCheckbox:      in.IsCheckbox,
		Done:            in.Done,
		ItemID:          in.ItemID,
		ItemCount:       in.ItemCount,
		TargetItemCount: in.TargetItemCount,
		StockpileID:     in.StockpileID,
	}
}

// handleUpdate обрабатывает PUT /api/tasks/{id} — частичное обновление.
func (c *Controller) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := lib.ParseID(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var in dto.UpdateTaskInputDTO
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

// handleDeposit обрабатывает POST /api/tasks/{id}/deposit — пополнение ресурсной задачи.
func (c *Controller) handleDeposit(w http.ResponseWriter, r *http.Request) {
	id, err := lib.ParseID(r)
	if err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var in dto.DepositTaskInputDTO
	if err := lib.DecodeJSON(r, &in); err != nil {
		lib.WriteError(w, http.StatusBadRequest, err)
		return
	}

	m, err := c.svc.Deposit(id, in.Count)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	lib.WriteJSON(w, http.StatusOK, toDTO(m))
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
