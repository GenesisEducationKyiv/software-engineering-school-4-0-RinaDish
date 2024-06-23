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
	"github.com/stretchr/testify/suite"
)

const email = "testemail@gmail.com"

type UserHandlerTestSuite struct {
	suite.Suite
	DB *gorm.DB
}

func (suite *UserHandlerTestSuite) SetupSuite() {
	suite.DB, _ = gorm.Open(postgres.Open(testdb.GetDBDSN()), &gorm.Config{})
}

func (suite *UserHandlerTestSuite) BeforeTest(suiteName, testName string) {
	db, _ := suite.DB.DB()
	db.Close()

	testdb.ResetDBTemplate(suite.T(), testdb.GetTemplateDBDSN(), testdb.DBName, testdb.TemplateDBName)

	suite.DB, _ = gorm.Open(postgres.Open(testdb.GetDBDSN()), &gorm.Config{})
}

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

func (suite *UserHandlerTestSuite) TestSuccessfulCreateSubscription() {
	response := sendMail(email, suite.DB)

	_, err := io.ReadAll(response.Body)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), http.StatusOK, response.StatusCode)

	actualEmail := testdb.GetEmail(suite.T(), testdb.GetDBDSN(), testdb.DBName, email)
	require.Equal(suite.T(), actualEmail, email)
}

func (suite *UserHandlerTestSuite) TestFailureDuplicateEmail() {
	_ = sendMail(email, suite.DB)
	response := sendMail(email, suite.DB)

	_, err := io.ReadAll(response.Body)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), http.StatusConflict, response.StatusCode)

	actualEmail := testdb.GetEmail(suite.T(), testdb.GetDBDSN(), testdb.DBName, email)
	require.Equal(suite.T(), actualEmail, email)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
