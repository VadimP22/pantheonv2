// Package lib содержит общие HTTP-хелперы и константы для контроллеров.
package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"backend/dto"
)

// Константы пагинации по умолчанию.
const (
	DefaultPage     = 1
	DefaultPageSize = 6
	MaxPageSize     = 100
)

// WriteJSON записывает тело в формате JSON с заданным статусом.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError записывает сообщение об ошибке в виде {"error":"..."}.
func WriteError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// DecodeJSON декодирует тело запроса; пустое тело считается ошибкой (EOF).
func DecodeJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json body: %w", err)
	}
	return nil
}

// ParseID извлекает и валидирует id из path-переменной {id}.
func ParseID(r *http.Request) (int64, error) {
	s := r.PathValue("id")
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid id: %q", s)
	}
	return id, nil
}

// ParseListQuery разбирает параметры пагинации page и page_size из query-строки.
func ParseListQuery(r *http.Request) (dto.ListQueryDTO, error) {
	q := dto.ListQueryDTO{Page: DefaultPage, PageSize: DefaultPageSize}

	if v := r.URL.Query().Get("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return q, fmt.Errorf("invalid page: %q", v)
		}
		q.Page = n
	}
	if v := r.URL.Query().Get("page_size"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > MaxPageSize {
			return q, fmt.Errorf("invalid page_size: %q", v)
		}
		q.PageSize = n
	}
	return q, nil
}
