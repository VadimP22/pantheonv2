// Package stockpile_resources реализует сервисный слой для ресурсов склада.
// Сервис инкапсулирует правила изменения количеств (deposit/withdraw),
// сортировку списка ресурсов и синхронизацию конкурентных операций,
// работая напрямую с типами модели.
package stockpile_resources

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	model "backend/model/stockpile_resources"
	repo "backend/repository/stockpile_resources"
	stockpilesrepo "backend/repository/stockpiles"
)

// Ошибки уровня сервиса.
var (
	// ErrNotFound — склад или ресурс не найден.
	ErrNotFound = errors.New("stockpile_resources: not found")
	// ErrInvalidInput — невалидные данные запроса.
	ErrInvalidInput = errors.New("stockpile_resources: invalid input")
	// ErrInsufficient — на складе меньше ресурса, чем просят списать.
	ErrInsufficient = errors.New("stockpile_resources: insufficient amount")
)

// Service обеспечивает операции пополнения/списания ресурсов склада.
type Service struct {
	resources  *repo.Repository
	stockpiles *stockpilesrepo.Repository
	// mu сериализует операции mutate/read, чтобы deposit/withdraw не гонялись
	// между чтением текущего количества и записью нового.
	mu sync.Mutex
}

// New возвращает новый Service поверх репозиториев ресурсов и складов.
func New(resources *repo.Repository, stockpiles *stockpilesrepo.Repository) *Service {
	return &Service{resources: resources, stockpiles: stockpiles}
}

// Deposit вносит count единиц ресурса itemID на склад stockpileID.
// Если записи для пары «склад + ресурс» нет, она создаётся; если есть —
// количество увеличивается. Возвращает актуальное состояние записи.
func (s *Service) Deposit(stockpileID int64, itemID string, count int64) (model.StockpileResource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	itemID, err := s.validate(stockpileID, itemID, count)
	if err != nil {
		return model.StockpileResource{}, err
	}

	rec, err := s.find(stockpileID, itemID)
	if err != nil {
		return model.StockpileResource{}, err
	}

	if rec == nil {
		created := model.StockpileResource{
			StockpileID: stockpileID,
			ItemID:      itemID,
			ItemCount:   count,
		}
		if err := s.resources.Create(&created); err != nil {
			return model.StockpileResource{}, fmt.Errorf("deposit resource: %w", err)
		}
		if err := s.stockpiles.Touch(stockpileID); err != nil {
			return model.StockpileResource{}, fmt.Errorf("deposit resource: %w", err)
		}
		return created, nil
	}

	total := rec.ItemCount + count
	updated, err := s.resources.Update(rec.ID, model.UpdateStockpileResourceInput{ItemCount: &total})
	switch {
	case err != nil:
		return model.StockpileResource{}, fmt.Errorf("deposit resource: %w", err)
	case !updated:
		return model.StockpileResource{}, fmt.Errorf("%w: resource %d disappeared", ErrNotFound, rec.ID)
	}

	if err := s.stockpiles.Touch(stockpileID); err != nil {
		return model.StockpileResource{}, fmt.Errorf("deposit resource: %w", err)
	}

	rec.ItemCount = total
	return *rec, nil
}

// Withdraw списывает count единиц ресурса itemID со склада stockpileID.
// Если ресурса не хватает, возвращается ErrInsufficient; если после списания
// количество становится нулевым, запись удаляется.
func (s *Service) Withdraw(stockpileID int64, itemID string, count int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	itemID, err := s.validate(stockpileID, itemID, count)
	if err != nil {
		return err
	}

	rec, err := s.find(stockpileID, itemID)
	if err != nil {
		return err
	}
	if rec == nil {
		return fmt.Errorf("%w: resource %q in stockpile %d", ErrNotFound, itemID, stockpileID)
	}
	if rec.ItemCount < count {
		return fmt.Errorf("%w: have %d, requested %d", ErrInsufficient, rec.ItemCount, count)
	}

	rest := rec.ItemCount - count
	if rest == 0 {
		if err := s.resources.Delete(rec.ID); err != nil {
			return fmt.Errorf("withdraw resource: %w", err)
		}
		if err := s.stockpiles.Touch(stockpileID); err != nil {
			return fmt.Errorf("withdraw resource: %w", err)
		}
		return nil
	}

	updated, err := s.resources.Update(rec.ID, model.UpdateStockpileResourceInput{ItemCount: &rest})
	switch {
	case err != nil:
		return fmt.Errorf("withdraw resource: %w", err)
	case !updated:
		return fmt.Errorf("%w: resource %d disappeared", ErrNotFound, rec.ID)
	}

	if err := s.stockpiles.Touch(stockpileID); err != nil {
		return fmt.Errorf("withdraw resource: %w", err)
	}
	return nil
}

// ListByStockpile возвращает ресурсы склада, упорядоченные по убыванию количества.
// Неизвестный склад — ErrNotFound; при пустом складе — пустой не-nil слайс.
func (s *Service) ListByStockpile(stockpileID int64) ([]model.StockpileResource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if stockpileID < 1 {
		return nil, fmt.Errorf("%w: stockpile id must be positive", ErrInvalidInput)
	}
	if err := s.ensureStockpile(stockpileID); err != nil {
		return nil, err
	}

	items, err := s.resources.ListByStockpile(stockpileID)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].ItemCount != items[j].ItemCount {
			return items[i].ItemCount > items[j].ItemCount
		}
		if items[i].ItemID != items[j].ItemID {
			return items[i].ItemID < items[j].ItemID
		}
		return items[i].ID < items[j].ID
	})
	return items, nil
}

// validate проверяет аргументы операции и возвращает нормализованный itemID.
func (s *Service) validate(stockpileID int64, itemID string, count int64) (string, error) {
	if stockpileID < 1 {
		return "", fmt.Errorf("%w: stockpile id must be positive", ErrInvalidInput)
	}
	trimmed := strings.TrimSpace(itemID)
	if trimmed == "" {
		return "", fmt.Errorf("%w: item id is required", ErrInvalidInput)
	}
	if count < 1 {
		return "", fmt.Errorf("%w: count must be positive", ErrInvalidInput)
	}
	if err := s.ensureStockpile(stockpileID); err != nil {
		return "", err
	}
	return trimmed, nil
}

// ensureStockpile проверяет существование склада с указанным id.
func (s *Service) ensureStockpile(id int64) error {
	_, err := s.stockpiles.GetByID(id)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, sql.ErrNoRows):
		return fmt.Errorf("%w: stockpile %d", ErrNotFound, id)
	default:
		return fmt.Errorf("get stockpile %d: %w", id, err)
	}
}

// find ищет запись по паре «склад + предмет». Возвращает nil, nil при отсутствии.
func (s *Service) find(stockpileID int64, itemID string) (*model.StockpileResource, error) {
	rec, err := s.resources.GetByStockpileAndItem(stockpileID, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return rec, nil
}
