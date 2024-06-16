package currencyrates

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/RinaDish/currency-rates/integrations_tests/pkg/testdb"
	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/RinaDish/currency-rates/internal/repo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupUserHandler(t *testing.T) (handlers.SubscribeHandler, *repo.Repository) {
	logger := zap.NewNop().Sugar()

	adminRepository, err := repo.NewAdminRepository(testdb.GetDBDSN(), logger)
	if err != nil {
		t.Fatalf("failed to create admin repository: %v", err)
	}

	handler := handlers.NewSubscribeHandler(logger, adminRepository)

	return handler, adminRepository
}

func TestUserHandler(main *testing.T) {
	main.Run("create user", func(t *testing.T) {
		testdb.Reset(t, testdb.GetTemplateDBDSN(), testdb.DBName, testdb.TemplateDBName)

		userHandler, _ := setupUserHandler(t)

		form := url.Values{}
		form.Add("email", "testuser")

		formData := form.Encode()

		req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(formData))

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		userHandler.CreateSubscription(w, req) 

		response := w.Result()
		defer response.Body.Close()
		responseBody, err := io.ReadAll(response.Body)
		fmt.Println(responseBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, response.StatusCode)
	})
}
