//go:build integration

package currencyrates

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/RinaDish/currency-rates/integrations_tests/pkg/testdb"
	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/RinaDish/currency-rates/internal/repo"
	"github.com/stretchr/testify/require"
)

const email = "testemail@gmail.com"

func setupUserHandler(db *gorm.DB) (handlers.SubscribeHandler, *repo.Repository) {
	logger := zap.NewNop().Sugar()

	adminRepository := repo.NewAdminRepository(db, logger)

	handler := handlers.NewSubscribeHandler(logger, adminRepository)

	return handler, adminRepository
}

func sendMail(email string, db *gorm.DB) *http.Response {
	userHandler, _ := setupUserHandler(db)
	form := url.Values{}
	form.Add("email", email)

	formData := form.Encode()

	req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(formData))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	userHandler.CreateSubscription(w, req)

	response := w.Result()
	defer response.Body.Close()

	return response
}

func TestUserHandler(main *testing.T) {
	main.Run("successful create subsctiption", func(t *testing.T) {
		testdb.ResetDBTemplate(t, testdb.GetTemplateDBDSN(), testdb.DBName, testdb.TemplateDBName)

		db, err := gorm.Open(postgres.Open(testdb.GetDBDSN()), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to create admin repository: %v", err)
		}

		defer func() {
			if db, err := db.DB(); err == nil {
				_ = db.Close()
			}
		}()

		response := sendMail(email, db)

		_, err = io.ReadAll(response.Body)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		rowEmail := testdb.GetEmail(t, testdb.GetDBDSN(), testdb.DBName, email)
		defer require.Equal(t, rowEmail, email)
	})

	main.Run("failure: duplicate email", func(t *testing.T) {
		testdb.ResetDBTemplate(t, testdb.GetTemplateDBDSN(), testdb.DBName, testdb.TemplateDBName)

		db, err := gorm.Open(postgres.Open(testdb.GetDBDSN()), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to create admin repository: %v", err)
		}

		defer func() {
			if db, err := db.DB(); err == nil {
				_ = db.Close()
			}
		}()

		_ = sendMail(email, db)
		response := sendMail(email, db)

		_, err = io.ReadAll(response.Body)

		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, response.StatusCode)

		rowEmail := testdb.GetEmail(t, testdb.GetDBDSN(), testdb.DBName, email)

		require.Equal(t, rowEmail, email)
	})
}
