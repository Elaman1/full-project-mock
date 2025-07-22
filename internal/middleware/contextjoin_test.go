package middleware

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestContextJoinMiddleware_JoinedContextCancelledOnShutdown(t *testing.T) {
	// Создаём shutdown-контекст с возможностью отмены
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	defer shutdownCancel()

	// Запрос и ответ
	req := httptest.NewRequest(http.MethodGet, "/any", nil)
	rec := httptest.NewRecorder()

	// Переменная для захвата переданного в handler контекста
	var joinedCtx context.Context

	// Обработчик, в который попадёт уже объединённый контекст
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		joinedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем handler мидлваром
	middleware := ContextJoinMiddleware(shutdownCtx)
	wrappedHandler := middleware(handler)

	// Выполняем запрос
	wrappedHandler.ServeHTTP(rec, req)

	// Проверяем, что статус был установлен как ожидается
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotNil(t, joinedCtx)

	// Отменяем shutdown-контекст
	shutdownCancel()

	// Проверяем, что объединённый контекст тоже завершился
	select {
	case <-joinedCtx.Done():
		// Успешно завершился — всё ок
	case <-time.After(500 * time.Millisecond):
		t.Fatal("joined context was not cancelled after shutdownCtx cancel")
	}
}
