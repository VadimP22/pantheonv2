// Package tasks реализует сервисный слой для задач.
// Сервис инкапсулирует бизнес-логику и валидацию поверх репозитория, работая
// напрямую с типами модели. DTO/JSON-представление вводятся на уровне API —
// отдельной задачей.
package tasks

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"backend/dto"
	model "backend/model/tasks"
	repo "backend/repository/tasks"
	svcresources "backend/services/stockpile_resources"
)

// Ошибки уровня сервиса: они оборачивают или дополняют ошибки репозитория,
// предоставляя вызывающему коду определённую семантику.
var (
	// ErrNotFound возвращается, когда запись с запрошенным id не найдена.
	ErrNotFound = errors.New("tasks: not found")
	// ErrInvalidInput возвращается при невалидных данных.
	ErrInvalidInput = errors.New("tasks: invalid input")
)

// Service обеспечивает операции над задачами с валидацией.
type Service struct {
	repo      *repo.Repository
	resources *svcresources.Service
}

// New возвращает новый Service, работающий поверх предоставленного репозитория
// и сервиса stockpile_resources.
func New(r *repo.Repository, resources *svcresources.Service) *Service {
	return &Service{repo: r, resources: resources}
}

// Create создаёт новую задачу и возвращает её идентификатор.
// Поле name обязательно (непустая строка).
//
// Правила по типу задачи:
//   - is_checkbox = false (ресурс): обязателен item_id; icon_item_name
//     заполняется сервисом из item_id (фронтенд его не собирает и не шлёт).
//   - is_checkbox = true (чекбокс): обязателен icon_item_name;
//     item_id и количества не требуются.
//
// stockpile_id опционален в обоих случаях (nil — целевой склад не задан).
func (s *Service) Create(in dto.CreateTaskInputDTO) (int64, error) {
	if strings.TrimSpace(in.Name) == "" {
		return 0, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}

	iconItemName := in.IconItemName
	if !in.IsCheckbox {
		if strings.TrimSpace(in.ItemID) == "" {
			return 0, fmt.Errorf("%w: item_id is required for resource tasks", ErrInvalidInput)
		}
		iconItemName = in.ItemID
	} else {
		if strings.TrimSpace(in.IconItemName) == "" {
			return 0, fmt.Errorf("%w: icon_item_name is required for checkbox tasks", ErrInvalidInput)
		}
	}

	id, err := s.repo.Create(in.Name, in.Description, iconItemName, in.IsCheckbox, false, in.ItemID, in.ItemCount, in.TargetItemCount, in.StockpileID)
	if err != nil {
		return 0, fmt.Errorf("create task: %w", err)
	}

	return id, nil
}

// GetByID возвращает задачу по идентификатору.
// Если запись не найдена, возвращается ErrNotFound.
func (s *Service) GetByID(id int64) (model.Task, error) {
	res, err := s.repo.GetByID(id)
	switch {
	case err == nil:
		return *res, nil
	case errors.Is(err, sql.ErrNoRows):
		return model.Task{}, ErrNotFound
	default:
		return model.Task{}, err
	}
}

// ListPaged возвращает страницу задач: не более limit записей,
// начиная со смещения offset. При пустой странице — пустой не-nil слайс.
func (s *Service) ListPaged(limit, offset int) ([]model.Task, error) {
	return s.repo.ListPaged(limit, offset)
}

// Count возвращает общее количество задач.
func (s *Service) Count() (int64, error) {
	return s.repo.Count()
}

// Deposit пополняет запас задачи: upsert в stockpile_resources и рост item_count.
func (s *Service) Deposit(id int64, count int64) (model.Task, error) {
	if count < 1 {
		return model.Task{}, fmt.Errorf("tasks: deposit count %d: %w", count, ErrInvalidInput)
	}
	t, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, fmt.Errorf("tasks: deposit %d: %w", id, ErrNotFound)
		}
		return model.Task{}, fmt.Errorf("tasks: deposit get %d: %w", id, err)
	}
	if t.IsCheckbox {
		return model.Task{}, fmt.Errorf("tasks: deposit %d is checkbox: %w", id, ErrInvalidInput)
	}
	if t.StockpileID == nil {
		return model.Task{}, fmt.Errorf("tasks: deposit %d has no stockpile: %w", id, ErrInvalidInput)
	}
	trimmed := strings.TrimSpace(t.ItemID)
	if trimmed == "" {
		return model.Task{}, fmt.Errorf("tasks: deposit %d has empty item_id: %w", id, ErrInvalidInput)
	}
	if _, err := s.resources.Deposit(*t.StockpileID, trimmed, count); err != nil {
		return model.Task{}, fmt.Errorf("tasks: deposit stockpile write %d: %w", id, err)
	}
	if ok, err := s.repo.AddItemCount(id, count); err != nil {
		return model.Task{}, fmt.Errorf("tasks: deposit add count %d: %w", id, err)
	} else if !ok {
		return model.Task{}, fmt.Errorf("tasks: deposit %d not found after write: %w", id, ErrNotFound)
	}
	res, err := s.repo.GetByID(id)
	if err != nil {
		return model.Task{}, fmt.Errorf("tasks: deposit get after write %d: %w", id, err)
	}
	return *res, nil
}

// Update обновляет у задачи id только те поля input, которые не nil,
// и возвращает обновлённую запись. Если в input нет ни одного не-nil поля,
// возвращается ErrInvalidInput; если запись не найдена — ErrNotFound.
func (s *Service) Update(id int64, in model.UpdateTaskInput) (model.Task, error) {
	if !hasUpdateFields(in) {
		return model.Task{}, ErrInvalidInput
	}

	// Валидация по типу задачи, если соответствующие поля переданы.
	if in.IsCheckbox != nil && !*in.IsCheckbox {
		if in.ItemID == nil || strings.TrimSpace(*in.ItemID) == "" {
			return model.Task{}, fmt.Errorf("%w: item_id is required for resource tasks", ErrInvalidInput)
		}
	}
	if in.IsCheckbox != nil && *in.IsCheckbox {
		if in.IconItemName == nil || strings.TrimSpace(*in.IconItemName) == "" {
			return model.Task{}, fmt.Errorf("%w: icon_item_name is required for checkbox tasks", ErrInvalidInput)
		}
	}
	if in.Name != nil && strings.TrimSpace(*in.Name) == "" {
		return model.Task{}, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}

	updated, err := s.repo.Update(id, in)
	switch {
	case errors.Is(err, repo.ErrNoUpdateFields):
		return model.Task{}, ErrInvalidInput
	case err != nil:
		return model.Task{}, err
	case !updated:
		return model.Task{}, ErrNotFound
	}

	return s.GetByID(id)
}

// hasUpdateFields определяет, есть ли в input хотя бы одно не-nil поле.
func hasUpdateFields(in model.UpdateTaskInput) bool {
	return in.Name != nil ||
		in.Description != nil ||
		in.IconItemName != nil ||
		in.IsCheckbox != nil ||
		in.Done != nil ||
		in.ItemID != nil ||
		in.ItemCount != nil ||
		in.TargetItemCount != nil ||
		in.StockpileID != nil
}
