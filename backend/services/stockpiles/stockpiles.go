// Package stockpiles реализует сервисный слой для складских позиций (запасов).
// Сервис инкапсулирует бизнес-логику и валидацию поверх репозитория, работая
// напрямую с типами модели. DTO/JSON-представление вводятся на уровне API —
// отдельной задачей.
package stockpiles

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	model "backend/model/stockpiles"
	repo "backend/repository/stockpiles"
)

// Ошибки уровня сервиса: они оборачивают или дополняют ошибки репозитория,
// предоставляя вызывающему коду определённую семантику.
var (
	// ErrNotFound возвращается, когда запись с запрошенным id не найдена.
	ErrNotFound = errors.New("stockpiles: not found")
	// ErrInvalidInput возвращается при невалидных данных либо когда обновлять нечего.
	ErrInvalidInput = errors.New("stockpiles: invalid input")
)

// Service обеспечивает CRUD-операции над складскими позициями с валидацией.
// Методы не принимают context.Context для согласованности с текущим репозиторием.
type Service struct {
	repo *repo.Repository
}

// New возвращает новый Service, работающий поверх предоставленного репозитория.
func New(r *repo.Repository) *Service {
	return &Service{repo: r}
}

// Create создаёт новую складскую позицию и возвращает её с заполненным ID.
// Поле ID во входной структуре игнорируется (присваивается БД).
func (s *Service) Create(in model.Stockpile) (model.Stockpile, error) {
	if err := validateRequired(in); err != nil {
		return model.Stockpile{}, err
	}

	m := model.Stockpile{
		Name:      in.Name,
		Dimension: in.Dimension,
		IconName:  in.IconName,
		X:         in.X,
		Y:         in.Y,
		Z:         in.Z,
	}

	if err := s.repo.Create(&m); err != nil {
		return model.Stockpile{}, fmt.Errorf("create stockpile: %w", err)
	}

	return m, nil
}

// GetByID возвращает складскую позицию по идентификатору.
// Если запись не найдена, возвращается ErrNotFound.
func (s *Service) GetByID(id int64) (model.Stockpile, error) {
	res, err := s.repo.GetByID(id)
	switch {
	case err == nil:
		return *res, nil
	case errors.Is(err, sql.ErrNoRows):
		return model.Stockpile{}, ErrNotFound
	default:
		return model.Stockpile{}, err
	}
}

// List возвращает все складские позиции, упорядоченные по id.
// При пустой таблице возвращается пустой не-nil слайс.
func (s *Service) List() ([]model.Stockpile, error) {
	return s.repo.List()
}

// ListPaged возвращает страницу складских позиций: не более limit записей,
// начиная со смещения offset. При пустой странице возвращается пустой не-nil слайс.
func (s *Service) ListPaged(limit, offset int) ([]model.Stockpile, error) {
	return s.repo.ListPaged(limit, offset)
}

// Count возвращает общее количество складских позиций.
func (s *Service) Count() (int64, error) {
	return s.repo.Count()
}

// Update обновляет у складской позиции id только те поля input, которые не nil,
// и возвращает обновлённую запись. Если в input нет ни одного не-nil поля,
// возвращается ErrInvalidInput; если запись не найдена — ErrNotFound.
func (s *Service) Update(id int64, in model.UpdateStockpileInput) (model.Stockpile, error) {
	if !hasUpdateFields(in) {
		return model.Stockpile{}, ErrInvalidInput
	}
	if in.Dimension != nil && strings.TrimSpace(*in.Dimension) == "" {
		return model.Stockpile{}, fmt.Errorf("%w: dimension is required", ErrInvalidInput)
	}

	updated, err := s.repo.Update(id, in)
	switch {
	case errors.Is(err, repo.ErrNoUpdateFields):
		return model.Stockpile{}, ErrInvalidInput
	case err != nil:
		return model.Stockpile{}, err
	case !updated:
		return model.Stockpile{}, ErrNotFound
	}

	return s.GetByID(id)
}

// Delete удаляет складскую позицию по идентификатору.
// Возвращает ошибку только при проблемах БД; отсутствие записи не является ошибкой.
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}

// validateRequired проверяет обязательные текстовые поля новой записи.
// Требуются непустые Name и Dimension; остальные поля могут быть пустыми строками
// (NOT NULL в БД допускает "").
func validateRequired(in model.Stockpile) error {
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if strings.TrimSpace(in.Dimension) == "" {
		return fmt.Errorf("%w: dimension is required", ErrInvalidInput)
	}
	return nil
}

// hasUpdateFields определяет, есть ли в input хотя бы одно не-nil поле.
func hasUpdateFields(in model.UpdateStockpileInput) bool {
	return in.Name != nil ||
		in.Dimension != nil ||
		in.IconName != nil ||
		in.X != nil ||
		in.Y != nil ||
		in.Z != nil
}
